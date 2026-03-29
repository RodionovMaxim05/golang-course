package adapters

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	pb "api/gen"
	"gateway/internal/domain"
)

const timeout = 5 * time.Second

type Client struct {
	grpcClient pb.RepoServiceClient
}

func NewClient(addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to collector: %w", err)
	}

	return &Client{
		grpcClient: pb.NewRepoServiceClient(conn),
	}, nil
}

func (c *Client) GetRepo(url string) (domain.Repository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	r, err := c.grpcClient.GetRepo(ctx, &pb.GetRepoRequest{Url: url})
	if err != nil {
		st, _ := status.FromError(err)
		switch st.Code() {
		case codes.NotFound:
			return domain.Repository{}, domain.ErrNotFound
		case codes.ResourceExhausted:
			return domain.Repository{}, domain.ErrRateLimited
		case codes.InvalidArgument:
			return domain.Repository{}, fmt.Errorf("%w: %s", domain.ErrInvalidArgument, st.Message())
		default:
			return domain.Repository{}, fmt.Errorf("grpc call failed: %w", err)
		}
	}

	return domain.Repository{
		Name:            r.Name,
		Description:     r.Description,
		StargazersCount: int(r.StargazersCount),
		ForksCount:      int(r.ForksCount),
		CreatedAt:       r.CreatedAt.AsTime(),
	}, nil
}
