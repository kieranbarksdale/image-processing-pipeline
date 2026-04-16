package worker

import (
	"context"
	"encoding/json"
	"bytes"
	"image"
	_ "image/jpeg"
	_ "image/png"
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

	err = db.PutStatus(ctx, workerHandler.db, payload.JobID, "processing")
	if err != nil {
		log.Printf("Error updating job status: %v", err)
		db.IncrementTries(ctx, workerHandler.db, payload.JobID)
		return err
	}

	// get photo from minio  
	imageStream, err := workerHandler.minioClient.GetImageStream(ctx, "images", payload.Key)
	if err != nil {
		log.Printf("Error getting object: %v", err)
		db.IncrementTries(ctx, workerHandler.db, payload.JobID)
		return err
	} 

	defer imageStream.Close()

	// we need to decode the image
	img, format, err := image.Decode(imageStream)
	if err != nil { 
		log.Printf("Error decoding image: %v", err)
		db.IncrementTries(ctx, workerHandler.db, payload.JobID)
		return err
	}

	//loop of all the sizes
	sizes := []int{200, 800, 1600}
	for _, size := range sizes {
		resizedImg, err := ResizeImage(img, size, 0)
		if err != nil {
			db.IncrementTries(ctx, workerHandler.db, payload.JobID)
			return err
		}

		ioImg := bytes.NewReader(resizedImg)

		key, err := minio.WriteImage(workerHandler.minioClient, ioImg, int64(len(resizedImg)))
		if err != nil {
			db.IncrementTries(ctx, workerHandler.db, payload.JobID)
			return err
		}

		err = db.CreateImage(ctx, workerHandler.db, payload.JobID, key, getSizeName(size))
		if err != nil {
			log.Printf("Error creating image in the database: %v", err)
			db.IncrementTries(ctx, workerHandler.db, payload.JobID)
			return err
		}
	}
	
	err = db.PutStatus(ctx, workerHandler.db, payload.JobID, "completed")
	if err != nil {
		log.Printf("Error updating job status: %v", err)
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