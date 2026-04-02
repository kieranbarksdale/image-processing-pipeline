package main

import (
	"fmt"
	"net/http"
	"image-processing-pipeline/internal/db" 
	"image-processing-pipeline/internal/config"
	"log"
	"github.com/joho/godotenv"
	"image-processing-pipeline/internal/api"
	"image-processing-pipeline/internal/minio"
)

func main() {

	// load env vairables 
	godotenv.Load()
	cfg := config.Load()

	// migrations will go here before we connecct to the DB
	if err := db.Migrate(cfg.DatabaseURL, cfg.MigrationsPath); err != nil {
		log.Fatal("Error running database migrations: ", err)
	}

	// setup database connection and connections 
	dbConnection, err := db.NewConnection(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	} 
	fmt.Println("Database connected:", dbConnection.Stats())

	// connect to minio 
	minioClient, err := minio.NewClient(cfg.MinioEndpoint, cfg.MinioAccessKey, cfg.MinioSecretKey, false)
	if err != nil {
		log.Fatal("Error connecting to minio: ", err)
	}
	fmt.Println("Minio connected:", minioClient)
	// do bucket testing here 

	// setup http handlers and routes
	http.HandleFunc("/health", api.HealthHandler)

	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		api.UploadHandler(dbConnection, minioClient, w, r)
	})

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		api.StatusHandler(dbConnection, w, r)
	})

	// start server
	err = http.ListenAndServe(cfg.ServerPort, nil)
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
