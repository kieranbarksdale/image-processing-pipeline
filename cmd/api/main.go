package main

import (
	"fmt"
	"net/http"
	"image-processing-pipeline/internal/db" 
	"image-processing-pipeline/internal/config"
	"log"
	"github.com/joho/godotenv"
	"image-processing-pipeline/internal/api"
)

func main() {

	// load env vairables 
	godotenv.Load()
	cfg := config.Load()

	// migrations will go here before we connecct to the DB
	if err := db.Migrate(cfg.DatabaseURL); err != nil {
		log.Fatal("Error running database migrations: ", err)
	}

	// setup database connection and connections 
	dbConnection, err := db.NewConnection(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	} 
	fmt.Println("Database connected:", dbConnection.Stats())

	// setup http handlers and routes
	http.HandleFunc("/health", api.HealthHandler)

	// start server
	err = http.ListenAndServe(cfg.ServerPort, nil)
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
