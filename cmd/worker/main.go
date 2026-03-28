package main

import (
	"fmt"
	"image-processing-pipeline/internal/config"
	"log"
	"github.com/joho/godotenv"
	"image-processing-pipeline/internal/db"
	"net/http"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
}

func main() {
	// load env variables
	godotenv.Load()
	
	// load config
	cfg := config.Load() 

	// Database migration 
	// Database connection 
	db, err := db.NewConnection(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	} 

	fmt.Printf("Database connection established: %+v\n", db.Stats())
	
	http.HandleFunc("/health", healthHandler)

	// TODO: start server on port 
}
