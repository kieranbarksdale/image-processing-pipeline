package main

import (
	"fmt"
	"image-processing-pipeline/internal/config"
	"log"
	"github.com/joho/godotenv"
	//"image-processing-pipeline/internal/db"
	"net/http"
	"image-processing-pipeline/internal/api"
	"image-processing-pipeline/internal/queue"
	"github.com/hibiken/asynq"
)

func main() {
	// load env variables
	godotenv.Load()
	
	// load config
	cfg := config.Load() 

	// Database migration 
	// Database connection 
	// db, err := db.NewConnection(cfg.DatabaseURL)
	// if err != nil {
	// 	log.Fatal(err)
	// } 
	// fmt.Printf("Database connection established: %+v\n", db.Stats())

	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr: fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		},
		asynq.Config{
			Concurrency: 10,
		},
	)


	mux := asynq.NewServeMux()
	mux.HandleFunc("image:process", queue.ProcessImage)
	
	http.HandleFunc("/health", api.HealthHandler)

	if err := srv.Run(mux); err != nil {
		log.Fatal(err)
	}
}
