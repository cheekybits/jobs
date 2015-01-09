package jobs

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// J is a job.
type J struct {
	// id is the unique job id.
	ID bson.ObjectId `bson:"_id" json:"id"`
	// Status is the current status of the job.
	Status Status `json:"status"`
	// Tries holds history of attempts.
	Tries []*Try `json:"tries"`
	// Created is when the job was created.
	Created time.Time `json:"created"`
	// RunAt is the time this job should run at (or after).
	RunAt time.Time `json:"runat" bson:"runat"`
	// Data is the user data for this job.
	Data map[string]interface{} `json:"data" bson:"data"`
	// Retries is the number of remaining attempts that will
	// be made to run this job.
	Retries int
	// RetryInterval is the time to wait after a failure before
	// trying to run the job again.
	RetryInterval time.Duration
}

// New creates a new job.
func New() *J {
	return &J{
		ID:      bson.NewObjectId(),
		Status:  StatusNew,
		Created: time.Now(),
		Data:    make(map[string]interface{}),
	}
}

// Put adds the jobs to the specified collection.
func Put(c *mgo.Collection, jobs ...*J) error {
	inters := make([]interface{}, len(jobs))
	for i, job := range jobs {
		inters[i] = job
	}
	return c.Insert(inters...)
}
