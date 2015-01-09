package jobs

// Status is the status of a job.
type Status int8

const (
	// StatusInvalid represents invalid status values.
	StatusInvalid Status = iota
	// StatusNew means a job has just been created.
	StatusNew
	// StatusWorking means the job is being worked on.
	StatusWorking
	// StatusWaiting means the job is waiting to retry.
	StatusWaiting
	// StatusSuccess means the job was successful.
	StatusSuccess
	// StatusFailed means the job failed.
	StatusFailed
)

// String gets the string for the status.
func (s Status) String() string {
	return statusStrs[s]
}

// MarshalText marshals the value into bytes.
func (s Status) MarshalText() (text []byte, err error) {
	return []byte(s.String()), nil
}

// UnmarshalText unmarshals the bytes into the value.
func (s *Status) UnmarshalText(text []byte) error {
	*s = parseStatus(string(text))
	return nil
}

// statusStrs are status strings.
var statusStrs = map[Status]string{
	StatusInvalid: "invalid",
	StatusNew:     "new",
	StatusWorking: "working",
	StatusWaiting: "waiting",
	StatusSuccess: "success",
	StatusFailed:  "failed",
}
var strStatuses = map[string]Status{
	"new":     StatusNew,
	"working": StatusWorking,
	"waiting": StatusWaiting,
	"success": StatusSuccess,
	"failed":  StatusFailed,
}

// parseStatus parses a status string.
func parseStatus(s string) Status {
	if st, ok := strStatuses[s]; ok {
		return st
	}
	return StatusInvalid
}
