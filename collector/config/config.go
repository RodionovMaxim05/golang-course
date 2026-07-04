package config

import (
	"time"

	"repo-watcher/collector/internal/adapter/github"
	"repo-watcher/platform/env"
	"repo-watcher/platform/grpcserver"
	"repo-watcher/platform/logger"
)

type App struct {
	AppName string `yaml:"app_name" env:"APP_NAME" env-default:"repo-watcher-collector"`
}

type Services struct {
	Subscriber string `yaml:"subscriber" env:"SUBSCRIBER_ADDRESS" env-default:"localhost:8081"`
}

type Kafka struct {
	Address       string `yaml:"address"        env:"KAFKA_ADDRESS"        env-default:"kafka:9092"`
	ProducerTopic string `yaml:"producer_topic" env:"KAFKA_PRODUCER_TOPIC" required:"true"`
	ConsumerTopic string `yaml:"consumer_topic" env:"KAFKA_CONSUMER_TOPIC" required:"true"`
	GroupID       string `yaml:"group_id"       env:"KAFKA_GROUP_ID"       env-default:"default-group"`
}

type SubscriptionUpdater struct {
	Interval time.Duration `yaml:"interval" env:"SUBSCRIPTION_UPDATER_INTERVAL" env-default:"15s"`
}

type Config struct {
	App                 App                 `yaml:"app"`
	GitHub              github.Config       `yaml:"github"`
	Services            Services            `yaml:"services"`
	GRPC                grpcserver.Config   `yaml:"grpc"`
	Kafka               Kafka               `yaml:"kafka"`
	SubscriptionUpdater SubscriptionUpdater `yaml:"subscription_updater"`
	Logger              logger.Config       `yaml:"logger"`
}

func MustLoad(path string) Config {
	var cfg Config
	env.MustLoad(path, &cfg)
	return cfg
}
