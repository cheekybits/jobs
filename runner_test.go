package jobs_test

import (
	"errors"
	"testing"
	"time"

	"github.com/cheekybits/is"
	"github.com/cheekybits/jobs"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func TestRunnerEverythingOK(t *testing.T) {
	test(is.New(t), func(is is.I, db *mgo.Database) {

		// make three things
		colThings := db.C("things")
		thing1ID := bson.NewObjectId()
		thing1 := bson.M{"_id": thing1ID, "thing": 1}
		thing2ID := bson.NewObjectId()
		thing2 := bson.M{"_id": thing2ID, "thing": 2}
		thing3ID := bson.NewObjectId()
		thing3 := bson.M{"_id": thing3ID, "thing": 3}
		is.NoErr(colThings.Insert(thing1, thing2, thing3))

		// make three jobs
		job1 := jobs.New()
		job1.Data["thing_id"] = thing1ID
		job2 := jobs.New()
		job2.Data["thing_id"] = thing2ID
		job3 := jobs.New()
		job3.Data["thing_id"] = thing3ID
		is.NoErr(jobs.Put(db.C("jobs"), job1, job2, job3))

		var things []bson.M
		runner := jobs.NewRunner("test-1", db.C("jobs"), func(job *jobs.J) error {
			var thing map[string]interface{}
			is.NoErr(db.C("things").FindId(job.Data["thing_id"].(bson.ObjectId)).One(&thing))
			things = append(things, thing)
			return nil
		})
		is.NoErr(runner.Start())

		go func() {
			time.Sleep(500 * time.Millisecond)
			runner.Stop()
		}()

		select {
		case <-time.After(1 * time.Second):
			is.Fail("timed out")
		case <-runner.StopChan():
		}

		is.NoErr(runner.Err())
		is.Equal(3, len(things))
		is.Equal(1, things[0]["thing"])
		is.Equal(2, things[1]["thing"])
		is.Equal(3, things[2]["thing"])

		var job *jobs.J
		iter := db.C("jobs").Find(nil).Iter()
		for iter.Next(&job) {
			is.Equal(job.Status, jobs.StatusSuccess)
			is.Equal(len(job.Tries), 1)
			for _, try := range job.Tries {
				is.NotEqual(time.Time{}, try.When)
				is.Equal(runner.Name(), try.Runner)
			}
		}
		is.NoErr(iter.Err())

	})
}

func TestRunnerJobFailure(t *testing.T) {
	test(is.New(t), func(is is.I, db *mgo.Database) {

		// make three things
		colThings := db.C("things")
		thing1ID := bson.NewObjectId()
		thing1 := bson.M{"_id": thing1ID, "thing": 1}
		thing2ID := bson.NewObjectId()
		thing2 := bson.M{"_id": thing2ID, "thing": 2}
		thing3ID := bson.NewObjectId()
		thing3 := bson.M{"_id": thing3ID, "thing": 3}
		is.NoErr(colThings.Insert(thing1, thing2, thing3))

		// make three jobs
		job1 := jobs.New()
		job1.Data["thing_id"] = thing1ID
		job1.Data["job_id"] = 1
		job2 := jobs.New()
		job2.Data["thing_id"] = thing2ID
		job2.Data["job_id"] = 2
		job2.Retries = 5
		job2.RetryInterval = 1 * time.Millisecond
		job3 := jobs.New()
		job3.Data["thing_id"] = thing3ID
		job3.Data["job_id"] = 3
		is.NoErr(jobs.Put(db.C("jobs"), job1, job2, job3))

		failureCount := 0
		var things []bson.M
		testErr := errors.New("something went wrong")
		runner := jobs.NewRunner("test-1", db.C("jobs"), func(job *jobs.J) error {

			var thing map[string]interface{}
			is.NoErr(db.C("things").FindId(job.Data["thing_id"].(bson.ObjectId)).One(&thing))

			if thing["thing"] == 2 && failureCount < 3 {
				failureCount++
				return testErr
			}

			things = append(things, thing)
			return nil
		})
		runner.Interval = 1 * time.Millisecond
		is.NoErr(runner.Start())

		go func() {
			time.Sleep(200 * time.Millisecond)
			runner.Stop()
		}()

		select {
		case <-time.After(300 * time.Millisecond):
			is.Fail("timed out")
		case <-runner.StopChan():
			// finished
		}

		is.NoErr(runner.Err())
		is.Equal(3, len(things))
		is.Equal(1, things[0]["thing"])
		is.Equal(3, things[1]["thing"])
		is.Equal(2, things[2]["thing"])

		var job *jobs.J
		is.NoErr(db.C("jobs").FindId(job2.ID).One(&job))
		is.Equal(len(job.Tries), 4)
		is.Equal(job.Status, jobs.StatusSuccess)

	})
}

func TestRunnerJobCompleteFailure(t *testing.T) {
	test(is.New(t), func(is is.I, db *mgo.Database) {

		// make three things
		colThings := db.C("things")
		thing1ID := bson.NewObjectId()
		thing1 := bson.M{"_id": thing1ID, "thing": 1}
		thing2ID := bson.NewObjectId()
		thing2 := bson.M{"_id": thing2ID, "thing": 2}
		thing3ID := bson.NewObjectId()
		thing3 := bson.M{"_id": thing3ID, "thing": 3}
		is.NoErr(colThings.Insert(thing1, thing2, thing3))

		// make three jobs
		job1 := jobs.New()
		job1.Data["thing_id"] = thing1ID
		job1.Data["job_id"] = 1
		job2 := jobs.New()
		job2.Data["thing_id"] = thing2ID
		job2.Data["job_id"] = 2
		job2.Retries = 5
		job2.RetryInterval = 1 * time.Millisecond
		job3 := jobs.New()
		job3.Data["thing_id"] = thing3ID
		job3.Data["job_id"] = 3
		is.NoErr(jobs.Put(db.C("jobs"), job1, job2, job3))

		var things []bson.M
		testErr := errors.New("something went wrong")
		runner := jobs.NewRunner("test-1", db.C("jobs"), func(job *jobs.J) error {

			var thing map[string]interface{}
			is.NoErr(db.C("things").FindId(job.Data["thing_id"].(bson.ObjectId)).One(&thing))

			if thing["thing"] == 2 {
				return testErr
			}

			things = append(things, thing)
			return nil
		})
		runner.Interval = 1 * time.Millisecond
		is.NoErr(runner.Start())

		go func() {
			time.Sleep(200 * time.Millisecond)
			runner.Stop()
		}()

		select {
		case <-time.After(300 * time.Millisecond):
			is.Fail("timed out")
		case <-runner.StopChan():
			// finished
		}

		is.NoErr(runner.Err())
		is.Equal(2, len(things))
		is.Equal(1, things[0]["thing"])
		is.Equal(3, things[1]["thing"])

		var job *jobs.J
		is.NoErr(db.C("jobs").FindId(job2.ID).One(&job))
		is.Equal(len(job.Tries), 5)
		is.Equal(job.Status, jobs.StatusFailed)

	})
}
