package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Dhorq/job-tracker-lab/internal/models"
)

func Create(db *sql.DB, app models.JobApplication) (int, error) {
	query := `
		INSERT INTO job_applications (company_name, position, status, applied_date, notes, created_at, updated_at)
		VALUES($1, $2, $3, $4, $5, $6, $7) 
		RETURNING id
	`

	var id int
	err := db.QueryRow(
		query,
		app.CompanyName,
		app.Position,
		app.Status,
		app.AppliedDate,
		app.Notes,
		time.Now(),
		time.Now(),
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("Gagal menyimpan data: %w", err)
	}

	return id, nil
}

func GetAll(db *sql.DB) ([]models.JobApplication, error) {
	query := `SELECT id, company_name, position, status, applied_date, notes, created_at, updated_at FROM job_applications ORDER BY applied_date DESC`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("Gagal mengambil data: %w", err)
	}
	defer rows.Close()

	var apps []models.JobApplication
	for rows.Next() {
		var a models.JobApplication
		err := rows.Scan(&a.ID, &a.CompanyName, &a.Position, &a.Status, &a.AppliedDate, &a.Notes, &a.CreatedAt, &a.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("Gagal mengecek data: %w", err)
		}
		apps = append(apps, a)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("Error saat iterasi: %w", err)
	}

	return apps, nil
}
