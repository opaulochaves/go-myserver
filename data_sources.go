package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/opaulochaves/myserver/config"
)

type dataSources struct {
	DBPool *pgxpool.Pool
}

// InitDS establishes connections to fields in dataSources
func initDS(ctx context.Context, cfg config.Config) (*dataSources, error) {
	log.Printf("Initializing data sources\n")

	log.Printf("Connecting to Postgresql\n")

	//	# Example URL
	//	postgres://jack:secret@pg.example.com:5432/mydb?sslmode=verify-ca&pool_max_conns=10
	config, err := pgxpool.ParseConfig(cfg.DatabaseUrl)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse db config: %w", err)
	}

	dbpool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to the database: %w", err)
	}

	return &dataSources{
		DBPool: dbpool,
	}, nil
}
