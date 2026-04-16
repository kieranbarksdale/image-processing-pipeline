package db

import (
	"database/sql"
	"context"
	"log"
)

func CreateImage(ctx context.Context, db *sql.DB, jobId string, key string, size string) error {
	query := "INSERT INTO images (job_id, img_key, size) VALUES ($1, $2, $3)"
	_, err := db.ExecContext(ctx, query, jobId, key, size)
	if err != nil {
		log.Printf("Error creating image: %v", err)
		return err
	}
	
	return nil
}