package db

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
)

func CreateJob(ctx context.Context, dbConnection *sql.DB, originalURL string) (uuid.UUID, error) {
	query := "INSERT INTO jobs (status, original_url) VALUES ('pending', $1) RETURNING id"
	var id uuid.UUID
	err := dbConnection.QueryRowContext(ctx, query, originalURL).Scan(&id)
	if err != nil { 
		println(err.Error())
		return uuid.Nil, err
	}
	return id, nil
}
