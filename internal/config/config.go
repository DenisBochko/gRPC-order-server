package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	PortGRPC string `yaml:"GRPC_PORT" env:"GRPC_PORT" env-default:"50051"`
}

func New() (*Config, error) {
	var cfg Config

	err := cleanenv.ReadConfig("./config/config.yaml", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
