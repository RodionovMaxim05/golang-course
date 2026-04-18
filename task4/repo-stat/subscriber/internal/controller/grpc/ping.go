package grpc

import (
	"context"
	"log/slog"

	subscriberpb "repo-stat/proto/subscriber"
	"repo-stat/subscriber/internal/usecase"
)

type Server struct {
	subscriberpb.UnimplementedSubscriberServer
	log              *slog.Logger
	ping             *usecase.Ping
	subscribe        *usecase.Subscribe
	unsubscribe      *usecase.Unsubscribe
	getSubscriptions *usecase.GetSubscriptions
}

func NewServer(log *slog.Logger, ping *usecase.Ping, subscribe *usecase.Subscribe, unsubscribe *usecase.Unsubscribe, getSubscriptions *usecase.GetSubscriptions) *Server {
	return &Server{
		log:              log,
		ping:             ping,
		subscribe:        subscribe,
		unsubscribe:      unsubscribe,
		getSubscriptions: getSubscriptions,
	}
}

func (s *Server) Ping(ctx context.Context, _ *subscriberpb.PingRequest) (*subscriberpb.PingResponse, error) {
	s.log.Debug("subscriber ping request received")

	return &subscriberpb.PingResponse{
		Reply: s.ping.Execute(ctx),
	}, nil
}
