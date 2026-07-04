package scheduler

import (
	"context"
	"log/slog"
	"time"
)

type SubscriptionUpdater struct {
	log      *slog.Logger
	interval time.Duration
	execute  func(ctx context.Context) error
}

func NewSubscriptionUpdater(log *slog.Logger, interval time.Duration, execute func(ctx context.Context) error) *SubscriptionUpdater {
	return &SubscriptionUpdater{log: log, interval: interval, execute: execute}
}

// Start blocks, triggering execute on each tick until the context is
// cancelled. Execution errors are logged and do not stop the ticker.
func (su *SubscriptionUpdater) Start(ctx context.Context) {
	ticker := time.NewTicker(su.interval)
	defer ticker.Stop()

	su.log.Info("background subscription updater started", "interval", su.interval)

	for {
		select {
		case <-ctx.Done():
			su.log.Info("background subscription updater stopped")
			return
		case <-ticker.C:
			su.log.Debug("ticker triggered: updating subscriptions...")
			if err := su.execute(ctx); err != nil {
				su.log.Error("failed to update subscriptions in background", "error", err)
			}
		}
	}
}
