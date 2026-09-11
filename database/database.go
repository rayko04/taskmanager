package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(url string) (*pgxpool.Pool, error) {

	pool, err := pgxpool.New(context.Background(), url)
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
