package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DBConfig struct {
	ConnectionString string
}

func Connect(cfg DBConfig) (*pgxpool.Pool, error) {
	ctx := context.Background()
	return pgxpool.New(ctx, cfg.ConnectionString)
}
