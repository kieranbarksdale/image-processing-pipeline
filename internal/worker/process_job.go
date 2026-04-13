package worker

import (
	"context"
	"encoding/json"
	"bytes"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"fmt"
	"image-processing-pipeline/internal/queue"
	"github.com/hibiken/asynq"
	"image-processing-pipeline/internal/db"
	"image-processing-pipeline/internal/minio"
)

func (workerHandler *WorkerHandler) HandleProcessJob(ctx context.Context, task *asynq.Task) error {

	marshalledPayload := task.Payload()
	var payload queue.ImageProcessPayload
	err := json.Unmarshal(marshalledPayload, &payload)
	if err != nil {
		db.IncrementTries(ctx, workerHandler.db, payload.JobID)
		return err
	}

	fmt.Println("Beginning processing image with job ID:", payload.JobID)

	err = db.PutStatus(ctx, workerHandler.db, payload.JobID, "processing")
	if err != nil {
		fmt.Println("Error updating job status:", err)
		db.IncrementTries(ctx, workerHandler.db, payload.JobID)
		return err
	} 

	fmt.Println("Job status updated to processing")

	// get photo from minio  
	data, err := workerHandler.minioClient.GetImage(ctx, "images", payload.Key)
	if err != nil {
		fmt.Println("Error getting object:", err)
		db.IncrementTries(ctx, workerHandler.db, payload.JobID)
		return err
	} 

	fmt.Println("Got image data:", len(data), "bytes")


	// we need to decode the image 
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil { 
		fmt.Println("Error decoding image:", err)
		db.IncrementTries(ctx, workerHandler.db, payload.JobID)
		return err
	}
	fmt.Println("Decoded image:", img.Bounds())

	//loop of all the sizes
	sizes := []int{200, 800, 1600}
	for _, size := range sizes {
		resizedImg, err := ResizeImage(img, size, 0)
		if err != nil {
			db.IncrementTries(ctx, workerHandler.db, payload.JobID)
			return err
		}
		fmt.Printf("Resized image to %dx%d (%d bytes)\n", size, size, len(resizedImg))

		ioImg := bytes.NewReader(resizedImg)

		key, err := minio.WriteImage(workerHandler.minioClient, ioImg, int64(len(resizedImg)))
		if err != nil {
			db.IncrementTries(ctx, workerHandler.db, payload.JobID)
			return err
		}
		fmt.Printf("Saved image to minio with key: %s\n", key)

		err = db.CreateImage(ctx, workerHandler.db, payload.JobID, key, getSizeName(size))
		if err != nil {
			fmt.Println("Error creating image in the database:", err)
			db.IncrementTries(ctx, workerHandler.db, payload.JobID)
			return err
		}
		fmt.Printf("Written image to database with key: %s\n", key)
	}
	
	err = db.PutStatus(ctx, workerHandler.db, payload.JobID, "completed")
	if err != nil {
		fmt.Println("Error updating job status:", err)
		db.IncrementTries(ctx, workerHandler.db, payload.JobID)
		return err
	}
	
	return nil
}

func getSizeName(size int) string {
	switch size {
	case 200:
		return "thumbnail"
	case 800:
		return "medium"
	case 1600:
		return "large"
	default:
		return "unknown"
	}
}