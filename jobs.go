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
	// TargetDB is the target database name.
	TargetDB string `json:"target_db" bson:"target_db"`
	// TargetCollection is the target collection name.
	TargetCollection string `json:"target_col" bson:"target_col"`
	// TargetID is the ID of the target, inside the TargetCollection.
	TargetID bson.ObjectId `json:"target_id" bson:"target_id"`
}

// New creates a new job for the specified target.
func New(target *mgo.Collection, targetID bson.ObjectId) *J {
	return &J{
		ID:               bson.NewObjectId(),
		Status:           StatusNew,
		Created:          time.Now(),
		TargetDB:         target.Database.Name,
		TargetCollection: target.Name,
		TargetID:         targetID,
	}
}

// Put adds a new job to the specified collection.
func Put(jobs *mgo.Collection, job *J) error {
	return jobs.Insert(job)
}

// Try contains details of an attempt to run
// a job.
type Try struct {
	When time.Time
}
