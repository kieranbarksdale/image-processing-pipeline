package db

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
)

func CreateJob(ctx context.Context, dbConnection *sql.DB, key string) (uuid.UUID, error) {
	query := "INSERT INTO jobs (status, img_key) VALUES ('pending', $1) RETURNING id"
	var id uuid.UUID
	err := dbConnection.QueryRowContext(ctx, query, key).Scan(&id)
	if err != nil { 
		println(err.Error())
		return uuid.Nil, err
	}
	return id, nil
}
