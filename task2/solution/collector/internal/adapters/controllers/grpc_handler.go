package controllers

import (
	"collector/internal/domain"
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"

	pb "api/gen"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type RepoUsecases interface {
	Execute(owner, name string) (domain.Repository, error)
}

type RepoHandler struct {
	pb.UnsafeRepoServiceServer
	repoUsecases RepoUsecases
}

func NewRepoHandler(repoUsecases RepoUsecases) *RepoHandler {
	return &RepoHandler{repoUsecases: repoUsecases}
}

func (rh *RepoHandler) GetRepo(ctx context.Context, req *pb.GetRepoRequest) (*pb.GetRepoResponse, error) {
	owner, name, err := parseGitHubURL(req.Url)
	if err != nil {
		return nil,
			status.Error(
				codes.InvalidArgument,
				fmt.Sprintf("invalid github url %q: expected \"https://github.com/owner/repo\"", req.Url),
			)
	}

	repo, err := rh.repoUsecases.Execute(owner, name)
	if err != nil {
		log.Printf("execute error for %q: %v", req.Url, err)
		return nil, mapError(err)
	}

	return &pb.GetRepoResponse{
		Name:            repo.Name,
		Description:     repo.Description,
		StargazersCount: int32(repo.StargazersCount),
		ForksCount:      int32(repo.ForksCount),
		CreatedAt:       timestamppb.New(repo.CreatedAt),
	}, nil
}

func mapError(err error) error {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, domain.ErrRateLimited):
		return status.Error(codes.ResourceExhausted, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}

func parseGitHubURL(rawURL string) (owner, name string, err error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", "", err
	}

	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("expected /owner/repo, got %s", u.Path)
	}

	return parts[0], parts[1], nil
}
