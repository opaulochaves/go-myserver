package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func DBConnection() *pgxpool.Pool {
	//	# Example URL
	//	postgres://jack:secret@pg.example.com:5432/mydb?sslmode=verify-ca&pool_max_conns=10
	config, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse db config: %v\n", err)
		os.Exit(1)
	}

	ctx := context.Background()
	dbpool, err := pgxpool.NewWithConfig(ctx, config)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to the database: %v\n", err)
		os.Exit(1)
	}

	// Try to ping database.
	if err := dbpool.Ping(ctx); err != nil {
		defer dbpool.Close()
		fmt.Fprintf(os.Stderr, "Unable to connect to the database: %v\n", err)
		os.Exit(1)
	}

	return dbpool
}
