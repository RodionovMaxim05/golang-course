package config

import (
	"fmt"

	"repo-stat/platform/env"
	"repo-stat/platform/grpcserver"
	"repo-stat/platform/logger"
	"repo-stat/subscriber/internal/adapter/github"
)

type App struct {
	AppName string `yaml:"app_name" env:"APP_NAME" env-default:"repo-stat-subscriber"`
}

type Database struct {
	Host     string `yaml:"host" env:"POSTGRES_HOST" env-default:"localhost"`
	Port     int    `yaml:"port" env:"POSTGRES_PORT" env-default:"5432"`
	User     string `yaml:"user" env:"POSTGRES_USER" env-default:"postgres"`
	Password string `yaml:"password" env:"POSTGRES_PASSWORD" env-default:"password"`
	Name     string `yaml:"name" env:"POSTGRES_DB" env-default:"repo_stat"`
	SSLMode  string `yaml:"sslmode" env:"POSTGRES_SSLMODE" env-default:"disable"`
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
	GRPC     grpcserver.Config `yaml:"grpc"`
	Logger   logger.Config     `yaml:"logger"`
	Database Database          `yaml:"database"`
	GitHub   github.Config     `yaml:"github"`
}

func MustLoad(path string) Config {
	var cfg Config
	env.MustLoad(path, &cfg)
	return cfg
}
