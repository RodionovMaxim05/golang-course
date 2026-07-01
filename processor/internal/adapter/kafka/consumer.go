package kafka

import (
	"context"
	"errors"
	"log/slog"

	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"

	"repo-watcher/processor/config"
	"repo-watcher/processor/internal/domain"
	commonpb "repo-watcher/proto/common"
)

type ConsumerClient struct {
	log     *slog.Logger
	storage domain.DataStorage
	reader  *kafka.Reader
}

func NewConsumerClient(cfg config.Kafka, storage domain.DataStorage, log *slog.Logger) *ConsumerClient {
	return &ConsumerClient{
		log:     log,
		storage: storage,
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:     []string{cfg.Address},
			Topic:       cfg.ConsumerTopic,
			GroupID:     cfg.GroupID,
			StartOffset: kafka.FirstOffset,
		}),
	}
}

func (cc *ConsumerClient) StartConsumer(ctx context.Context) {
	cc.log.Info("Processor Kafka consumer started, listening for results...")

	for {
		message, err := cc.reader.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				cc.log.Info("Kafka consumer stopped by context signal")
				return
			}

			cc.log.Error("error reading message from Kafka", "error", err)
			continue
		}

		cc.log.Debug("Received raw message from Kafka", "bytes_len", len(message.Value))

		var event commonpb.GetRepoResultEvent
		if err := proto.Unmarshal(message.Value, &event); err != nil {
			cc.log.Error("message deserialization error",
				"error", err,
				"raw_bytes", message.Value,
			)
			continue
		}

		cc.log.Debug("Successfully unmarshaled event", "repo", event.FullName, "has_result", event.Result != nil)

		switch res := event.Result.(type) {
		case *commonpb.GetRepoResultEvent_Error:
			cc.log.Warn("collector failed to process repository, saving error status to DB",
				"repo", event.FullName,
				"code", res.Error.GetCode(),
				"message", res.Error.GetMessage(),
			)

			repoEntity := domain.Repository{
				FullName:  event.FullName,
				Status:    "ERROR",
				ErrorCode: res.Error.GetCode(),
			}

			err = cc.storage.UpdateRepoStatus(ctx, &repoEntity)
			if err != nil {
				cc.log.Error("failed to save error status to the database",
					"repo", event.FullName,
					"error", err,
				)
			}
			continue

		case *commonpb.GetRepoResultEvent_Success:
			successData := res.Success

			repoEntity := domain.Repository{
				FullName:        event.FullName,
				Description:     successData.Description,
				StargazersCount: int(successData.StargazersCount),
				ForksCount:      int(successData.ForksCount),
				CreatedAt:       successData.GetCreatedAt().AsTime(),
				Status:          "SUCCESS",
				ErrorCode:       "",
			}

			err = cc.storage.InsertRepo(ctx, &repoEntity)
			if err != nil {
				cc.log.Error("failed to save data from Kafka to the database",
					"repo", repoEntity.FullName,
					"error", err,
				)

				continue
			}

			cc.log.Info("repository data has been successfully updated to the database from Kafka",
				"repo", repoEntity.FullName,
			)

		case nil:
			cc.log.Error("received event with empty result field (oneof is nil)")
			continue
		}
	}
}

func (c *ConsumerClient) Close() error {
	return c.reader.Close()
}
