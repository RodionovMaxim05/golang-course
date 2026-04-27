package processor

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"repo-stat/api/internal/domain"
	processorpb "repo-stat/proto/processor"
)

type Client struct {
	log  *slog.Logger
	conn *grpc.ClientConn
	pb   processorpb.ProcessorClient
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
		pb:   processorpb.NewProcessorClient(conn),
	}, nil
}

func (c *Client) Ping(ctx context.Context) domain.PingStatus {
	_, err := c.pb.Ping(ctx, &processorpb.PingRequest{})
	if err != nil {
		c.log.Error("processor ping failed", "error", err)
		return domain.PingStatusDown
	}

	return domain.PingStatusUp
}

func (c *Client) GetRepo(ctx context.Context, name, repo string) (*processorpb.GetRepoResponse, error) {
	req := &processorpb.GetRepoRequest{
		Name: name,
		Repo: repo,
	}
	resp, err := c.pb.GetRepo(ctx, req)
	if err != nil {
		c.log.Error("processor get repo failed", "error", err)
		return nil, err
	}
	return resp, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
