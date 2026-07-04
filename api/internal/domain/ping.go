package domain

type PingStatus string

// PingStatus represents the liveness state of a downstream service.
const (
	PingStatusUp   PingStatus = "up"
	PingStatusDown PingStatus = "down"
)

// PingResult holds the liveness status of each downstream service.
type PingResult struct {
	Subscriber PingStatus
	Processor  PingStatus
}
