package config

import (
	"repo-watcher/platform/env"
	"repo-watcher/platform/httpserver"
	"repo-watcher/platform/logger"
)

type App struct {
	AppName string `yaml:"app_name" env:"APP_NAME" env-default:"repo-watcher-api"`
}

type Services struct {
	Subscriber string `yaml:"subscriber" env:"SUBSCRIBER_ADDRESS" env-default:"localhost:8081"`
	Processor  string `yaml:"processor" env:"PROCESSOR_ADDRESS" env-default:"localhost:8083"`
}

type Config struct {
	App      App               `yaml:"app"`
	Services Services          `yaml:"services"`
	HTTP     httpserver.Config `yaml:"http"`
	Logger   logger.Config     `yaml:"logger"`
}

func MustLoad(path string) Config {
	var cfg Config
	env.MustLoad(path, &cfg)
	return cfg
}
