package http

import (
	"context"
	"log/slog"
	"net/http"

	"repo-watcher/api/config"
	"repo-watcher/api/internal/adapter/processor"
	"repo-watcher/api/internal/adapter/subscriber"
	"repo-watcher/api/internal/usecase"
)

func NewHandler(ctx context.Context, log *slog.Logger, cfg config.Config) (http.Handler, error) {
	subscriberClient, err := subscriber.NewClient(cfg.Services.Subscriber, log)
	if err != nil {
		log.Error("cannot init subscriber adapter", "error", err)
		return nil, err
	}

	processorClient, err := processor.NewClient(cfg.Services.Processor, log)
	if err != nil {
		log.Error("cannot init processor adapter", "error", err)
		return nil, err
	}

	pingUseCase := usecase.NewPing(subscriberClient, processorClient)
	subscribeUseCase := usecase.NewSubscriber(subscriberClient)
	unsubscribeUseCase := usecase.NewUnsubscriber(subscriberClient)
	getSubscriptionsUseCase := usecase.NewGetSubscriptions(subscriberClient)
	subscriptionsInfoUseCase := usecase.NewGetSubscriptionsInfo(processorClient)
	getRepoUseCase := usecase.NewGetRepo(processorClient)

	mux := http.NewServeMux()
	AddRoutes(
		mux,
		log,
		pingUseCase,
		getRepoUseCase,
		subscribeUseCase,
		unsubscribeUseCase,
		getSubscriptionsUseCase,
		subscriptionsInfoUseCase,
	)

	return mux, nil
}
