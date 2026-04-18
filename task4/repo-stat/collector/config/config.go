package config

import (
	"repo-stat/collector/internal/adapter"
	"repo-stat/platform/env"
	"repo-stat/platform/grpcserver"
	"repo-stat/platform/logger"
)

type App struct {
	AppName string `yaml:"app_name" env:"APP_NAME" env-default:"repo-stat-collector"`
}

type Config struct {
	App        App               `yaml:"app"`
	GitHub     adapter.Config    `yaml:"github"`
	Subscriber SubscriberConfig `yaml:"subscriber"`
	GRPC       grpcserver.Config `yaml:"grpc"`
	Logger     logger.Config     `yaml:"logger"`
}

type SubscriberConfig struct {
	Address string `yaml:"address" env:"SUBSCRIBER_ADDRESS" env-default:"subscriber:8081"`
}

func MustLoad(path string) Config {
	var cfg Config
	env.MustLoad(path, &cfg)
	return cfg
}
