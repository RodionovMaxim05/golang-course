package kafka

import (
	"context"
	"log/slog"

	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"repo-stat/collector/config"
	"repo-stat/collector/internal/domain"
	collectorpb "repo-stat/proto/collector"
)

type ProducerClient struct {
	log    *slog.Logger
	writer *kafka.Writer
}

func NewProducerClient(cfg config.Kafka, log *slog.Logger) *ProducerClient {
	return &ProducerClient{
		log: log,
		writer: &kafka.Writer{
			Addr:         kafka.TCP(cfg.Address),
			Topic:        cfg.ProducerTopic,
			Balancer:     &kafka.LeastBytes{},
			RequiredAcks: kafka.RequireAll,
		},
	}
}

func (pc *ProducerClient) SendRepoResult(ctx context.Context, repo domain.Repository) error {
	getRepoResponse := &collectorpb.GetRepoResponse{
		FullName:        repo.FullName,
		Description:     repo.Description,
		StargazersCount: int32(repo.StargazersCount),
		ForksCount:      int32(repo.ForksCount),
		CreatedAt:       timestamppb.New(repo.CreatedAt),
	}

	responseBytes, err := proto.Marshal(getRepoResponse)
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

func (pc *ProducerClient) Close() error {
	return pc.writer.Close()
}
