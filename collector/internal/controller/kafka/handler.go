package kafkacontroller

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"

	"repo-watcher/collector/internal/domain"
	collectorpb "repo-watcher/proto/gen/go/collector/v1"
)

const readErrorBackoff = time.Second

type GetRepoUsecase interface {
	Execute(ctx context.Context, owner, repo string) (domain.Repository, error)
}

type ResultProducer interface {
	SendSuccess(ctx context.Context, repo domain.Repository) error
	SendFailure(ctx context.Context, fullName string, cause error) error
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

// Start blocks, continuously reading collection tasks from Kafka and
// processing them until the context is cancelled. Transient read errors
// are logged and retried after a short backoff.
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
			time.Sleep(readErrorBackoff)
			continue
		}

		rw.processTask(ctx, message.Value)
	}
}

// processTask decodes a single collection task and executes it, publishing
// either the collected repository data or the failure reason back to
// Kafka. Decoding failures are logged and skipped without affecting
// subsequent messages.
func (rw *RepoWorker) processTask(ctx context.Context, rawMessage []byte) {
	var req collectorpb.CollectRepoCmd
	if err := proto.Unmarshal(rawMessage, &req); err != nil {
		rw.log.Error("error deserializing incoming task", "error", err)
		return
	}

	rw.log.Debug("task received", "owner", req.Owner, "repo", req.Repo)

	repo, err := rw.repoUsecase.Execute(ctx, req.Owner, req.Repo)
	if err != nil {
		fullName := req.Owner + "/" + req.Repo
		rw.log.Error("failed to collect data from GitHub API",
			"owner", req.Owner,
			"repo", req.Repo,
			"error", err,
		)

		if sendErr := rw.producer.SendFailure(ctx, fullName, err); sendErr != nil {
			rw.log.Error("failed to send failure result to Kafka", "repo", fullName, "error", sendErr)
			return
		}
		rw.log.Info("error status successfully sent to Kafka", "repo", fullName)
		return
	}

	if sendErr := rw.producer.SendSuccess(ctx, repo); sendErr != nil {
		rw.log.Error("failed to send result to Kafka", "repo", repo.FullName, "error", sendErr)
		return
	}
	rw.log.Info("result successfully processed and sent to Kafka", "repo", repo.FullName)
}
