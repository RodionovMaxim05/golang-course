package domain

// Status represents the processing state of a repository within the system.
type Status string

const (
	StatusSuccess Status = "SUCCESS"
	StatusPending Status = "PENDING"
	StatusError   Status = "ERROR"
	StatusUnknown Status = "UNKNOWN"
)

func (s Status) String() string {
	return string(s)
}

func ParseStatus(s string) Status {
	status := Status(s)

	switch status {
	case StatusSuccess, StatusPending, StatusError:
		return status
	default:
		return StatusUnknown
	}
}
