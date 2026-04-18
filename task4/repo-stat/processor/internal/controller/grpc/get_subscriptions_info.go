package grpc

import (
	"context"

	processorpb "repo-stat/proto/processor"
)

func (s *Server) GetSubscriptionsInfo(ctx context.Context, req *processorpb.GetSubsInfoRequest) (*processorpb.GetSubsInfoResponse, error) {
	s.log.Debug("processor get subscriptions info request received")

	return s.getSubscriptionsInfo.Execute(ctx)
}
