package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	PortGRPC string `yaml:"GRPC_PORT" env:"GRPC_PORT" env-default:"50051"`
	PortHttp string `yaml:"HTTP_PORT" env:"HTTP_PORT" env-default:"8080"`
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
