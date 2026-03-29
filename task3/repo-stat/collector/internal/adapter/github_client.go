package adapter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"repo-stat/collector/internal/domain"
)

type Config struct {
	Timeout time.Duration `yaml:"timeout" env:"TIMEOUT" env-default:"5s"`
}

type GitHubClient struct {
	httpClient *http.Client
}

func NewGitHubClient(cfg Config) *GitHubClient {
	return &GitHubClient{
		httpClient: &http.Client{Timeout: cfg.Timeout},
	}
}

type RepoInfo struct {
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	StargazersCount int       `json:"stargazers_count"`
	ForksCount      int       `json:"forks_count"`
	CreatedAt       time.Time `json:"created_at"`
}

func (gc *GitHubClient) GetRepo(owner, name string) (domain.Repository, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, name)

	resp, err := gc.httpClient.Get(url)
	if err != nil {
		return domain.Repository{}, fmt.Errorf("request error: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to close response body: %s\n", err)
		}
	}()

	switch resp.StatusCode {
	case http.StatusOK:
		// ok
	case http.StatusNotFound:
		return domain.Repository{}, domain.ErrNotFound
	case http.StatusForbidden:
		return domain.Repository{}, domain.ErrRateLimited
	default:
		return domain.Repository{}, fmt.Errorf("api error: %s", resp.Status)
	}

	var info RepoInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return domain.Repository{}, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return domain.Repository{
		Name:            info.Name,
		Description:     info.Description,
		StargazersCount: info.StargazersCount,
		ForksCount:      info.ForksCount,
		CreatedAt:       info.CreatedAt,
	}, nil
}
