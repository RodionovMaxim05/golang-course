package usecase

import (
	"context"

	"repo-stat/collector/internal/domain"
)

type SubscriberClient interface {
	GetSubscriptions(ctx context.Context) ([]domain.Subscription, error)
}

type GitHubClient interface {
	GetRepo(ctx context.Context, owner, name string) (domain.Repository, error)
}

type GetSubscriptionsInfoUsecase struct {
	subscriberClient SubscriberClient
	githubClient     GitHubClient
}

func NewGetSubscriptionsInfoUsecase(subscriberClient SubscriberClient, githubClient GitHubClient) *GetSubscriptionsInfoUsecase {
	return &GetSubscriptionsInfoUsecase{
		subscriberClient: subscriberClient,
		githubClient:     githubClient,
	}
}

func (gsiu *GetSubscriptionsInfoUsecase) Execute(ctx context.Context) ([]domain.Repository, error) {
	// Get subscriptions from subscriber service
	subscriptions, err := gsiu.subscriberClient.GetSubscriptions(ctx)
	if err != nil {
		return nil, err
	}

	// Get repository info for each subscription
	results := make([]domain.Repository, 0, len(subscriptions))
	for _, sub := range subscriptions {
		repo, err := gsiu.githubClient.GetRepo(ctx, sub.Owner, sub.Repo)
		if err != nil {
			return nil, err
		}
		results = append(results, repo)
	}

	return results, nil
}
