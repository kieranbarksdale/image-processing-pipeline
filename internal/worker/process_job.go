package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"image-processing-pipeline/internal/queue"
	"github.com/minio/minio-go/v7"
	"github.com/hibiken/asynq"
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
	_, err = workerHandler.db.ExecContext(ctx, "UPDATE jobs SET status = 'processing' WHERE id = ?", payload.JobID)
	if err != nil {
		fmt.Println("Error updating job status:", err)
		return err
	} 

	// get photo from minio  
	obj, err := workerHandler.minioClient.GetObject(ctx, "images",payload.URL, minio.GetObjectOptions{})
	if err != nil {
		fmt.Println("Error getting object:", err)
		return err
	}
	defer obj.Close()

	//loop of 3 
	for i := 0; i < 3; i++ {
		// do something with that photo
		// save result to minio 
	}
	
	// update job status to finsihed 
	// return nil

	return nil
}