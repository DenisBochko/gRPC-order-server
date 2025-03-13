package config

import (
	"order-server/pkg/postgres"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	postgres.PostgresCfg
	PortGRPC    string `yaml:"GRPC_PORT" env:"GRPC_PORT"`
	PortHttp    string `yaml:"HTTP_PORT" env:"HTTP_PORT"`
	Environment string `yaml:"ENV" env:"ENV"`
}

func NewYAML() (*Config, error) {
	var cfg Config

	err := cleanenv.ReadConfig("./config/config.yaml", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func NewENV() (*Config, error) {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
