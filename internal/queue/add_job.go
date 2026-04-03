package queue

import (
	"github.com/hibiken/asynq"
)

func AddToQueue(taskDistributor *TaskDistributor, taskID []byte) error {
	task := asynq.NewTask("image:process", taskID)
	_, err := taskDistributor.client.Enqueue(task)
	if err != nil {
		return err
	}
	return nil
}