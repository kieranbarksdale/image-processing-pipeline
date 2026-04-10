package db

import (
	"database/sql"
	"context"
	"fmt"
)

func CreateImage(ctx context.Context, db *sql.DB, jobId string, key string, size string) error {

	fmt.Println("Adding to imgaes table with job_id:", jobId, "and key:", key)
	query := "INSERT INTO images (job_id, img_key, size) VALUES ($1, $2, $3)"
	_, err := db.ExecContext(ctx, query, jobId, key, size)
	if err != nil {
		fmt.Println("Error creating image:", err)
		return err
	}
	
	return nil
}