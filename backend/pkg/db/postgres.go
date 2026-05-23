package db

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

const (
	maxOpenConns      = 25
	maxIdleConns      = 5
	connMaxLifetime   = 5 * time.Minute
	connectionTimeout = 5 * time.Second
)

// NewPostgresDB создает и настраивает подключение к PostgreSQL.
// Ожидает DSN в формате: "postgres://user:password@host:port/dbname?sslmode=disable".
func NewPostgresDB(ctx context.Context, dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(connMaxLifetime)

	// Создаем контекст с таймаутом для проверки подключения
	pingCtx, cancel := context.WithTimeout(ctx, connectionTimeout)
	defer cancel()

	if err := db.PingContext(pingCtx); err != nil {
		if closeErr := db.Close(); closeErr != nil {
			zap.L().Error("Failed to close DB connection after ping failure", zap.Error(closeErr))
		}
		return nil, err
	}

	zap.L().Info("Database connection established",
		zap.Int("max_open_conns", maxOpenConns),
		zap.Int("max_idle_conns", maxIdleConns),
		zap.Duration("conn_max_lifetime", connMaxLifetime),
	)

	return db, nil
}
