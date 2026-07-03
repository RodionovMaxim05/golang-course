package usecase

import (
	"context"
	"log/slog"

	"repo-watcher/processor/internal/domain"
)

type SubscriptionsClient interface {
	GetActiveSubscriptions(ctx context.Context) ([]string, error)
}

type SubscriptionsRepoReader interface {
	GetReposByNames(ctx context.Context, names []string) ([]*domain.Repository, error)
}

type GetSubscriptionsInfo struct {
	log         *slog.Logger
	subClient   SubscriptionsClient
	dataStorage SubscriptionsRepoReader
}

func NewGetSubscriptionsInfo(subClient SubscriptionsClient, storage SubscriptionsRepoReader, log *slog.Logger) *GetSubscriptionsInfo {
	return &GetSubscriptionsInfo{
		log:         log,
		subClient:   subClient,
		dataStorage: storage,
	}
}

// Execute retrieves the current list of active subscriptions from the
// Subscriber service and resolves their cached repository metrics from
// local storage. Returns an empty slice if there are no active subscriptions.
func (gsi *GetSubscriptionsInfo) Execute(ctx context.Context) ([]*domain.Repository, error) {
	gsi.log.Debug("executing GetSubscriptionsInfo usecase")

	activeSubs, err := gsi.subClient.GetActiveSubscriptions(ctx)
	if err != nil {
		gsi.log.Error("failed to fetch active subscriptions from subscribe service", "error", err)
		return nil, err
	}

	if len(activeSubs) == 0 {
		gsi.log.Debug("no active subscriptions found")
		return []*domain.Repository{}, nil
	}

	repos, err := gsi.dataStorage.GetReposByNames(ctx, activeSubs)
	if err != nil {
		gsi.log.Error("failed to fetch repos by names from local db", "error", err)
		return nil, err
	}

	return repos, nil
}
