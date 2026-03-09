package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

const timeout = 5 * time.Second

type repoInfo struct {
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	StargazersCount int       `json:"stargazers_count"`
	ForksCount      int       `json:"forks_count"`
	CreatedAt       time.Time `json:"created_at"`
}

func (rI repoInfo) print() {
	fmt.Println("=== Repository Information ===")
	fmt.Printf("%-12s %s\n", "Name:", rI.Name)
	fmt.Printf("%-12s %s\n", "Description:", func() string {
		if rI.Description == "" {
			return "No description"
		}
		return rI.Description
	}())
	fmt.Printf("%-12s %d\n", "Stars:", rI.StargazersCount)
	fmt.Printf("%-12s %d\n", "Forks:", rI.ForksCount)
	fmt.Printf("%-12s %s\n", "Created:", rI.CreatedAt.Format(time.DateTime))
}

func fetchRepoInfo(owner, repo string) (*repoInfo, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)

	client := &http.Client{Timeout: timeout}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
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
		return nil, fmt.Errorf("repository %s/%s not found", owner, repo)
	case http.StatusForbidden:
		return nil, fmt.Errorf("rate limit exceeded. try again later")
	default:
		return nil, fmt.Errorf("api error: %s", resp.Status)
	}

	var info repoInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &info, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage:\n\tgo run . <owner>/<repo>\n\n")
		fmt.Fprintf(os.Stderr, "Example:\n\tgo run . golang/go\n")
		os.Exit(1)
	}

	address := strings.Split(os.Args[1], "/")
	if len(address) != 2 {
		fmt.Fprintf(os.Stderr, "Error: argument must be in format <owner>/<repo>\n")
		os.Exit(1)
	}

	owner, repo := address[0], address[1]
	if owner == "" || repo == "" {
		fmt.Fprintf(os.Stderr, "Error: owner and repo cannot be empty\n")
		os.Exit(1)
	}

	info, err := fetchRepoInfo(owner, repo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	info.print()
}
