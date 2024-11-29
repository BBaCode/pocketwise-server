package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBConfig struct {
	ConnectionString string
}

func Connect(cfg DBConfig) (*pgxpool.Pool, error) {
	// Parse the connection string into a pgxpool.Config
	poolConfig, err := pgxpool.ParseConfig(cfg.ConnectionString)
	if err != nil {
		log.Fatalf("Failed to parse pool config: %v", err)
	}

	// Set PreferSimpleProtocol to true for compatibility with certain queries
	poolConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	// Create a new connection pool with the customized configuration
	ctx := context.Background()
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Fatalf("Failed to create pool: %v", err)
	}
	return pool, nil
}
