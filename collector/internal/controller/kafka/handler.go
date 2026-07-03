package kafkacontroller

import (
	"context"
	"errors"
	"log/slog"

	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"

	"repo-watcher/collector/internal/domain"
	collectorpb "repo-watcher/proto/gen/go/collector/v1"
)

type GetRepoUsecase interface {
	Execute(ctx context.Context, owner, repo string) (domain.Repository, error)
}

type ResultProducer interface {
	SendRepoResult(ctx context.Context, repo *domain.Repository, workErr error, fullname string) error
}

type RepoWorker struct {
	log         *slog.Logger
	reader      *kafka.Reader
	repoUsecase GetRepoUsecase
	producer    ResultProducer
}

func NewRepoWorker(
	log *slog.Logger,
	reader *kafka.Reader,
	repoUsecase GetRepoUsecase,
	producer ResultProducer,
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

		var req collectorpb.CollectRepoCmd
		if err := proto.Unmarshal(message.Value, &req); err != nil {
			rw.log.Error("error deserializing incoming task", "error", err)
			continue
		}

		rw.log.Debug("task received", "owner", req.Owner, "repo", req.Repo)

		repo, err := rw.repoUsecase.Execute(ctx, req.Owner, req.Repo)
		if err != nil {
			rw.log.Error("failed to collect data from GitHub API",
				"owner", req.Owner,
				"repo", req.Repo,
				"err", err,
			)

			rw.sendResultToKafka(ctx, nil, err, req.Owner+"/"+req.Repo)
			continue
		}

		rw.sendResultToKafka(ctx, &repo, nil, repo.FullName)
	}
}

func (rw *RepoWorker) sendResultToKafka(ctx context.Context, repo *domain.Repository, workErr error, repoIdentifier string) {
	err := rw.producer.SendRepoResult(ctx, repo, workErr, repoIdentifier)
	if err != nil {
		rw.log.Error("failed to send result or error to Kafka",
			"repo", repoIdentifier,
			"error", err,
		)
		return
	}

	if workErr != nil {
		rw.log.Info("error status successfully sent to Kafka", "repo", repoIdentifier)
	} else {
		rw.log.Info("result successfully processed and sent to Kafka", "repo", repoIdentifier)
	}
}
