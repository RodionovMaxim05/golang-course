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
	subscriberpb "repo-stat/proto/subscriber"
	"repo-stat/subscriber/config"
	"repo-stat/subscriber/internal/adapter/github"
	"repo-stat/subscriber/internal/adapter/postgres"
	"repo-stat/subscriber/internal/controller/grpc"
	"repo-stat/subscriber/internal/usecase"
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

	// database
	db, err := pgxpool.New(ctx, cfg.Database.URL())
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	// handlers
	githubClient := github.NewGitHubClient(cfg.GitHub, log)
	subscriptionRepo := postgres.NewSubscriptionRepository(db)

	pingUseCase := usecase.NewPing()
	subscribeUseCase := usecase.NewSubscriptionService(githubClient, subscriptionRepo)
	unsubscribeUseCase := usecase.NewUnsubscribe(subscriptionRepo)
	getSubsUseCase := usecase.NewGetSubscriptions(subscriptionRepo)

	grpcHandler := grpc.NewServer(
		log,
		pingUseCase,
		subscribeUseCase,
		unsubscribeUseCase,
		getSubsUseCase,
	)

	// server
	srv, err := grpcserver.New(cfg.GRPC.Address)
	if err != nil {
		return fmt.Errorf("create grpc server: %w", err)
	}

	subscriberpb.RegisterSubscriberServer(srv.GRPC(), grpcHandler)

	if err := srv.Run(ctx); err != nil {
		return fmt.Errorf("run grpc server: %w", err)
	}

	return nil
}

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	if err := run(ctx); err != nil {
		_, err = fmt.Fprintln(os.Stderr, err)
		if err != nil {
			fmt.Printf("launching server error: %s\n", err)
		}
		cancel()
		os.Exit(1)
	}
}
