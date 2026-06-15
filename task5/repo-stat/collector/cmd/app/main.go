package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
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

	// handlers
	githubClient := github.NewGitHubClient(cfg.GitHub, log)
	subscriberClient, err := subscriber.NewSubscriberClient(cfg.Services.Subscriber, log)
	if err != nil {
		log.Error("cannot init subscriber adapter", "error", err)
		return err
	}

	// kafka producer client
	producerClient := kafkaAdapter.NewProducerClient(cfg.Kafka, log)
	defer func() {
		if err := producerClient.Close(); err != nil {
			log.Error("failed to close producer client", "error", err)
		}
	}()

	// kafka consumer client
	tasksReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{cfg.Kafka.Address},
		Topic:   cfg.Kafka.ConsumerTopic,
		GroupID: cfg.Kafka.GroupID,
	})
	defer func() {
		if err := tasksReader.Close(); err != nil {
			log.Error("failed to close tasks reader", "error", err)
		}
	}()

	getRepoUsecase := usecase.NewRepoUsecase(githubClient)
	getSubscriptionsInfoUsecase := usecase.NewGetSubscriptionsInfoUsecase(subscriberClient, githubClient, producerClient)
	kafkaHandler := kafkaController.NewRepoWorker(tasksReader, getRepoUsecase, producerClient, log)

	go kafkaHandler.Start(ctx)

	go func() {
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

	return nil
}

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	if err := run(ctx); err != nil {
		_, err = fmt.Fprintln(os.Stderr, err)
		if err != nil {
			panic(err)
		}
		cancel()
		os.Exit(1)
	}
}
