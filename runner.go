package jobs

import (
	"sync"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Try contains details of an attempt to run
// a job.
type Try struct {
	// Runner is the name of the runner that tried to
	// run this job.
	Runner string `json:"runner" bson:"runner"`
	// When is a timestamp of when the attempt took place.
	When time.Time `json:"when" bson:"when"`
	// Err is the error that was returned by JobFunc.
	Err string `json:"err,omitempty" bson:"err,omitempty"`
}

// JobFunc is the function that gets called for each
// job.
type JobFunc func(job *J) error

// Runner runs jobs.
type Runner struct {
	c        *mgo.Collection
	fn       JobFunc
	stop     chan struct{}
	Interval time.Duration
	err      error
	stoponce sync.Once
	name     string
	kind     string
}

// Start starts the process.
func (r *Runner) Start() error {
	r.stop = make(chan struct{})
	r.stoponce = sync.Once{}
	go func() {
		var job *J
	outside:
		for {
			iter := r.c.Find(bson.M{
				"status": bson.M{"$in": []interface{}{StatusNew, StatusWaiting}},
				"runat":  bson.M{"$lte": time.Now()},
				"kind":   r.kind,
			}).Sort("created").Iter()
			for iter.Next(&job) {

				var err error
				var changeInfo *mgo.ChangeInfo
				change := mgo.Change{
					Update:    bson.M{"$set": bson.M{"status": StatusWorking}},
					ReturnNew: true,
				}
				if changeInfo, err = r.c.Find(bson.M{
					"_id":    job.ID,
					"status": bson.M{"$in": []interface{}{StatusNew, StatusWaiting}},
				}).Apply(change, &job); err != nil {
					if err == mgo.ErrNotFound {
						// skip this one - someone else is dealing with it
						continue
					}
					r.err = err
					break
				}
				if changeInfo.Updated != 1 {
					// skip this one - someone else is dealing with it
					continue
				}

				jobErr := r.fn(job)

				// record this attempt
				try := &Try{
					When:   time.Now(),
					Runner: r.name,
				}
				job.Tries = append(job.Tries, try)

				if jobErr != nil {

					try.Err = jobErr.Error()

					job.RunAt = time.Now().Add(job.RetryInterval)
					job.Retries--
					job.Status = StatusWaiting
					if job.Retries == 0 {
						job.Status = StatusFailed
					}

				} else {
					// success
					job.Status = StatusSuccess
				}

				if err := r.c.UpdateId(job.ID, job); err != nil {
					r.err = err
					break
				}

			}
			if err := iter.Close(); err != nil {
				r.err = err
			}
			if r.err != nil {
				r.Stop()
			}
			select {
			case <-r.stop:
				// stop
				break outside
			case <-time.After(r.Interval):
				// carry on
			}
		}

	}()

	return nil
}

// Stop stops the runner. Callers should then block on StopChan()
// to be notified of when the runner has stopped.
func (r *Runner) Stop() {
	r.stoponce.Do(func() {
		close(r.stop)
	})
}

// StopChan is a channel that gets closed when the runner
// has stopped. Callers should block on this after calling
// Stop to ensure the runner has properly stopped.
//     <-runner.StopChan()
func (r *Runner) StopChan() <-chan struct{} { return r.stop }

// Err is the last error that occurred.
func (r *Runner) Err() error { return r.err }

// Name is the name of the runner.
func (r *Runner) Name() string { return r.name }

// Kind is the name of the types of jobs this runner will process.
func (r *Runner) Kind() string { return r.kind }

// NewRunner makes a new Runner capable of running jobs. The name should be unique
// across a system. The mgo.Collection is where the job records live in MongoDB.
// The kind string should match J.Kind values for jobs that this runner should execute.
func NewRunner(name string, c *mgo.Collection, kind string, fn JobFunc) *Runner {
	return &Runner{
		c:        c,
		fn:       fn,
		Interval: 500 * time.Millisecond,
		name:     name,
		kind:     kind,
	}
}
