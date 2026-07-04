package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/segmentio/kafka-go"

	"repo-watcher/collector/config"
	"repo-watcher/collector/internal/adapter/github"
	kafkaAdapter "repo-watcher/collector/internal/adapter/kafka"
	"repo-watcher/collector/internal/adapter/subscriber"
	kafkaController "repo-watcher/collector/internal/controller/kafka"
	"repo-watcher/collector/internal/controller/scheduler"
	"repo-watcher/collector/internal/usecase"
	"repo-watcher/platform/logger"
)

func run(ctx context.Context) error {
	// config
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "server configuration file")
	flag.Parse()

	cfg := config.MustLoad(configPath)

	// logger
	log := logger.MustMakeLogger(cfg.Logger.LogLevel)
	log.Info("starting server...")
	log.Debug("debug messages are enabled")

	// adapters
	githubClient := github.NewGitHubClient(cfg.GitHub, log)
	subscriberClient, err := subscriber.NewSubscriberClient(cfg.Services.Subscriber, log)
	if err != nil {
		log.Error("cannot init subscriber adapter", "error", err)
		return err
	}
	defer func() {
		if err := subscriberClient.Close(); err != nil {
			log.Error("failed to close subscriber client", "error", err)
		}
	}()

	// kafka result producer client
	resultProducerClient := kafkaAdapter.NewResultProducerAdapter(cfg.Kafka, log)
	defer func() {
		if err := resultProducerClient.Close(); err != nil {
			log.Error("failed to close result producer client", "error", err)
		}
	}()

	// kafka consumer adapter
	tasksReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{cfg.Kafka.Address},
		Topic:       cfg.Kafka.ConsumerTopic,
		GroupID:     cfg.Kafka.GroupID,
		StartOffset: kafka.FirstOffset,
	})
	defer func() {
		if err := tasksReader.Close(); err != nil {
			log.Error("failed to close tasks reader", "error", err)
		}
	}()

	// kafka task producer adapter
	taskProducerAdapter := kafkaAdapter.NewTaskProducerAdapter(cfg.Kafka, log)
	defer func() {
		if err := taskProducerAdapter.Close(); err != nil {
			log.Error("failed to close task producer adapter", "error", err)
		}
	}()

	// usecases

	getRepoUsecase := usecase.NewGetRepoUsecase(githubClient)
	getSubscriptionsInfoUsecase := usecase.NewGetSubscriptionsInfoUsecase(log, subscriberClient, taskProducerAdapter)

	// controllers

	kafkaHandler := kafkaController.NewRepoWorker(log, tasksReader, getRepoUsecase, resultProducerClient)
	subscriptionUpdater := scheduler.NewSubscriptionUpdater(
		log,
		cfg.SubscriptionUpdater.Interval,
		getSubscriptionsInfoUsecase.Execute,
	)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		kafkaHandler.Start(ctx)
		log.Info("kafka handler worker stopped")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		subscriptionUpdater.Start(ctx)
	}()

	<-ctx.Done()
	log.Info("stopping collector server gracefully...")

	wg.Wait()
	log.Info("all background workers stopped cleanly")

	return nil
}

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)

	if err := run(ctx); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		cancel()
		os.Exit(1)
	}
}
