package db

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq" // драйвер PostgreSQL
)

// NewPostgresDB создаёт и настраивает подключение к PostgreSQL.
// Принимает DSN (Data Source Name) вида:
// "postgres://user:password@host:port/dbname?sslmode=disable"
func NewPostgresDB(ctx context.Context, dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Настройка пула подключений
	db.SetMaxOpenConns(25)                 // макс. открытых соединений
	db.SetMaxIdleConns(5)                  // макс. простаивающих соединений
	db.SetConnMaxLifetime(5 * time.Minute) // время жизни соединения

	// Проверка подключения
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
