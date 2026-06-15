package kafkaController

import (
	"context"
	"errors"
	"log/slog"

	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"

	"repo-stat/collector/internal/domain"
	collectorpb "repo-stat/proto/collector"
)

type GetRepoUsecase interface {
	Execute(ctx context.Context, owner, name string) (domain.Repository, error)
}

type ResultProducer interface {
	SendRepoResult(ctx context.Context, repo domain.Repository) error
}

type RepoWorker struct {
	log         *slog.Logger
	reader      *kafka.Reader
	repoUsecase GetRepoUsecase
	producer    ResultProducer
}

func NewRepoWorker(
	reader *kafka.Reader,
	repoUsecase GetRepoUsecase,
	producer ResultProducer,
	log *slog.Logger,
) *RepoWorker {
	return &RepoWorker{
		log:         log,
		reader:      reader,
		repoUsecase: repoUsecase,
		producer:    producer,
	}
}

func (rw *RepoWorker) Start(ctx context.Context) {
	rw.log.Info("Collector Kafka listens for tasks...")

	for {
		message, err := rw.reader.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				rw.log.Info("kafka worker stopped by context signal")
				return
			}

			rw.log.Error("error reading task from Kafka", "error", err)
			continue
		}

		rw.log.Debug("get message from Kafka")

		var req collectorpb.GetRepoRequest
		if err := proto.Unmarshal(message.Value, &req); err != nil {
			rw.log.Error("error deserializing incoming task", "err", err)
			continue
		}

		rw.log.Debug("a task for collecting data from Kafka has been received", "owner", req.Name, "repo", req.Repo)

		repo, err := rw.repoUsecase.Execute(ctx, req.Name, req.Repo)
		if err != nil {
			rw.log.Error("failed to collect data from GitHub API",
				"owner", req.Name,
				"repo", req.Repo,
				"err", err,
			)
			continue
		}

		err = rw.producer.SendRepoResult(ctx, repo)
		if err != nil {
			rw.log.Error("failed to send collection result to Kafka",
				"repo", repo.FullName,
				"err", err,
			)
			continue
		}

		rw.log.Info("the result has been successfully collected and sent to Kafka", "repo", repo.FullName)
	}
}
