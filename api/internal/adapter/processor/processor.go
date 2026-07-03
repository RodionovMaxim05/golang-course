package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	grpcAdapter "repo-watcher/api/internal/adapter/grpc"
	"repo-watcher/api/internal/domain"
	processorpb "repo-watcher/proto/gen/go/processor/v1"
)

type CacheClient interface {
	Get(ctx context.Context, key string) ([]byte, bool, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
}

type Client struct {
	log      *slog.Logger
	conn     *grpc.ClientConn
	pb       processorpb.ProcessorClient
	cache    CacheClient
	cacheTTL time.Duration
}

func NewClient(address string, cache CacheClient, ttl int, log *slog.Logger) (*Client, error) {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		log:      log,
		conn:     conn,
		pb:       processorpb.NewProcessorClient(conn),
		cache:    cache,
		cacheTTL: time.Duration(ttl) * time.Second,
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

func (c *Client) GetRepo(ctx context.Context, owner, repo string) (domain.Repository, error) {
	cacheKey := fmt.Sprintf("processor:repo:%s:%s", owner, repo)

	if cacheBytes, exists, err := c.cache.Get(ctx, cacheKey); err == nil && exists {
		var repository domain.Repository
		if err := json.Unmarshal(cacheBytes, &repository); err == nil {
			c.log.Debug("cache hit for repo", "key", cacheKey)
			return repository, nil
		}
		c.log.Warn("corrupted cache data for repo, bypass to gRPC", "key", cacheKey, "error", err)
	} else if err != nil {
		c.log.Error("cache error during GetRepo", "error", err)
	}

	req := &processorpb.GetRepoRequest{
		Owner: owner,
		Repo:  repo,
	}

	resp, err := c.pb.GetRepo(ctx, req)
	if err != nil {
		c.log.Error("processor get repo failed", "error", err)
		return domain.Repository{}, grpcAdapter.ErrToDomain(err)
	}

	switch resp.Status {
	case processorpb.GetRepoResponse_STATUS_PENDING:
		return domain.Repository{}, domain.ErrPending

	case processorpb.GetRepoResponse_STATUS_SUCCESS:
		result := domain.Repository{
			FullName:    resp.FullName,
			Description: resp.Description,
			Stargazers:  int(resp.StargazersCount),
			Forks:       int(resp.ForksCount),
			CreatedAt:   resp.CreatedAt.AsTime(),
		}

		if jsonBytes, err := json.Marshal(result); err == nil {
			if err := c.cache.Set(ctx, cacheKey, jsonBytes, c.cacheTTL); err != nil {
				c.log.Error("failed to update repo cache", "key", cacheKey, "error", err)
			}
		}

		return result, nil

	default:
		c.log.Error("unexpected repository status from processor", "status", resp.Status)
		return domain.Repository{}, domain.ErrInternal
	}
}

func (c *Client) GetSubscriptionsInfo(ctx context.Context) ([]domain.Repository, error) {
	cacheKey := "processor:subscriptions:all"

	if cacheBytes, exists, err := c.cache.Get(ctx, cacheKey); err == nil && exists {
		var cachedSubs []domain.Repository
		if err := json.Unmarshal(cacheBytes, &cachedSubs); err == nil {
			c.log.Debug("cache hit for subscriptions info")
			return cachedSubs, nil
		}
		c.log.Warn("corrupted cache data for subscriptions, bypass to gRPC", "error", err)
	} else if err != nil {
		c.log.Error("cache error during GetSubscriptionsInfo", "error", err)
	}

	req := &processorpb.GetSubsInfoRequest{}

	resp, err := c.pb.GetSubscriptionsInfo(ctx, req)
	if err != nil {
		c.log.Error("processor get subscriptions info failed", "error", err)
		return nil, grpcAdapter.ErrToDomain(err)
	}

	result := make([]domain.Repository, 0, len(resp.Repositories))
	for _, item := range resp.Repositories {
		result = append(result, domain.Repository{
			FullName:    item.FullName,
			Description: item.Description,
			Stargazers:  int(item.StargazersCount),
			Forks:       int(item.ForksCount),
			CreatedAt:   item.CreatedAt.AsTime(),
		})
	}

	if jsonBytes, err := json.Marshal(result); err == nil {
		if err := c.cache.Set(ctx, cacheKey, jsonBytes, c.cacheTTL); err != nil {
			c.log.Error("failed to update subscriptions cache", "key", cacheKey, "error", err)
		}
	}

	return result, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
