package usecase

import (
	"context"
	"log/slog"

	"repo-stat/processor/internal/domain"
)

type SubscribeClient interface {
	GetActiveSubscriptions(ctx context.Context) ([]string, error)
}

type DataStorage interface {
	GetReposByNames(ctx context.Context, names []string) ([]*domain.Repository, error)
}

type GetSubscriptionsInfo struct {
	log         *slog.Logger
	subClient   SubscribeClient
	dataStorage DataStorage
}

func NewGetSubscriptionsInfo(subClient SubscribeClient, storage DataStorage, log *slog.Logger) *GetSubscriptionsInfo {
	return &GetSubscriptionsInfo{
		log:         log,
		subClient:   subClient,
		dataStorage: storage,
	}
}

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
