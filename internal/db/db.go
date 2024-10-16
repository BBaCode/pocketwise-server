package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

func Connect(cfg Config) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)

	ctx := context.Background()
	return pgxpool.New(ctx, connString)
}
