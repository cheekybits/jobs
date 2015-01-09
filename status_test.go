package jobs

import (
	"encoding/json"
	"testing"

	"github.com/cheekybits/is"
)

func TestStatus(t *testing.T) {
	is := is.New(t)

	is.Equal(StatusNew.String(), "new")
	is.Equal(StatusFailed.String(), "failed")
	is.Equal(StatusSuccess.String(), "success")
	is.Equal(StatusWaiting.String(), "waiting")
	is.Equal(StatusWorking.String(), "working")

	is.Equal(StatusInvalid, parseStatus("bollocks"))
	is.Equal(StatusNew, parseStatus("new"))
	is.Equal(StatusSuccess, parseStatus("success"))
	is.Equal(StatusFailed, parseStatus("failed"))
	is.Equal(StatusWaiting, parseStatus("waiting"))
	is.Equal(StatusWorking, parseStatus("working"))

	b, err := json.Marshal(StatusNew)
	is.NoErr(err)
	is.Equal(string(b), `"new"`)

	// var s Status
	// is.NoErr(json.Unmarshal([]byte(`"success"`), &s))
	// is.Equal(s, StatusSuccess)

	// obj := bson.M{"status": StatusWaiting}
	// b, err = bson.Marshal(obj)
	// is.NoErr(err)
	// is.Equal(true, bytes.Contains(b, []byte(StatusWaiting.String())))
	// var obj2 bson.M
	// is.NoErr(bson.Unmarshal(b, &obj2))
	// is.Equal(obj2["status"], StatusWaiting)

}
