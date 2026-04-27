package grpc

import (
	"context"

	processorpb "repo-stat/proto/processor"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) GetSubscriptionsInfo(ctx context.Context, req *processorpb.GetSubsInfoRequest) (*processorpb.GetSubsInfoResponse, error) {
	s.log.Debug("processor get subscriptions info request received")

	resp, err := s.getSubscriptionsInfo.Execute(ctx)
	if err != nil {
		return nil, err
	}

	repositories := make([]*processorpb.GetRepoResponse, 0, len(resp))
	for _, repo := range resp {
		repositories = append(repositories, &processorpb.GetRepoResponse{
			FullName:        repo.FullName,
			Description:     repo.Description,
			StargazersCount: int32(repo.StargazersCount),
			ForksCount:      int32(repo.ForksCount),
			CreatedAt:       timestamppb.New(repo.CreatedAt),
		})
	}

	return &processorpb.GetSubsInfoResponse{Repositories: repositories}, nil
}
