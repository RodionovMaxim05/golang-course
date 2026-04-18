package subscriber

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"repo-stat/collector/internal/domain"
	subscriberpb "repo-stat/proto/subscriber"
)

type SubscriberClient struct {
	log    *slog.Logger
	client subscriberpb.SubscriberClient
}

func NewSubscriberClient(address string, log *slog.Logger) (*SubscriberClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &SubscriberClient{log: log, client: subscriberpb.NewSubscriberClient(conn)}, nil
}

func (sc *SubscriberClient) GetSubscriptions(ctx context.Context) ([]domain.Subscription, error) {
	sc.log.Debug("fetching subscriptions from subscriber service")

	resp, err := sc.client.GetSubscriptions(ctx, &subscriberpb.GetSubsRequest{})
	if err != nil {
		sc.log.Error("subscriber get subscriptions failed", "error", err)
		return nil, err
	}

	subscriptions := make([]domain.Subscription, 0, len(resp.Subscriptions))
	for _, s := range resp.Subscriptions {
		subscriptions = append(subscriptions, domain.Subscription{
			Owner:     s.Owner,
			Repo:      s.Repo,
			CreatedAt: s.CreatedAt.AsTime(),
		})
	}

	return subscriptions, nil
}
