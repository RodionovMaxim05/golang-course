package usecase

import (
	"context"

	"repo-watcher/api/internal/domain"
)

type Pinger interface {
	Ping(ctx context.Context) domain.PingStatus
}

type Ping struct {
	subscriberPinger Pinger
	processorPinger  Pinger
}

func NewPing(subscriberPinger, processorPinger Pinger) *Ping {
	return &Ping{
		subscriberPinger: subscriberPinger,
		processorPinger:  processorPinger,
	}
}

// Execute checks the liveness of the Subscriber and Processor services
// and returns their individual statuses.
func (u *Ping) Execute(ctx context.Context) domain.PingResult {
	return domain.PingResult{
		Subscriber: u.subscriberPinger.Ping(ctx),
		Processor:  u.processorPinger.Ping(ctx),
	}
}
