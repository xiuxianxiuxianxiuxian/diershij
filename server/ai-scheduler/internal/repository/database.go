package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/cultivation-world/shared/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	pool *pgxpool.Pool
}

func NewDatabase(cfg *config.DatabaseConfig) (*Database, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return &Database{pool: pool}, nil
}

func (db *Database) Close() {
	db.pool.Close()
}

func (db *Database) Pool() *pgxpool.Pool {
	return db.pool
}
