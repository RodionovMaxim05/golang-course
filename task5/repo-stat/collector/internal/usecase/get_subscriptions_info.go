package usecase

import (
	"context"
	"log/slog"

	"repo-stat/collector/internal/domain"
)

type SubscriberClient interface {
	GetSubscriptions(ctx context.Context) ([]domain.Subscription, error)
}

type TaskProducer interface {
	SendCollectionTask(ctx context.Context, owner, repo string) error
}

type GetSubscriptionsInfoUsecase struct {
	log              *slog.Logger
	subscriberClient SubscriberClient
	taskProducer     TaskProducer
}

func NewGetSubscriptionsInfoUsecase(log *slog.Logger, subscriberClient SubscriberClient, taskProducer TaskProducer,
) *GetSubscriptionsInfoUsecase {
	return &GetSubscriptionsInfoUsecase{
		log:              log,
		subscriberClient: subscriberClient,
		taskProducer:     taskProducer,
	}
}

func (gsiu *GetSubscriptionsInfoUsecase) Execute(ctx context.Context) error {
	// Get subscriptions from subscriber service
	subscriptions, err := gsiu.subscriberClient.GetSubscriptions(ctx)
	if err != nil {
		gsiu.log.Error("failed to get subscriptions from subscriber client", "error", err)
		return err
	}

	// Submit a request to receive subscription data
	for _, sub := range subscriptions {
		err = gsiu.taskProducer.SendCollectionTask(ctx, sub.Owner, sub.Repo)
		if err != nil {
			gsiu.log.Error("failed to publish collection task to producer",
				"owner", sub.Owner,
				"repo", sub.Repo,
				"error", err,
			)
			continue
		}
	}

	return nil
}
