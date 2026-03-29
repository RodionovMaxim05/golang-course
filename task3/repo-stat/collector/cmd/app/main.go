package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"repo-stat/collector/config"
	"repo-stat/collector/internal/adapter"
	"repo-stat/collector/internal/controller"
	"repo-stat/collector/internal/usecase"
	"repo-stat/platform/grpcserver"
	"repo-stat/platform/logger"
	collectorpb "repo-stat/proto/collector"
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
	githubClient := adapter.NewGitHubClient(cfg.GitHub, log)
	getRepoUsecase := usecase.NewRepoUsecase(githubClient)
	grpcHandler := controller.NewRepoHandler(log, getRepoUsecase)

	// server
	srv, err := grpcserver.New(cfg.GRPC.Address)
	if err != nil {
		return fmt.Errorf("create grpc server: %w", err)
	}

	collectorpb.RegisterRepoServiceServer(srv.GRPC(), grpcHandler)

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
			panic(err)
		}
		cancel()
		os.Exit(1)
	}
}
