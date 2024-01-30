package impl

import (
	"github.com/yellowb/simple-task-dispatch-demo/internal/app/dispatcher/model"
)

// RedisQueueDeliverier 对接Redis List的Job投递器实现
type RedisQueueDeliverier struct {
	// TODO：redis client
}

func (r *RedisQueueDeliverier) Deliver(job *model.Job) error {
	//TODO implement me
	panic("implement me")
}
