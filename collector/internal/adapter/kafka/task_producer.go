package kafka

import (
	"context"
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

func (tpa *TaskProducerAdapter) SendCollectionTask(ctx context.Context, owner, repo string) error {
	getRepoRequest := &collectorpb.CollectRepoCmd{
		Owner: owner,
		Repo:  repo,
	}

	tpa.log.Debug("task producer get request")

	requestBytes, err := proto.Marshal(getRepoRequest)
	if err != nil {
		tpa.log.Error("failed to Marshal get repo request", "error", err)
		return err
	}

	err = tpa.writer.WriteMessages(ctx, kafka.Message{Value: requestBytes})
	if err != nil {
		tpa.log.Error("failed to write messages in Kafka", "error", err)
		return err
	}

	tpa.log.Debug("task producer send request", "owner", owner, "repo", repo)

	return nil
}

func (tpa *TaskProducerAdapter) Close() error {
	return tpa.writer.Close()
}
