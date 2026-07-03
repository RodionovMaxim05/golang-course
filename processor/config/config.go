package config

import (
	"fmt"
	"repo-watcher/platform/env"
	"repo-watcher/platform/grpcserver"
	"repo-watcher/platform/logger"
)

type App struct {
	AppName string `yaml:"app_name" env:"APP_NAME" env-default:"repo-watcher-processor"`
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

type Database struct {
	User     string `yaml:"user"     env:"PROC_POSTGRES_USER"     env-default:"proc_user"`
	Password string `yaml:"password" env:"PROC_POSTGRES_PASSWORD" env-default:"proc_password"`
	Host     string `yaml:"host"     env:"PROC_POSTGRES_HOST"     env-default:"localhost"`
	Port     int    `yaml:"port"     env:"PROC_POSTGRES_PORT"     env-default:"5432"`
	Name     string `yaml:"name"     env:"PROC_POSTGRES_DB"       env-default:"processor_db"`
	SSLMode  string `yaml:"sslmode"  env:"PROC_POSTGRES_SSLMODE"  env-default:"disable"`
}

func (d Database) URL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		d.User,
		d.Password,
		d.Host,
		d.Port,
		d.Name,
		d.SSLMode,
	)
}

type Config struct {
	App      App               `yaml:"app"`
	Services Services          `yaml:"services"`
	GRPC     grpcserver.Config `yaml:"grpc"`
	Logger   logger.Config     `yaml:"logger"`
	Database Database          `yaml:"database"`
	Kafka    Kafka             `yaml:"kafka"`
}

func MustLoad(path string) Config {
	var cfg Config
	env.MustLoad(path, &cfg)
	return cfg
}
