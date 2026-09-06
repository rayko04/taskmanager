package database

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"context"
)

func NewPool() (*pgxpool.Pool, error) {

	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}

	err = pool.Ping(context.Background())
	if err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}
