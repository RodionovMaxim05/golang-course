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

type KafkaResultProducer interface {
	SendRepoResult(ctx context.Context, repo domain.Repository) error
}

type GetSubscriptionsInfoUsecase struct {
	subscriberClient SubscriberClient
	githubClient     GitHubClient
	kafkaProducer    KafkaResultProducer
}

func NewGetSubscriptionsInfoUsecase(subscriberClient SubscriberClient, githubClient GitHubClient, kafkaProducer KafkaResultProducer) *GetSubscriptionsInfoUsecase {
	return &GetSubscriptionsInfoUsecase{
		subscriberClient: subscriberClient,
		githubClient:     githubClient,
		kafkaProducer:    kafkaProducer,
	}
}

func (gsiu *GetSubscriptionsInfoUsecase) Execute(ctx context.Context) error {
	// Get subscriptions from subscriber service
	subscriptions, err := gsiu.subscriberClient.GetSubscriptions(ctx)
	if err != nil {
		return err
	}

	for _, sub := range subscriptions {
		// Get repository info for each subscription
		repo, err := gsiu.githubClient.GetRepo(ctx, sub.Owner, sub.Repo)
		if err != nil {
			continue
		}

		// Send the finished result to Kafka
		err = gsiu.kafkaProducer.SendRepoResult(ctx, repo)
		if err != nil {
			continue
		}
	}

	return nil
}
