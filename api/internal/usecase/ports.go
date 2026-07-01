package usecase

import (
	"context"

	"repo-watcher/api/internal/domain"
)

type SubscriberPinger interface {
	Ping(ctx context.Context) (domain.PingStatus, error)
}
