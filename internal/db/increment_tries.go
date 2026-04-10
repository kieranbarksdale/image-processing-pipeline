package db

import (
	"context"
	"database/sql"
)

func (workerHandler *WorkerHandler) IncrementTries(ctx context.Context) error {
	_, err := workerHandler.db.ExecContext(ctx, "UPDATE jobs SET tries = tries + 1 WHERE id = $1", workerHandler.jobID)
	if err != nil {
		return err
	}
	return nil
}
