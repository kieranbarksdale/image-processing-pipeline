package db

import (
	"context"
	"database/sql"
)

func GetImageKeys(ctx context.Context, db *sql.DB, jobId string) ([]string, error) {
	rows, err := db.QueryContext(ctx, "SELECT key FROM images WHERE job_id = ?", jobId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var keys []string
	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}
	return keys, nil
}
