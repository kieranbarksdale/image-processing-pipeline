package db

import (
	"context"
	"database/sql"
)

func PutStatus(ctx context.Context, db *sql.DB, jobID string, status string) error {
	_, err := db.ExecContext(ctx, "UPDATE jobs SET status = $1 WHERE id = $2", status, jobID)
	if err != nil {
		return err
	}
	return nil
}
