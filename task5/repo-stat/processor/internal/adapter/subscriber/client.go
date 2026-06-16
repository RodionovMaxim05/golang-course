package subscriber

import (
	"context"
	"fmt"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	subscribepb "repo-stat/proto/subscriber"
)

type Client struct {
	log    *slog.Logger
	client subscribepb.SubscriberClient
	conn   *grpc.ClientConn
}

func NewClient(addr string, log *slog.Logger) (*Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("connect to subscribe service: %w", err)
	}

	return &Client{
		log:    log,
		client: subscribepb.NewSubscriberClient(conn),
		conn:   conn,
	}, nil
}

func (c *Client) GetActiveSubscriptions(ctx context.Context) ([]string, error) {
	resp, err := c.client.GetSubscriptions(ctx, &subscribepb.GetSubsRequest{})
	if err != nil {
		return nil, err
	}

	fullNames := make([]string, 0, len(resp.GetSubscriptions()))
	for _, sub := range resp.GetSubscriptions() {
		fullNames = append(fullNames, sub.Owner+"/"+sub.Repo)
	}

	return fullNames, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
