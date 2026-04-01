package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"image-processing-pipeline/internal/db"
	"image-processing-pipeline/internal/minio"
)

const MAX_UPLOAD_SIZE = 1024 * 1024 * 5 // 5 MB

func HealthHandler(w http.ResponseWriter, r *http.Request, ) {
	fmt.Fprintln(w, "ok")
}

func UploadHandler(dbConnection *sql.DB, minioClient *minio.MinioClient, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.Body == nil {
		http.Error(w, "Request body and image file are required", http.StatusBadRequest)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		http.Error(w, "Request body too large", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Image file is required", http.StatusBadRequest)
		return
	}
	defer file.Close()



	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	originalURL, err := minio.WriteImage(minioClient, file, handler.Size)
	if err != nil {
		http.Error(w, "Failed to upload file to minio", http.StatusInternalServerError)
		return
	}
	
	// we need to write this originalURL to database and create a new entry,
	ctx := r.Context()
	jobID, err := db.CreateJob(ctx, dbConnection, originalURL)
	if err != nil {
		fmt.Println("DATABASE ERROR:", err)
		http.Error(w, "Failed to create job", http.StatusInternalServerError)
		return
	}
	
	// return the job id
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("%d", jobID)))
}

