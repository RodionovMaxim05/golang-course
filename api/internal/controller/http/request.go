package http

import (
	"fmt"
	"net/url"
	"strings"
)

// parseGitHubURL extracts the owner and repository name from a GitHub
// repository URL (e.g. "https://github.com/golang/go"). Returns an error
// if the URL is not a valid https://github.com/{owner}/{repo} URL.
func parseGitHubURL(rawURL string) (owner, repo string, err error) {
	u, err := url.Parse(rawURL)
	if err != nil || u.Scheme != "https" || u.Host != "github.com" {
		return "", "", fmt.Errorf("unsupported host or invalid url")
	}

	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid path: %s", u.Path)
	}

	return parts[0], parts[1], nil
}
