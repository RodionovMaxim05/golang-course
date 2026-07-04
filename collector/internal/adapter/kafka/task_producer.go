package kafka

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"

	"repo-watcher/collector/config"
	collectorpb "repo-watcher/proto/gen/go/collector/v1"
)

type TaskProducerAdapter struct {
	log    *slog.Logger
	writer *kafka.Writer
}

func NewTaskProducerAdapter(cfg config.Kafka, log *slog.Logger) *TaskProducerAdapter {
	return &TaskProducerAdapter{
		log: log,
		writer: &kafka.Writer{
			Addr:         kafka.TCP(cfg.Address),
			Topic:        cfg.ConsumerTopic,
			Balancer:     &kafka.LeastBytes{},
			RequiredAcks: kafka.RequireOne,
		},
	}
}

// SendCollectionTask publishes a request to collect data for the given owner/repo.
func (tpa *TaskProducerAdapter) SendCollectionTask(ctx context.Context, owner, repo string) error {
	getRepoRequest := &collectorpb.CollectRepoCmd{
		Owner: owner,
		Repo:  repo,
	}

	requestBytes, err := proto.Marshal(getRepoRequest)
	if err != nil {
		return fmt.Errorf("marshal collection task: %w", err)
	}

	err = tpa.writer.WriteMessages(ctx, kafka.Message{Value: requestBytes})
	if err != nil {
		return fmt.Errorf("write collection task to kafka: %w", err)
	}

	tpa.log.Debug("task producer send request", "owner", owner, "repo", repo)

	return nil
}

func (tpa *TaskProducerAdapter) Close() error {
	return tpa.writer.Close()
}
