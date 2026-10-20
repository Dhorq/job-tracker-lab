package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@127.0.0.1:5432/jobtracker?sslmode=disable"
	}

	var db *sql.DB
	var err error

	for i := 0; i < 10; i++ {
		db, err = sql.Open("postgres", dsn)
		if err != nil {
			time.Sleep(2 * time.Second)
			continue
		}

		err = db.Ping()
		if err == nil {
			return db, nil
		}

		db.Close()
		time.Sleep(2 * time.Second)
	}

	if db != nil {
		db.Close()
	}

	return nil, fmt.Errorf("gagal terhubung ke database setelah 10 percobaan: %w", err)
}
