package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@127.0.0.1:5432/jobtracker?sslmode=disable"
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("Gagal membuka koneksi: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("Gagal memanggil database: %w", err)
	}

	return db, nil
}
