package db

import (
	"context"
	"database/sql"
)

func IncrementTries(ctx context.Context, db *sql.DB, jobID string) error {
	_, err := db.ExecContext(ctx, "UPDATE jobs SET tries = tries + 1 WHERE id = $1", jobID)
	if err != nil {
		return err
	}
	return nil
}
