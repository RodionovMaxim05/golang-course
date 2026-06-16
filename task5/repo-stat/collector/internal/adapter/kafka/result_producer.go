package kafka

import (
	"context"
	"log/slog"

	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"repo-stat/collector/config"
	"repo-stat/collector/internal/domain"
	collectorpb "repo-stat/proto/processor"
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

func (pc *ResultProducerAdapter) SendRepoResult(ctx context.Context, repo *domain.Repository, workerErr error, fullname string) error {
	event := &collectorpb.GetRepoResultEvent{FullName: fullname}

	if workerErr != nil {
		event.Result = &collectorpb.GetRepoResultEvent_Error{
			Error: &collectorpb.ErrorResponse{
				Code:    mapDomainErrorToCode(workerErr),
				Message: workerErr.Error(),
			},
		}
	} else {
		event.Result = &collectorpb.GetRepoResultEvent_Success{
			Success: &collectorpb.GetRepoResponse{
				Description:     repo.Description,
				StargazersCount: int32(repo.StargazersCount),
				ForksCount:      int32(repo.ForksCount),
				CreatedAt:       timestamppb.New(repo.CreatedAt),
			},
		}
	}

	responseBytes, err := proto.Marshal(event)
	if err != nil {
		pc.log.Error("failed to Marshal get repo response", "error", err)
		return err
	}

	err = pc.writer.WriteMessages(ctx, kafka.Message{Value: responseBytes})
	if err != nil {
		pc.log.Error("failed to write messages in Kafka", "error", err)
		return err
	}

	return nil
}

func (pc *ResultProducerAdapter) Close() error {
	return pc.writer.Close()
}
