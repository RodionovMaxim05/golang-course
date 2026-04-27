package collector

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"repo-stat/processor/internal/domain"
	collectorpb "repo-stat/proto/collector"
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

func (c *Client) GetRepo(ctx context.Context, name, repo string) (*domain.Repository, error) {
	req := &collectorpb.GetRepoRequest{
		Name: name,
		Repo: repo,
	}

	resp, err := c.pb.GetRepo(ctx, req)
	if err != nil {
		c.log.Error("collector get repo failed", "error", err)
		return nil, err
	}

	return &domain.Repository{
		FullName:        resp.FullName,
		Description:     resp.Description,
		StargazersCount: int(resp.StargazersCount),
		ForksCount:      int(resp.ForksCount),
		CreatedAt:       resp.CreatedAt.AsTime(),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
