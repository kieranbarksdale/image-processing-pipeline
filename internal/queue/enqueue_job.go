package queue

import (
	"encoding/json"

	"github.com/hibiken/asynq"
)

type ImageProcessPayload struct {
	JobID string `json:"job_id"`
	URL   string `json:"url"`
}

func (taskDistributor *TaskDistributor) AddToQueue(jobId string, originalURL string) error {
	payload := ImageProcessPayload{
		JobID: jobId,
		URL:   originalURL,
	}
	marshalledPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	task := asynq.NewTask("image:process", marshalledPayload)
	_, err = taskDistributor.client.Enqueue(task)
	if err != nil {
		return err
	}
	return nil
}