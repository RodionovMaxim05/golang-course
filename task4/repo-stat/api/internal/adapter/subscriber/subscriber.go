package subscriber

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	grpcAdapter "repo-stat/api/internal/adapter/grpc"
	"repo-stat/api/internal/domain"
	subscriberpb "repo-stat/proto/subscriber"
)

type Client struct {
	log  *slog.Logger
	conn *grpc.ClientConn
	pb   subscriberpb.SubscriberClient
}

func NewClient(address string, log *slog.Logger) (*Client, error) {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		log:  log,
		conn: conn,
		pb:   subscriberpb.NewSubscriberClient(conn),
	}, nil
}

func (c *Client) Ping(ctx context.Context) domain.PingStatus {
	_, err := c.pb.Ping(ctx, &subscriberpb.PingRequest{})
	if err != nil {
		c.log.Error("subscriber ping failed", "error", err)
		return domain.PingStatusDown
	}

	return domain.PingStatusUp
}

func (c *Client) Subscribe(ctx context.Context, owner, repo string) (*domain.Subscription, error) {
	req := &subscriberpb.SubscribeRequest{
		Subscription: &subscriberpb.Subscription{
			Owner: owner,
			Repo:  repo,
		},
	}

	resp, err := c.pb.Subscribe(ctx, req)
	if err != nil {
		c.log.Error("subscriber subscribe failed", "error", err)
		return nil, grpcAdapter.ErrToDomain(err)
	}

	return &domain.Subscription{
		Owner:     resp.Subscription.Owner,
		Repo:      resp.Subscription.Repo,
		CreatedAt: resp.Subscription.CreatedAt.AsTime(),
	}, nil
}

func (c *Client) Unsubscribe(ctx context.Context, owner, repo string) error {
	req := &subscriberpb.UnsubscribeRequest{
		Subscription: &subscriberpb.Subscription{
			Owner: owner,
			Repo:  repo,
		},
	}

	_, err := c.pb.Unsubscribe(ctx, req)
	if err != nil {
		c.log.Error("subscriber unsubscribe failed", "error", err)
		return grpcAdapter.ErrToDomain(err)
	}

	return nil
}

func (c *Client) GetSubscriptions(ctx context.Context) ([]domain.Subscription, error) {
	resp, err := c.pb.GetSubscriptions(ctx, &subscriberpb.GetSubsRequest{})
	if err != nil {
		c.log.Error("subscriber getSubscriptions failed", "error", err)
		return nil, grpcAdapter.ErrToDomain(err)
	}

	subs := make([]domain.Subscription, 0, len(resp.Subscriptions))
	for _, s := range resp.Subscriptions {
		subs = append(subs, domain.Subscription{
			Owner:     s.Owner,
			Repo:      s.Repo,
			CreatedAt: s.CreatedAt.AsTime(),
		})
	}

	return subs, nil
}
func (c *Client) Close() error {
	return c.conn.Close()
}
