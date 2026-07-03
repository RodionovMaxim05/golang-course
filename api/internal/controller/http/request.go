package http

import (
	"fmt"
	"net/url"
	"strings"
)

func parseGitHubURL(rawURL string) (owner, repo string, err error) {
	u, err := url.Parse(rawURL)
	if err != nil || u.Host != "github.com" {
		return "", "", fmt.Errorf("unsupported host or invalid url")
	}

	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid path: %s", u.Path)
	}

	return parts[0], parts[1], nil
}
