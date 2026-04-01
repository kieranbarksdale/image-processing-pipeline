package db

import (
	"context"
	"database/sql"
)

func CreateJob(ctx context.Context, dbConnection *sql.DB, originalURL string) (int64, error) {
	query := "INSERT INTO jobs (status, original_url) VALUES ('pending', $1) RETURNING id"
	var id int64
	err := dbConnection.QueryRowContext(ctx, query, originalURL).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}
