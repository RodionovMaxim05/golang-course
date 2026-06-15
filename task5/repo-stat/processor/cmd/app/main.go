package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/jackc/pgx/v5/pgxpool"

	"repo-stat/platform/grpcserver"
	"repo-stat/platform/logger"
	"repo-stat/processor/config"
	"repo-stat/processor/internal/adapter/kafka"
	"repo-stat/processor/internal/adapter/postgres"
	grpccontroller "repo-stat/processor/internal/controller/grpc"
	"repo-stat/processor/internal/usecase"
	processorpb "repo-stat/proto/processor"
)

func run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// config
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "server configuration file")
	flag.Parse()

	cfg := config.MustLoad(configPath)

	// logger
	log := logger.MustMakeLogger(cfg.Logger.LogLevel)
	log.Info("starting processor server...")
	log.Debug("debug messages are enabled")

	// repo database
	pool, err := pgxpool.New(ctx, cfg.Database.URL())
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer pool.Close()

	// postgresql client
	dbClient := postgres.NewDBRepository(pool)

	// kafka producer client
	producerClient := kafka.NewProducerClient(cfg.Kafka, log)
	defer producerClient.Close()

	// kafka consumer client
	consumerClient := kafka.NewConsumerClient(cfg.Kafka, dbClient, log)
	defer consumerClient.Close()

	go consumerClient.StartConsumer(ctx)

	// handlers
	pingUseCase := usecase.NewPing()
	getRepoUseCase := usecase.NewGetRepo(producerClient, dbClient, log)
	getSubscriptionsInfoUseCase := usecase.NewGetSubscriptionsInfo(dbClient)
	processorServer := grpccontroller.NewServer(log, pingUseCase, getRepoUseCase, getSubscriptionsInfoUseCase)

	// server
	srv, err := grpcserver.New(cfg.GRPC.Address)
	if err != nil {
		return fmt.Errorf("create grpc server: %w", err)
	}

	processorpb.RegisterProcessorServer(srv.GRPC(), processorServer)

	if err := srv.Run(ctx); err != nil {
		return fmt.Errorf("run grpc server: %w", err)
	}

	cancel()
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
