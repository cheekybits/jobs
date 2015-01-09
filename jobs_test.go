package jobs_test

import (
	"github.com/cheekybits/is"
	"github.com/cheekybits/jobs"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"testing"
	"time"
)

func test(isI is.I, fn func(is is.I, db *mgo.Database)) {
	session, err := mgo.DialWithTimeout("localhost", 500*time.Millisecond)
	isI.NoErr(err)
	defer session.Close()
	db := session.DB("jobs-database")
	isI.NoErr(db.DropDatabase())
	fn(isI, db)
}

func TestNew(t *testing.T) {
	test(is.New(t), func(is is.I, db *mgo.Database) {

		thingID := bson.NewObjectId()
		thing := bson.M{"_id": thingID}
		is.NoErr(db.C("things").Insert(thing))

		job := jobs.New("things")
		now := time.Now()
		job.RunAt = now
		job.Data["something"] = true
		err := jobs.Put(db.C("jobs"), job)
		is.NoErr(err)

		var result map[string]interface{}
		is.NoErr(db.C("jobs").Find(nil).Limit(1).One(&result))

		is.OK(result)
		is.Equal(job.ID, result["_id"])
		is.NoErr(err)
		is.Equal(result["runat"].(time.Time).Format("20060102"), now.Format("20060102"))
		is.Equal(result["status"], jobs.StatusNew)
		is.OK(result["tries"])
		is.Equal(0, len(result["tries"].([]interface{})))
		is.Equal(true, result["data"].(map[string]interface{})["something"])

	})
}
