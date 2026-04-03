package queue

import (
	"context"
	"encoding/json"
	"github.com/hibiken/asynq"
)

func ProcessImage(ctx context.Context, task *asynq.Task) error {
	marshalledPayload := task.Payload()
	var payload ImageProcessPayload
	err := json.Unmarshal(marshalledPayload, &payload)
	if err != nil {
		return err
	}
	
	println("Processing job with data:", payload.JobID, payload.URL)
	return nil
}
