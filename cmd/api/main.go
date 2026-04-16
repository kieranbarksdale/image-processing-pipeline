package main

import (
	"fmt"
	"net/http"
	"image-processing-pipeline/internal/db" 
	"image-processing-pipeline/internal/config"
	"log"
	"github.com/joho/godotenv"
	"image-processing-pipeline/internal/api"
	"image-processing-pipeline/internal/api/health"
	"image-processing-pipeline/internal/minio"
	"image-processing-pipeline/internal/queue"
	"image-processing-pipeline/internal/cache"
	"github.com/go-chi/chi/v5"
)

func main() {

	// load env vairables 
	godotenv.Load()
	cfg := config.Load()

	// Get Chi router 
	router := chi.NewRouter()

	// migrations will go here before we connecct to the DB
	if err := db.Migrate(cfg.DatabaseURL, cfg.MigrationsPath); err != nil {
		log.Fatal("Error running database migrations: ", err)
	}

	// setup database connection and connections 
	dbConnection, err := db.NewConnection(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	} 
	log.Println("Database connected:", dbConnection.Stats())

	// connect to minio 
	minioClient, err := minio.NewClient(cfg.MinioEndpoint, cfg.MinioAccessKey, cfg.MinioSecretKey, false)
	if err != nil {
		log.Fatal("Error connecting to minio: ", err)
	}
	log.Println("Minio connected:", minioClient)
	
	redisClient, err := cache.CreateRedisClient(cfg.RedisURL)
	if err != nil {
		log.Fatal("Error connecting to redis: ", err)
	}
	log.Println("Redis connected:", redisClient)

	taskDistributor := queue.CreateQueue(fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort))
	log.Println("Queue connected:", taskDistributor)

	// setup http handlers and routes
	router.HandleFunc("/health", health.HealthHandler)

	router.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		health.HealthzHandler(dbConnection, minioClient, taskDistributor, redisClient, w, r)
	})

	router.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		api.UploadHandler(dbConnection, minioClient, taskDistributor, redisClient, w, r)
	})

	router.HandleFunc("/status/{jobId}", func(w http.ResponseWriter, r *http.Request) {
		api.StatusHandler(dbConnection, redisClient, w, r)
	})

	router.HandleFunc("/images/{jobId}", func(w http.ResponseWriter, r *http.Request) {
		api.GetImageHandler(dbConnection, minioClient, redisClient, w, r)
	})

	// start server
	err = http.ListenAndServe(cfg.ServerPort, router)
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
