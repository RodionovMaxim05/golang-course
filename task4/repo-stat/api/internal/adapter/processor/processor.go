package processor

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	grpcAdapter "repo-stat/api/internal/adapter/grpc"
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

func (c *Client) GetRepo(ctx context.Context, name, repo string) (domain.Repository, error) {
	req := &processorpb.GetRepoRequest{
		Name: name,
		Repo: repo,
	}

	resp, err := c.pb.GetRepo(ctx, req)
	if err != nil {
		c.log.Error("processor get repo failed", "error", err)
		return domain.Repository{}, grpcAdapter.ErrToDomain(err)
	}

	return domain.Repository{
		Owner:       resp.Name,
		Repo:        resp.Repo,
		Description: resp.Description,
		Stargazers:  int(resp.StargazersCount),
		Forks:       int(resp.ForksCount),
		CreatedAt:   resp.CreatedAt.AsTime(),
	}, nil
}

func (c *Client) GetSubscriptionsInfo(ctx context.Context) ([]domain.Repository, error) {
	req := &processorpb.GetSubsInfoRequest{}

	resp, err := c.pb.GetSubscriptionsInfo(ctx, req)
	if err != nil {
		c.log.Error("processor get subscriptions info failed", "error", err)
		return nil, grpcAdapter.ErrToDomain(err)
	}

	result := make([]domain.Repository, 0, len(resp.Repositories))
	for _, item := range resp.Repositories {
		result = append(result, domain.Repository{
			Owner:       item.Name,
			Repo:        item.Repo,
			Description: item.Description,
			Stargazers:  int(item.StargazersCount),
			Forks:       int(item.ForksCount),
			CreatedAt:   item.CreatedAt.AsTime(),
		})
	}

	return result, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
