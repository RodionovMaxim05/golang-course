package grpc

import (
	"context"

	processorpb "repo-stat/proto/processor"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) GetRepo(ctx context.Context, req *processorpb.GetRepoRequest) (*processorpb.GetRepoResponse, error) {
	s.log.Debug("processor get repo request received", "name", req.Name, "repo", req.Repo)

	resp, err := s.getRepo.Execute(ctx, req.Name, req.Repo)
	if err != nil {
		return nil, err
	}

	return &processorpb.GetRepoResponse{
		FullName:        resp.FullName,
		Description:     resp.Description,
		StargazersCount: int32(resp.StargazersCount),
		ForksCount:      int32(resp.ForksCount),
		CreatedAt:       timestamppb.New(resp.CreatedAt),
	}, nil
}
