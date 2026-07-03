package grpc

import (
	"context"

	processorpb "repo-watcher/proto/gen/go/processor/v1"

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
		if repo.Status != "SUCCESS" {
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
