package database

import (
	"database/sql"
	"fmt"
)

func Migrate(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS job_applications (
		id SERIAL PRIMARY KEY,
		company_name VARCHAR(255) NOT NULL,
		position VARCHAR(255) NOT NULL,
		status VARCHAR(50) NOT NULL,
		applied_date TIMESTAMP NOT NULL,
		notes TEXT,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW()
	);`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("Gagal memigrasi database: %w", err)
	}

	return nil
}
