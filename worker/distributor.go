package worker

import (
	"context"

	"github.com/hibiken/asynq"
)

//create tasks and distribute them to worker via redis queue

type TaskDistributor interface {
	DistributeTaskSendVerifyEmail(
		ctx context.Context,
		payload *PayloadSendVerifyEmail,
		opts ...asynq.Option,
	) error
}

type RedisTaskDistributor struct {
	client *asynq.Client //use to send task to redis
}

func NewRedisDistributor(redisOpt asynq.RedisClientOpt) TaskDistributor {
	client := asynq.NewClient(redisOpt)
	return &RedisTaskDistributor{ //by returning RedisTaskDistributor we are forcing RedisTaskDistributor to implement TaskDistributor
		client: client,
	}
}
