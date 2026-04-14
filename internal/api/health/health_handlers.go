package health

import (
	"database/sql"
	"fmt"
	"net/http"
	"image-processing-pipeline/internal/minio"
	"image-processing-pipeline/internal/queue"
	"image-processing-pipeline/internal/cache"
)

type HealthChecker struct {
	name string
	Ping func() error
}

func HealthHandler(w http.ResponseWriter, r *http.Request, ) {
	fmt.Fprintln(w, "ok")
}

func HealthzHandler(dbConnection *sql.DB, minioClient *minio.MinioClient, taskDistributor *queue.TaskDistributor, redisClient *cache.RedisClient, w http.ResponseWriter, r *http.Request) {
	services := []HealthChecker{
		{name: "Database", Ping: dbConnection.Ping},
		{name: "Minio", Ping: minioClient.Ping},
		{name: "TaskDistributor", Ping: taskDistributor.Ping},
		{name: "Redis", Ping: redisClient.Ping},
	}

	for _, service := range services {
		if err := service.Ping(); err != nil {
			http.Error(w, fmt.Sprintf("%s connection failed: %v", service.name, err), http.StatusInternalServerError)
			return
		}
	}

	fmt.Fprintln(w, "ok")
}