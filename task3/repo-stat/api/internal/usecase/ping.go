package usecase

import (
	"context"

	"repo-stat/api/internal/domain"
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

func (u *Ping) Execute(ctx context.Context) (domain.PingStatus, domain.PingStatus) {
	subscriberStatus := u.subscriberPinger.Ping(ctx)
	processorStatus := u.processorPinger.Ping(ctx)
	return subscriberStatus, processorStatus
}
