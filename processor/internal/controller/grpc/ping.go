package grpc

import (
	"context"
	"log/slog"

	"repo-watcher/processor/internal/usecase"
	processorpb "repo-watcher/proto/gen/go/processor/v1"
)

type Server struct {
	processorpb.UnimplementedProcessorServer
	log                  *slog.Logger
	ping                 *usecase.Ping
	getRepo              *usecase.GetRepo
	getSubscriptionsInfo *usecase.GetSubscriptionsInfo
}

func NewServer(log *slog.Logger, ping *usecase.Ping, getRepo *usecase.GetRepo, getSubscriptionsInfo *usecase.GetSubscriptionsInfo) *Server {
	return &Server{
		log:                  log,
		ping:                 ping,
		getRepo:              getRepo,
		getSubscriptionsInfo: getSubscriptionsInfo,
	}
}

// Ping is a liveness check endpoint, returning a static reply from the
// underlying Ping use case.
func (s *Server) Ping(ctx context.Context, _ *processorpb.PingRequest) (*processorpb.PingResponse, error) {
	s.log.Debug("processor ping request received")

	return &processorpb.PingResponse{
		Reply: s.ping.Execute(ctx),
	}, nil
}
