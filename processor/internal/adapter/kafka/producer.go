package kafka

import (
	"context"
	"log/slog"

	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"

	"repo-watcher/processor/config"
	processorpb "repo-watcher/proto/gen/go/processor/v1"
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
			RequiredAcks: kafka.RequireOne,
		},
	}
}

func (pc *ProducerClient) SendRepoRequest(ctx context.Context, owner, repo string) error {
	getRepoRequest := &processorpb.GetRepoRequest{
		Owner: owner,
		Repo:  repo,
	}

	pc.log.Debug("producer get request")

	requestBytes, err := proto.Marshal(getRepoRequest)
	if err != nil {
		pc.log.Error("failed to Marshal get repo request", "error", err)
		return err
	}

	err = pc.writer.WriteMessages(ctx, kafka.Message{Value: requestBytes})
	if err != nil {
		pc.log.Error("failed to write messages in Kafka", "error", err)
		return err
	}

	pc.log.Debug("message successfully written to Kafka")

	return nil
}

func (pc *ProducerClient) Close() error {
	return pc.writer.Close()
}
