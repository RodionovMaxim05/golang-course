// @title       Repo Stat API
// @version     1.0
// @description API Gateway for repository statistics
// @host        localhost:28080
// @BasePath    /
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/redis/go-redis/v9"

	"repo-watcher/api/config"
	_ "repo-watcher/api/docs"
	"repo-watcher/api/internal/adapter/cache"
	"repo-watcher/api/internal/adapter/processor"
	"repo-watcher/api/internal/adapter/ratelimiter"
	"repo-watcher/api/internal/adapter/subscriber"
	"repo-watcher/api/internal/controller/http"
	"repo-watcher/api/internal/usecase"
	"repo-watcher/platform/httpserver"
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

	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.Redis.Address,
	})
	defer func() {
		if err := redisClient.Close(); err != nil {
			log.Error("failed to close redis client", "error", err)
		}
	}()

	redisCache := cache.NewCache(redisClient, log)
	redisRateLimiter := ratelimiter.NewRedisRateLimiter(redisClient, cfg.RateLimit, log)

	inMemoryRateLimiter := ratelimiter.NewInMemoryRateLimiter(cfg.RateLimit, log)

	fallbackRateLimiter := ratelimiter.NewFallbackRateLimiter(redisRateLimiter, inMemoryRateLimiter, log)

	subscriberClient, err := subscriber.NewClient(cfg.Services.Subscriber, log)
	if err != nil {
		log.Error("cannot init subscriber adapter", "error", err)
		return err
	}
	defer func() {
		if err := subscriberClient.Close(); err != nil {
			log.Error("failed to close subscriber client", "error", err)
		}
	}()

	processorClient, err := processor.NewClient(cfg.Services.Processor, redisCache, cfg.Cache.TTL, log)
	if err != nil {
		log.Error("cannot init processor adapter", "error", err)
		return err
	}
	defer func() {
		if err := processorClient.Close(); err != nil {
			log.Error("failed to close processor client", "error", err)
		}
	}()

	// usercases

	pingUseCase := usecase.NewPing(subscriberClient, processorClient)
	subscribeUseCase := usecase.NewSubscriber(subscriberClient)
	unsubscribeUseCase := usecase.NewUnsubscriber(subscriberClient)
	getSubscriptionsUseCase := usecase.NewGetSubscriptions(subscriberClient)
	subscriptionsInfoUseCase := usecase.NewGetSubscriptionsInfo(processorClient)
	getRepoUseCase := usecase.NewGetRepo(processorClient)

	// handler
	handler := http.NewHandler(
		log, fallbackRateLimiter, pingUseCase, getRepoUseCase, subscribeUseCase, unsubscribeUseCase,
		getSubscriptionsUseCase, subscriptionsInfoUseCase,
	)

	// server
	srv := httpserver.New(cfg.HTTP, handler)
	if err := srv.Run(ctx); err != nil {
		return fmt.Errorf("run http server: %w", err)
	}
	return nil
}

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)

	if err := run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		cancel()
		os.Exit(1)
	}
}
