package db

import (
	"database/sql"
	"context"
)

func CreateImage(ctx context.Context, db *sql.DB, jobId string, imageURL string) error {

	query := "INSERT INTO images (job_id, image_url) VALUES ($1, $2)"
	_, err := db.ExecContext(ctx, query, jobId, imageURL)
	if err != nil {
		return err
	}
	
	return nil
}