package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

// New открывает подключение к Postgres по DSN
func New(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
