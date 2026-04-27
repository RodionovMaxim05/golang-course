package github

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"repo-stat/collector/internal/domain"
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

type RepoInfo struct {
	FullName        string    `json:"full_name"`
	Description     string    `json:"description"`
	StargazersCount int       `json:"stargazers_count"`
	ForksCount      int       `json:"forks_count"`
	CreatedAt       time.Time `json:"created_at"`
}

func (gc *GitHubClient) GetRepo(ctx context.Context, owner, name string) (domain.Repository, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, name)
	gc.log.Debug("fetching repository from github", "owner", owner, "repo", name, "url", url)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		gc.log.Error("failed to create github request", "error", err, "owner", owner, "repo", name)
		return domain.Repository{}, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := gc.httpClient.Do(req)
	if err != nil {
		gc.log.Error("failed to execute github request", "error", err, "owner", owner, "repo", name)
		return domain.Repository{}, fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			gc.log.Error("failed to close response body", "error", err)
		}
	}()

	if err := gc.checkStatus(resp); err != nil {
		return domain.Repository{}, err
	}

	var info RepoInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		gc.log.Error("failed to parse github response", "error", err, "owner", owner, "repo", name)
		return domain.Repository{}, fmt.Errorf("failed to parse JSON: %w", err)
	}

	gc.log.Debug("repository fetched successfully", "owner", owner, "repo", name, "stars", info.StargazersCount, "created_at", info.CreatedAt)

	return domain.Repository{
		FullName:        info.FullName,
		Description:     info.Description,
		StargazersCount: info.StargazersCount,
		ForksCount:      info.ForksCount,
		CreatedAt:       info.CreatedAt,
	}, nil
}

func (gc *GitHubClient) checkStatus(resp *http.Response) error {
	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		gc.log.Warn("repository not found in github", "status_code", resp.StatusCode)
		return domain.ErrNotFound
	case http.StatusForbidden:
		gc.log.Warn("github rate limit exceeded", "status_code", resp.StatusCode)
		return domain.ErrRateLimited
	default:
		gc.log.Error("github api error", "status_code", resp.StatusCode, "status", resp.Status)
		return fmt.Errorf("api error: status %d (%s)", resp.StatusCode, resp.Status)
	}
}
