package api

import (
	"database/sql"
	"log"
	"net/http"
	"image-processing-pipeline/internal/db"
	"image-processing-pipeline/internal/minio"
	"image-processing-pipeline/internal/queue"
	"image-processing-pipeline/internal/cache"
	"github.com/google/uuid"
	"encoding/json"
	"io"
	"archive/zip"
	"context"
	"github.com/go-chi/chi/v5"
)

const MAX_UPLOAD_SIZE = 1024 * 1024 * 5 // 5 MB

func UploadHandler(dbConnection *sql.DB, minioClient *minio.MinioClient, taskDistributor *queue.TaskDistributor, redisClient *cache.RedisClient, w http.ResponseWriter, r *http.Request) {
	
	limited, err := redisClient.IsRateLimited(r.RemoteAddr) 
	if err != nil {
		log.Println("CACHE ERROR:", err)
		http.Error(w, "Cache service unavailable", http.StatusServiceUnavailable)
		return
	}
	if limited {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	if err != nil {
		http.Error(w, "Rate limit check failed", http.StatusInternalServerError)
		return
	}
	if limited {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}
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

	log.Printf("Uploaded File: %+v\n", handler.Filename)
	log.Printf("File Size: %+v\n", handler.Size)
	log.Printf("MIME Header: %+v\n", handler.Header)

	originalKey, err := minio.WriteImage(minioClient, file, handler.Size)
	if err != nil {
		http.Error(w, "Failed to upload file to minio", http.StatusInternalServerError)
		return
	}
	
	// we need to write this originalURL to database and create a new entry,
	ctx := r.Context()
	jobID, err := db.CreateJob(ctx, dbConnection, originalKey)
	if err != nil {
		log.Println("DATABASE ERROR:", err)
		http.Error(w, "Failed to create job", http.StatusInternalServerError)
		return
	}

	err = taskDistributor.AddToQueue(jobID.String(), originalKey)
	if err != nil {
		log.Println("QUEUE ERROR:", err)
		http.Error(w, "Queue is unavailable", http.StatusServiceUnavailable) 
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"job_id": jobID.String(),
		"status": "pending",
	})
}

func StatusHandler(dbConnection *sql.DB, redisClient *cache.RedisClient, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	limited, err := redisClient.IsRateLimited(r.RemoteAddr) 
	if err != nil {
		log.Println("CACHE ERROR:", err)
		http.Error(w, "Cache service unavailable", http.StatusServiceUnavailable)
		return
	}
	if limited {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}
	
	jobIdStr := chi.URLParam(r, "jobId")
	jobId, err := uuid.Parse(jobIdStr)
	if err != nil {
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		return
	}

	status, err := db.GetStatus(dbConnection, jobId)
	if err != nil {
		http.Error(w, "Failed to get status", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"job_id": jobId.String(),
		"status": status,
	})
}

func GetImageHandler(dbConnection *sql.DB, minioClient *minio.MinioClient, redisClient *cache.RedisClient, w http.ResponseWriter, r *http.Request) {
	
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	limited, err := redisClient.IsRateLimited(r.RemoteAddr) 
	if err != nil {
		log.Println("CACHE ERROR:", err)
		http.Error(w, "Cache service unavailable", http.StatusServiceUnavailable)
		return
	}
	if limited {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}
	
	jobIdStr := chi.URLParam(r, "jobId")
	jobId, err := uuid.Parse(jobIdStr)
	if err != nil {
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		return
	}

	status, err := db.GetStatus(dbConnection, jobId)
	if err != nil {
		http.Error(w, "Failed to get status", http.StatusInternalServerError)
		return
	}
	if status != "completed" {
		http.Error(w, "Image not ready", http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename="+jobId.String()+".zip")

	imgKeys, err := db.GetImageKeys(context.Background(), dbConnection, jobId.String())
	if err != nil { 
		log.Println("DEBUG: Failed to get image keys")
		http.Error(w, "Failed to getting images", http.StatusInternalServerError)
		return
	}
	if imgKeys == nil { 
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte("No images found"))
		return
	}
	log.Println("DEBUG: Found", len(imgKeys), "images")

	zipWriter := zip.NewWriter(w)
	defer zipWriter.Close()
	log.Println("DEBUG: Created ZIP writer")

	for _, key := range imgKeys {
		// get the image from minio
		imageStream, err := minioClient.GetImageStream(context.Background(), "images", key)
		if err != nil {
			log.Printf("Failed to get image for key %s: %v", key, err)
			return
		}
		log.Println("DEBUG: Got image stream for key", key)

		// Now logic for zip handler

		// write the image to the zip file
		writer, err := zipWriter.Create(key)
		if err != nil {
			log.Printf("Failed to create zip entry for key %s: %v", key, err)
			continue
		}
		log.Println("DEBUG: Creating zip entry for key", key)
		
		n, err := io.Copy(writer, imageStream)
		log.Println("DEBUG: Zipped", n, "bytes for key", key)
		imageStream.Close()
		if err != nil {
			return
		}
	}
}
