package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"repo-stat/platform/grpcserver"
	"repo-stat/platform/logger"
	"repo-stat/processor/config"
	"repo-stat/processor/internal/adapter/collector"
	grpccontroller "repo-stat/processor/internal/controller/grpc"
	"repo-stat/processor/internal/usecase"
	processorpb "repo-stat/proto/processor"
)

func run(ctx context.Context) error {
	// config
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "server configuration file")
	flag.Parse()

	cfg := config.MustLoad(configPath)

	// logger
	log := logger.MustMakeLogger(cfg.Logger.LogLevel)
	log.Info("starting processor server...")
	log.Debug("debug messages are enabled")

	// collector client
	collectorClient, err := collector.NewClient(cfg.Services.Collector, log)
	if err != nil {
		log.Error("cannot init collector adapter", "error", err)
		return err
	}

	// handlers
	pingUseCase := usecase.NewPing()
	getRepoUseCase := usecase.NewGetRepo(collectorClient)
	processorServer := grpccontroller.NewServer(log, pingUseCase, getRepoUseCase)

	// server
	srv, err := grpcserver.New(cfg.GRPC.Address)
	if err != nil {
		return fmt.Errorf("create grpc server: %w", err)
	}

	processorpb.RegisterProcessorServer(srv.GRPC(), processorServer)

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
