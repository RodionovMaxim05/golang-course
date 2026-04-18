package collector

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	collectorpb "repo-stat/proto/collector"
	processorpb "repo-stat/proto/processor"
)

type Client struct {
	log  *slog.Logger
	conn *grpc.ClientConn
	pb   collectorpb.RepoServiceClient
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
		pb:   collectorpb.NewRepoServiceClient(conn),
	}, nil
}

func (c *Client) GetRepo(ctx context.Context, name, repo string) (*collectorpb.GetRepoResponse, error) {
	req := &collectorpb.GetRepoRequest{
		Name: name,
		Repo: repo,
	}
	resp, err := c.pb.GetRepo(ctx, req)
	if err != nil {
		c.log.Error("collector get repo failed", "error", err)
		return nil, err
	}
	return resp, nil
}

func (c *Client) GetSubscriptionsInfo(ctx context.Context, req *processorpb.GetSubsInfoRequest) (*processorpb.GetSubsInfoResponse, error) {
	resp, err := c.pb.GetSubscriptionsInfo(ctx, &collectorpb.GetSubsInfoRequest{})
	if err != nil {
		c.log.Error("collector get subscriptions info failed", "error", err)
		return nil, err
	}

	repositories := make([]*processorpb.GetRepoResponse, 0, len(resp.Repositories))
	for _, repo := range resp.Repositories {
		repositories = append(repositories, &processorpb.GetRepoResponse{
			FullName:        repo.FullName,
			Description:     repo.Description,
			StargazersCount: repo.StargazersCount,
			ForksCount:      repo.ForksCount,
			CreatedAt:       repo.CreatedAt,
		})
	}

	return &processorpb.GetSubsInfoResponse{Repositories: repositories}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
