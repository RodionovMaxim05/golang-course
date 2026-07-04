package grpc

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	subscriberpb "repo-watcher/proto/gen/go/subscriber/v1"
)

// Subscribe creates a new subscription for the requested repository,
// validating its existence on GitHub via the Subscribe use case.
func (s *Server) Subscribe(ctx context.Context, req *subscriberpb.SubscribeRequest) (*subscriberpb.SubscribeResponse, error) {
	s.log.Debug("subscriber subscribe request received")

	subscription, err := s.subscribe.Execute(ctx, req.Subscription.Owner, req.Subscription.Repo)
	if err != nil {
		s.log.Error("subscribe usecase failed", "owner", req.Subscription.Owner, "repo", req.Subscription.Repo, "error", err)
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

// Unsubscribe removes an existing subscription for the requested repository.
func (s *Server) Unsubscribe(ctx context.Context, req *subscriberpb.UnsubscribeRequest) (*subscriberpb.UnsubscribeResponse, error) {
	s.log.Debug("subscriber unsubscribe request received")

	err := s.unsubscribe.Execute(ctx, req.Subscription.Owner, req.Subscription.Repo)
	if err != nil {
		s.log.Error("unsubscribe usecase failed", "owner", req.Subscription.Owner, "repo", req.Subscription.Repo, "error", err)
		return nil, grpcError(err)
	}

	return &subscriberpb.UnsubscribeResponse{}, nil
}

// GetSubscriptions returns the full list of currently active subscriptions.
func (s *Server) GetSubscriptions(ctx context.Context, _ *subscriberpb.GetSubsRequest) (*subscriberpb.GetSubsResponse, error) {
	s.log.Debug("subscriber get subscriptions request received")

	subscriptions, err := s.getSubscriptions.Execute(ctx)
	if err != nil {
		s.log.Error("get subscriptions usecase failed", "error", err)
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
