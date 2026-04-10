package queue

import (
	"encoding/json"
	"github.com/hibiken/asynq"
)

type ImageProcessPayload struct {
	JobID string `json:"job_id"`
	Key   string `json:"key"`
}

func (taskDistributor *TaskDistributor) AddToQueue(jobId string, originalKey string) error {
	payload := ImageProcessPayload{
		JobID: jobId,
		Key:   originalKey,
	}
	marshalledPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	task := asynq.NewTask("image:process", marshalledPayload)
	_, err = taskDistributor.client.Enqueue(task, asynq.MaxRetry(3))
	if err != nil {
		return err
	}
	return nil
}