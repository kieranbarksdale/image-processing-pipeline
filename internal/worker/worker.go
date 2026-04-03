package worker

import (
	"database/sql"
	"image-processing-pipeline/internal/minio"
)

type WorkerHandler struct {
	db *sql.DB
	minioClient *minio.MinioClient
}

func NewWorker(db *sql.DB, minioClient *minio.MinioClient) *WorkerHandler {
	return &WorkerHandler{
		db: db,
		minioClient: minioClient,
	}
}
