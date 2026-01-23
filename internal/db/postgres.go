package postgres

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func DbPool() (*pgxpool.Pool, error) {
	ctx := context.Background()
	dbHost := os.Getenv("DATABASE_URL")
	dbPool, err := pgxpool.New(ctx, dbHost)

	return dbPool, err
}
