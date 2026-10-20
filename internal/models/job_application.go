package models

import "time"

type JobApplication struct {
	ID          int       `json:"id"`
	CompanyName string    `json:"company_name"`
	Position    string    `json:"position"`
	Status      string    `json:"status"`
	AppliedDate time.Time `json:"applied_date"`
	Notes       string    `json:"notes"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
