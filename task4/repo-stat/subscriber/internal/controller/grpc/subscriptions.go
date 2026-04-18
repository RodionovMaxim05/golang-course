package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	subscriberpb "repo-stat/proto/subscriber"
	"repo-stat/subscriber/internal/domain"
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

func grpcError(err error) error {
	switch {
	case errors.Is(err, domain.ErrInvalidArgument):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, domain.ErrRepoNotFound), errors.Is(err, domain.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, domain.ErrAlreadySubscribed):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, domain.ErrRateLimited):
		return status.Error(codes.ResourceExhausted, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
