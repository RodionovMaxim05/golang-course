package kafka

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"

	"repo-watcher/processor/config"
	"repo-watcher/processor/internal/domain"
	commonpb "repo-watcher/proto/gen/go/common/v1"
)

const readErrorBackoff = time.Second

type RepoUpdater interface {
	InsertRepo(ctx context.Context, repo *domain.Repository) error
	UpdateRepoStatus(ctx context.Context, repo *domain.Repository) error
}

type ConsumerClient struct {
	log     *slog.Logger
	storage RepoUpdater
	reader  *kafka.Reader
}

func NewConsumerClient(cfg config.Kafka, storage RepoUpdater, log *slog.Logger) *ConsumerClient {
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

// StartConsumer blocks, continuously reading messages from Kafka and
// processing them until the context is cancelled. Transient read errors
// are logged and retried after a short backoff.
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
			time.Sleep(readErrorBackoff)
			continue
		}

		cc.processMessage(ctx, message.Value)
	}
}

// processMessage decodes a single Kafka message into a GetRepoResultEvent
// and dispatches it to the appropriate handler based on its result type.
// Decoding failures and unknown/empty result types are logged and
// skipped without affecting subsequent messages.
func (cc *ConsumerClient) processMessage(ctx context.Context, rawMessage []byte) {
	cc.log.Debug("received raw message from Kafka", "bytes_len", len(rawMessage))

	var event commonpb.GetRepoResultEvent
	if err := proto.Unmarshal(rawMessage, &event); err != nil {
		cc.log.Error("message deserialization error",
			"error", err,
			"raw_bytes", rawMessage,
		)
		return
	}

	cc.log.Debug("successfully unmarshaled event", "repo", event.FullName, "has_result", event.Result != nil)

	switch res := event.Result.(type) {
	case *commonpb.GetRepoResultEvent_Error:
		cc.handleErrorEvent(ctx, event.FullName, res.Error)

	case *commonpb.GetRepoResultEvent_Success:
		cc.handleSuccessEvent(ctx, event.FullName, res.Success)

	case nil:
		cc.log.Error("received event with empty result field (oneof is nil)")
	}
}

// handleErrorEvent persists a repository collection failure reported by
// the Collector, marking the repository with an ERROR status and the
// reported error code. Storage failures are logged but not retried.
func (cc *ConsumerClient) handleErrorEvent(ctx context.Context, fullName string, repoErr *commonpb.ErrorResponse) {
	code := toDomainErrorCode(repoErr.GetCode())

	cc.log.Warn("collector failed to process repository, saving error status to DB",
		"repo", fullName,
		"code", code,
		"message", repoErr.GetMessage(),
	)

	repoEntity := domain.Repository{
		FullName:  fullName,
		Status:    domain.StatusError,
		ErrorCode: code,
	}

	if err := cc.storage.UpdateRepoStatus(ctx, &repoEntity); err != nil {
		cc.log.Error("failed to save error status to the database",
			"repo", fullName,
			"error", err,
		)
	}
}

// handleSuccessEvent persists successfully collected repository data
// reported by the Collector, marking the repository with a SUCCESS
// status. Storage failures are logged but not retried.
func (cc *ConsumerClient) handleSuccessEvent(ctx context.Context, fullName string, success *commonpb.CollectRepoSuccess) {
	repoEntity := domain.Repository{
		FullName:        fullName,
		Description:     success.GetDescription(),
		StargazersCount: int(success.GetStargazersCount()),
		ForksCount:      int(success.GetForksCount()),
		CreatedAt:       success.GetCreatedAt().AsTime(),
		Status:          domain.StatusSuccess,
		ErrorCode:       domain.ErrorCodeUnspecified,
	}

	if err := cc.storage.InsertRepo(ctx, &repoEntity); err != nil {
		cc.log.Error("failed to save data from Kafka to the database",
			"repo", repoEntity.FullName,
			"error", err,
		)
		return
	}

	cc.log.Info("repository data has been successfully updated to the database from Kafka",
		"repo", repoEntity.FullName,
	)
}

// toDomainErrorCode converts a wire-level commonpb.ErrorCode into the
// domain's ErrorCode representation. Unknown or unspecified codes map to
// domain.ErrorCodeUnspecified rather than failing, so that a newly added
// error code on the Collector side doesn't break the Processor.
func toDomainErrorCode(code commonpb.ErrorCode) domain.ErrorCode {
	switch code {
	case commonpb.ErrorCode_ERROR_CODE_REPOSITORY_NOT_FOUND:
		return domain.ErrorCodeRepositoryNotFound
	case commonpb.ErrorCode_ERROR_CODE_GITHUB_RATE_LIMIT_EXCEEDED:
		return domain.ErrorCodeGitHubRateLimitExceeded
	case commonpb.ErrorCode_ERROR_CODE_INTERNAL_COLLECTOR_ERROR:
		return domain.ErrorCodeInternalCollectorError
	default:
		return domain.ErrorCodeUnspecified
	}
}

func (cc *ConsumerClient) Close() error {
	return cc.reader.Close()
}
