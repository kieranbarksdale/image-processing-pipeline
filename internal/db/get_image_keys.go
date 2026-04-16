package db

import (
	"context"
	"database/sql"
	"log"
)

func GetImageKeys(ctx context.Context, db *sql.DB, jobId string) ([]string, error) {
	rows, err := db.QueryContext(ctx, "SELECT img_key FROM images WHERE job_id = $1", jobId)
	if err != nil {
		log.Println("DEBUG: Failed to get image keys", err)
		return nil, err
	}
	defer rows.Close()
	
	var keys []string
	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			log.Println("DEBUG: Failed to scan image key", err)
			return nil, err
		}
		keys = append(keys, key)
	}
	return keys, nil
}
