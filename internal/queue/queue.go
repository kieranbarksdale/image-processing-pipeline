package queue

import (
	"github.com/hibiken/asynq"
)

type TaskDistributor struct {
	client *asynq.Client
}

func CreateQueue(redisAddr string) *TaskDistributor {
	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr: redisAddr,
	})
	return &TaskDistributor{
		client: client,
	}
}
