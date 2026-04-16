package main

import (
	"fmt"
	"image-processing-pipeline/internal/config"
	"log"
	"github.com/joho/godotenv"
	"image-processing-pipeline/internal/db"
	"image-processing-pipeline/internal/worker"
	"github.com/hibiken/asynq"
	"image-processing-pipeline/internal/minio"
)


func main() {
	// load env variables
	godotenv.Load()
	cfg := config.Load() 

	// Database migration 
	if err := db.Migrate(cfg.DatabaseURL, cfg.MigrationsPath); err != nil {
		log.Fatal("Error running database migrations: ", err)
	}

	// Database connection 
	db, err := db.NewConnection(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	} 
	log.Printf("Database connection established: %+v\n", db.Stats())

	//Minio connection 
	minioClient, err := minio.NewClient(cfg.MinioEndpoint, cfg.MinioAccessKey, cfg.MinioSecretKey, false)
	if err != nil {
		log.Fatal("Error connecting to minio: ", err)
	}
	log.Println("Minio connected:", minioClient)

	// Asynq server setup
	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr: fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		},
		asynq.Config{
			Concurrency: 10,
		},
	)

	// Set up the worker
	var workerHandler *worker.WorkerHandler
	workerHandler = worker.NewWorker(db, minioClient)

	// Asynq mux setup
	mux := asynq.NewServeMux()

	mux.HandleFunc("image:process", workerHandler.HandleProcessJob)	

	// Run the server
	if err := srv.Run(mux); err != nil {
		log.Fatal(err)
	}
}
