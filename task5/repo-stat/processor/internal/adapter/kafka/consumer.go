package kafka

import (
	"context"
	"errors"
	"log/slog"

	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"

	"repo-stat/processor/config"
	"repo-stat/processor/internal/domain"
	processorpb "repo-stat/proto/processor"
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
			Brokers: []string{cfg.Address},
			Topic:   cfg.ConsumerTopic,
			GroupID: cfg.GroupID,
		}),
	}
}

func (cc *ConsumerClient) StartConsumer(ctx context.Context) {
	for {
		message, err := cc.reader.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}

			cc.log.Error("error reading message from Kafka", "error", err)
			continue
		}

		var resultProto processorpb.GetRepoResponse
		if err := proto.Unmarshal(message.Value, &resultProto); err != nil {
			cc.log.Error("message deserialization error",
				"error", err,
				"raw_bytes", message.Value,
			)
			continue
		}

		repoEntity := domain.Repository{
			FullName:        resultProto.FullName,
			Description:     resultProto.Description,
			StargazersCount: int(resultProto.StargazersCount),
			ForksCount:      int(resultProto.ForksCount),
			CreatedAt:       resultProto.GetCreatedAt().AsTime(),
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
	}
}

func (c *ConsumerClient) Close() error {
	return c.reader.Close()
}
