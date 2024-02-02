package impl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/yellowb/simple-task-dispatch-demo/internal/app/dispatcher/model"
	"github.com/yellowb/simple-task-dispatch-demo/internal/global"
)

type RedisProvider struct {
	redisCli  *redis.Client
	queueName string
	taskChan  chan *model.Job
}

func (r *RedisProvider) WatchTask() {
	// 循环获取消息
	for {
		// 从Redis队列中获取消息,没有消息会阻塞等待
		result, err := r.redisCli.BLPop(context.Background(), 0, r.queueName).Result()
		if err != nil {
			logrus.Errorf("getting message from Redis error :%v", err)
			continue
		}
		// 提取队列名和消息内容
		queue := result[0]
		message := result[1]
		fmt.Printf("message received from queue '%s': %s\n", queue, message)
		job := &model.Job{}
		err = json.Unmarshal([]byte(message), &job)
		if err != nil {
			logrus.Errorf("unmarshal message from Redis error:%v", err)
			continue
		}
		// 将消息放入chan中处理
		r.taskChan <- job
		fmt.Println("job deliver to taskChan")
	}
}

func (r *RedisProvider) GetTaskChan() chan *model.Job {
	return r.taskChan
}

func NewRedisQueueProvider(mq string) (*RedisProvider, error) {
	if global.GetRedisCli() == nil {
		return nil, errors.New("redis client error")
	}
	return &RedisProvider{
		redisCli:  global.GetRedisCli(),
		queueName: mq,
	}, nil
}
