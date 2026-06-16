package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"

	"repo-stat/collector/config"
	"repo-stat/collector/internal/adapter/github"
	kafkaAdapter "repo-stat/collector/internal/adapter/kafka"
	"repo-stat/collector/internal/adapter/subscriber"
	kafkaController "repo-stat/collector/internal/controller/kafka"
	"repo-stat/collector/internal/usecase"
	"repo-stat/platform/logger"
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

	// usecases & controllers
	getRepoUsecase := usecase.NewRepoUsecase(githubClient)
	getSubscriptionsInfoUsecase := usecase.NewGetSubscriptionsInfoUsecase(log, subscriberClient, taskProducerAdapter)
	kafkaHandler := kafkaController.NewRepoWorker(log, tasksReader, getRepoUsecase, resultProducerClient)

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
		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()

		log.Info("background subscription updater ticker started (15s)")

		for {
			select {
			case <-ctx.Done():
				log.Info("background subscription updater ticker stopped")
				return
			case <-ticker.C:
				log.Debug("ticker triggered: updating subscriptions...")
				if err := getSubscriptionsInfoUsecase.Execute(ctx); err != nil {
					log.Error("failed to update subscriptions in background", "error", err)
				}
			}
		}
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
	defer cancel()

	if err := run(ctx); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
