package impl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/yellowb/simple-task-dispatch-demo/internal/app/consumer/iface"
	"github.com/yellowb/simple-task-dispatch-demo/internal/app/consumer/status"
	"github.com/yellowb/simple-task-dispatch-demo/internal/global"
	"github.com/yellowb/simple-task-dispatch-demo/internal/model"
	"sync"
	"sync/atomic"
)

type RedisReceiver struct {
	cfg       *global.ReceiverConfig
	redisCli  *redis.Client
	queueName string
	jobChan   chan *model.Job
	isStop    atomic.Bool
	status    status.Status
	lock      sync.Mutex
}

func NewRedisQueueReceiver() *RedisReceiver {
	return &RedisReceiver{
		status: status.New,
		lock:   sync.Mutex{},
	}
}

func (r *RedisReceiver) WatchTask() {
	// 循环获取消息
	for !r.isStop.Load() {
		// 从Redis队列中获取消息,没有消息会阻塞等待
		result, err := r.redisCli.BRPop(context.Background(), 0, r.queueName).Result()
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
		r.jobChan <- job
		fmt.Println("job deliver to taskChan")
	}
	r.isStop.Store(false)
}

func (r *RedisReceiver) Config(cfg *global.ReceiverConfig) iface.Receiver {
	r.cfg = cfg
	return r
}

func (r *RedisReceiver) Init() error {
	r.lock.Lock()
	defer r.lock.Unlock()
	err := status.CheckStatus(r.status, status.New)
	if err != nil {
		return err
	}
	if r.cfg == nil {
		return errors.New("cfg is nil")
	}
	r.redisCli = global.GetRedisCli()
	if r.redisCli == nil {
		return errors.New("redis client is nil, connect error")
	}
	r.queueName = r.cfg.QueueName
	r.jobChan = make(chan *model.Job, r.cfg.JobChanSize)
	r.status = status.Initialized
	return nil
}

func (r *RedisReceiver) Run() error {
	r.lock.Lock()
	defer r.lock.Unlock()
	err := status.CheckStatus(r.status, status.Initialized)
	if err != nil {
		return err
	}
	go r.WatchTask()
	r.status = status.Running
	return nil
}

func (r *RedisReceiver) Shutdown() error {
	r.lock.Lock()
	defer r.lock.Unlock()
	err := status.CheckStatus(r.status, status.Running)
	if err != nil {
		return err
	}
	r.status = status.Stopped
	r.isStop.Store(true)
	return nil
}

func (r *RedisReceiver) GetJobChan() <-chan *model.Job {
	return r.jobChan
}
