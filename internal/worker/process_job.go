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
		return err
	}

	fmt.Println("Beginning processing image with job ID:", payload.JobID)
	// write to DB saying its in progress
	_, err = workerHandler.db.ExecContext(ctx, "UPDATE jobs SET status = 'processing' WHERE id = $1", payload.JobID)
	if err != nil {
		fmt.Println("Error updating job status:", err)
		return err
	} 

	fmt.Println("Job status updated to processing")

	// get photo from minio  
	data, err := workerHandler.minioClient.GetImage(ctx, "images", payload.Key)
	if err != nil {
		fmt.Println("Error getting object:", err)
		return err
	} 

	fmt.Println("Got image data:", len(data), "bytes")


	// we need to decode the image 

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil { 
		fmt.Println("Error decoding image:", err)
		return err
	}
	fmt.Println("Decoded image:", img.Bounds())

	//loop of 3 
	sizes := []int{200, 800, 1600}
	for _, size := range sizes {
		resizedImg, err := ResizeImage(img, size, 0)
		if err != nil {
			return err
		}
		fmt.Printf("Resized image to %dx%d (%d bytes)\n", size, size, len(resizedImg))

		ioImg := bytes.NewReader(resizedImg)

		key, err := minio.WriteImage(workerHandler.minioClient, ioImg, int64(len(resizedImg)))
		if err != nil {
			return err
		}
		fmt.Printf("Saved image to minio with key: %s\n", key)

		err = db.CreateImage(ctx, workerHandler.db, payload.JobID, key, getSizeName(size))
		if err != nil {
			fmt.Println("Error creating image in the database:", err)
			return err
		}
		fmt.Printf("Written image to database with key: %s\n", key)

		
		// write to images table
	}
	
	// update job status to finsihed in jobs table 

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