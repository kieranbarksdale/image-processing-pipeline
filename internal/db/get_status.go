package db

import (
	"database/sql"
	"github.com/google/uuid"
)

func GetStatus(db *sql.DB, jobID uuid.UUID) (string, error) {
	var status string
	err := db.QueryRow("SELECT status FROM jobs WHERE id = $1", jobID).Scan(&status)
	if err != nil {
		return "", err
	}
	return status, nil
}
