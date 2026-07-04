package kafka

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"repo-watcher/collector/config"
	"repo-watcher/collector/internal/domain"
	commonpb "repo-watcher/proto/gen/go/common/v1"
)

type ResultProducerAdapter struct {
	log    *slog.Logger
	writer *kafka.Writer
}

func NewResultProducerAdapter(cfg config.Kafka, log *slog.Logger) *ResultProducerAdapter {
	return &ResultProducerAdapter{
		log: log,
		writer: &kafka.Writer{
			Addr:         kafka.TCP(cfg.Address),
			Topic:        cfg.ProducerTopic,
			Balancer:     &kafka.LeastBytes{},
			RequiredAcks: kafka.RequireOne,
		},
	}
}

// SendSuccess publishes successfully collected repository data to Kafka.
func (pc *ResultProducerAdapter) SendSuccess(ctx context.Context, repo domain.Repository) error {
	event := &commonpb.GetRepoResultEvent{
		FullName: repo.FullName,
		Result: &commonpb.GetRepoResultEvent_Success{
			Success: &commonpb.CollectRepoSuccess{
				Description:     repo.Description,
				StargazersCount: int32(repo.StargazersCount),
				ForksCount:      int32(repo.ForksCount),
				CreatedAt:       timestamppb.New(repo.CreatedAt),
			},
		},
	}

	return pc.publish(ctx, event)
}

// SendFailure publishes a repository collection failure to Kafka, mapping
// the domain error to a wire-level error code.
func (pc *ResultProducerAdapter) SendFailure(ctx context.Context, fullName string, cause error) error {
	event := &commonpb.GetRepoResultEvent{
		FullName: fullName,
		Result: &commonpb.GetRepoResultEvent_Error{
			Error: &commonpb.ErrorResponse{
				Code:    mapDomainErrorToCode(cause),
				Message: cause.Error(),
			},
		},
	}

	return pc.publish(ctx, event)
}

// Publish marshals and writes a single result event to Kafka.
func (pc *ResultProducerAdapter) publish(ctx context.Context, event *commonpb.GetRepoResultEvent) error {
	eventBytes, err := proto.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal repo result event: %w", err)
	}

	if err := pc.writer.WriteMessages(ctx, kafka.Message{Value: eventBytes}); err != nil {
		return fmt.Errorf("write repo result event to kafka: %w", err)
	}

	return nil
}

func (pc *ResultProducerAdapter) Close() error {
	return pc.writer.Close()
}
