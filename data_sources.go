package main

import (
	"context"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/opaulochaves/myserver/config"
	"log"
)

type dataSources struct {
	DB *sqlx.DB
}

// InitDS establishes connections to fields in dataSources
func initDS(ctx context.Context, cfg config.Config) (*dataSources, error) {
	log.Printf("Initializing data sources\n")

	log.Printf("Connecting to Postgresql\n")

	//	postgres://jack:secret@pg.example.com:5432/mydb?sslmode=verify-ca&pool_max_conns=10
	db, err := sqlx.Connect("pgx", cfg.DatabaseUrl)
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to the database: %w", err)
	}

	return &dataSources{DB: db}, nil
}
