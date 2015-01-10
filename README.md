# jobs

Job management for MongoDB (via mgo).

## Usage

```
// create a job and set some data
job := jobs.New("notifications")
job.Data["message"] = "Hello world"

// put the job
jobs.Put(db.C("jobs"), job)
```

Meanwhile, in a process far, far away:

```
r := jobs.NewRunner("runner-1", db.C("jobs"), "notifications", func(j *jobs.J) error {
	log.Println("TODO: process this message -", j.Data["message"])
	return nil
})
r.Start()
```