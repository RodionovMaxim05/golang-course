package grpc

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"repo-watcher/processor/internal/domain"
	processorpb "repo-watcher/proto/gen/go/processor/v1"
)

func (s *Server) GetSubscriptionsInfo(ctx context.Context, req *processorpb.GetSubsInfoRequest) (*processorpb.GetSubsInfoResponse, error) {
	s.log.Debug("processor get subscriptions info request received")

	resp, err := s.getSubscriptionsInfo.Execute(ctx)
	if err != nil {
		return nil, err
	}

	repositories := make([]*processorpb.GetRepoResponse, 0, len(resp))
	for _, repo := range resp {
		if repo.Status != domain.StatusSuccess {
			continue
		}

		repositories = append(repositories, &processorpb.GetRepoResponse{
			Status:          processorpb.GetRepoResponse_STATUS_SUCCESS,
			FullName:        repo.FullName,
			Description:     repo.Description,
			StargazersCount: int32(repo.StargazersCount),
			ForksCount:      int32(repo.ForksCount),
			CreatedAt:       timestamppb.New(repo.CreatedAt),
		})
	}

	return &processorpb.GetSubsInfoResponse{Repositories: repositories}, nil
}
