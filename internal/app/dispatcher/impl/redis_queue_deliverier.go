package impl

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/yellowb/simple-task-dispatch-demo/internal/model"
	"time"
)

// RedisQueueDeliverier 对接Redis List的Job投递器实现
type RedisQueueDeliverier struct {
	redisCli  *redis.Client
	queueName string // Key of redis list
}

func NewRedisQueueDeliverier(addr, port, password string, db int, queueName string) (*RedisQueueDeliverier, error) {
	cli := redis.NewClient(&redis.Options{
		Addr:        fmt.Sprintf("%s:%s", addr, port),
		Password:    password,
		DB:          db,
		ReadTimeout: 10 * time.Second, // go-redis默认值为3s, 可能在BRPOP中堵塞超时设置稍微大一点就会有问题，所以调大它
	})
	_, err := cli.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	return &RedisQueueDeliverier{
		redisCli:  cli,
		queueName: queueName,
	}, nil
}

func (r *RedisQueueDeliverier) Deliver(job *model.Job) error {
	jsonBytes, _ := json.Marshal(job) // demo简单起见忽略了error
	// 简单的把job序列化成json字符串后，往List中append
	_, err := r.redisCli.LPush(context.Background(), r.queueName, string(jsonBytes)).Result()
	return err
}

func (r *RedisQueueDeliverier) Len() (int64, error) {
	return r.redisCli.LLen(context.Background(), r.queueName).Result()
}
