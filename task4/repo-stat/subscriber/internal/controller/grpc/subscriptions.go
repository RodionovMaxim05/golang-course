package grpc

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	subscriberpb "repo-stat/proto/subscriber"
)

func (s *Server) Subscribe(ctx context.Context, req *subscriberpb.SubscribeRequest) (*subscriberpb.SubscribeResponse, error) {
	s.log.Debug("subscriber subscribe request received")

	subscription, err := s.subscribe.Execute(ctx, req.Subscription.Owner, req.Subscription.Repo)
	if err != nil {
		return nil, grpcError(err)
	}

	return &subscriberpb.SubscribeResponse{
		Subscription: &subscriberpb.SubscriptionResponse{
			Owner:     subscription.Owner,
			Repo:      subscription.Repo,
			CreatedAt: timestamppb.New(subscription.CreatedAt),
		},
	}, nil
}

func (s *Server) Unsubscribe(ctx context.Context, req *subscriberpb.UnsubscribeRequest) (*subscriberpb.UnsubscribeResponse, error) {
	s.log.Debug("subscriber unsubscribe request received")

	err := s.unsubscribe.Execute(ctx, req.Subscription.Owner, req.Subscription.Repo)
	if err != nil {
		return nil, grpcError(err)
	}

	return &subscriberpb.UnsubscribeResponse{}, nil
}

func (s *Server) GetSubscriptions(ctx context.Context, _ *subscriberpb.GetSubsRequest) (*subscriberpb.GetSubsResponse, error) {
	s.log.Debug("subscriber get subscriptions request received")

	subscriptions, err := s.getSubscriptions.Execute(ctx)
	if err != nil {
		return nil, grpcError(err)
	}

	result := make([]*subscriberpb.SubscriptionResponse, 0, len(subscriptions))
	for _, subscription := range subscriptions {
		result = append(result, &subscriberpb.SubscriptionResponse{
			Owner:     subscription.Owner,
			Repo:      subscription.Repo,
			CreatedAt: timestamppb.New(subscription.CreatedAt),
		})
	}

	return &subscriberpb.GetSubsResponse{Subscriptions: result}, nil
}
