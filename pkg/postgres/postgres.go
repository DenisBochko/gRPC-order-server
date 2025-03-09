package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type Postgres struct {
	Host     string `yaml:"POSTGRES_HOST" env:"POSTGRES_HOST" env-default:"localhost"`
	Port     string `yaml:"POSTGRES_PORT" env:"POSTGRES_PORT" env-default:"5432"`
	Username string `yaml:"POSTGRES_USER" env:"POSTGRES_USER" env-default:"postgres"`
	Password string `yaml:"POSTGRES_PASS" env:"POSTGRES_PASS" env-default:"postgres"`
	Database string `yaml:"POSTGRES_DB" env:"POSTGRES_DB" env-default:"test_db"`
	Sslmode  string `yaml:"POSTGRES_SSLMODE" env:"POSTGRES_SSLMODE" env-default:"disable"`
}

func New(config Postgres) (*pgx.Conn, error) {
	// urlExample := "postgres://username:password@localhost:5432/database_name"
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", config.Username, config.Password, config.Host, config.Port, config.Database)
	conn, err := pgx.Connect(context.Background(), connString)

	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	return conn, nil
}
