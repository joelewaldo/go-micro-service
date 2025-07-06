package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/joelewaldo/go-micro-service/internal/oauth2/config"
	_ "github.com/lib/pq"
)

func Connect(cfg *config.Config) (*sql.DB, error) {
	if cfg.DATABASE_URL == "" {
		return nil, errors.New("config: DATABASE_URL is empty")
	}

	db, err := sql.Open("postgres", cfg.DATABASE_URL)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("db.Ping: %w", err)
	}

	return db, nil
}
