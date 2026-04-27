package github

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"repo-stat/subscriber/internal/domain"
)

type Config struct {
	Timeout time.Duration `yaml:"timeout" env:"TIMEOUT" env-default:"5s"`
}

type GitHubClient struct {
	log        *slog.Logger
	httpClient *http.Client
}

func NewGitHubClient(cfg Config, log *slog.Logger) *GitHubClient {
	return &GitHubClient{
		log:        log,
		httpClient: &http.Client{Timeout: cfg.Timeout},
	}
}

func (gc *GitHubClient) RepoExists(ctx context.Context, owner, repo string) (bool, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)
	gc.log.Debug("fetching repository from github", "owner", owner, "repo", repo, "url", url)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		gc.log.Error("failed to create github request", "error", err, "owner", owner, "repo", repo)
		return false, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := gc.httpClient.Do(req)
	if err != nil {
		gc.log.Error("failed to execute github request", "error", err, "owner", owner, "repo", repo)
		return false, fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			gc.log.Error("failed to close response body", "error", err)
		}
	}()

	switch resp.StatusCode {
	case http.StatusOK:
		return true, nil
	case http.StatusNotFound:
		gc.log.Warn("repository not found in github", "status_code", resp.StatusCode)
		return false, domain.ErrRepoNotFound
	case http.StatusForbidden:
		gc.log.Warn("github rate limit exceeded", "status_code", resp.StatusCode)
		return false, domain.ErrRateLimited
	default:
		gc.log.Error("github api error", "status_code", resp.StatusCode, "status", resp.Status)
		return false, fmt.Errorf("github api error: status %d", resp.StatusCode)
	}
}
