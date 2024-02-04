package impl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/yellowb/simple-task-dispatch-demo/internal/app/consumer/iface"
	"github.com/yellowb/simple-task-dispatch-demo/internal/app/consumer/receiver_status"
	"github.com/yellowb/simple-task-dispatch-demo/internal/global"
	"github.com/yellowb/simple-task-dispatch-demo/internal/model"
	"log"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	defaultJobChLen        = 10 * 1024
	defaultRedisPopTimeout = 2 * time.Second
)

type RedisQueueReceiver struct {
	// 配置信息
	cfg *global.ReceiverConfig

	// 内部属性
	status    receiver_status.ReceiverStatus // Receiver状态
	lock      sync.Mutex
	redisCli  *redis.Client
	queueName string          // Key of redis list
	jobCh     chan *model.Job // 从Redis Queue接收后的Job缓冲在这个channel等待Consumer取走
	stopFlag  atomic.Bool     // 跳出从Redis Queue获取数据的循环的标记位
}

func NewRedisQueueReceiver() *RedisQueueReceiver {
	return &RedisQueueReceiver{
		status: receiver_status.New,
		lock:   sync.Mutex{},
	}
}

func (r *RedisQueueReceiver) Config(cfg *global.ReceiverConfig) iface.Receiver {
	r.cfg = cfg
	return r
}

func (r *RedisQueueReceiver) Init() error {
	r.lock.Lock()
	defer r.lock.Unlock()

	// 禁止重复初始化
	err := r.checkStatus(receiver_status.New)
	if err != nil {
		return err
	}

	// 简单检查一下依赖项。正式代码中应该对每个字段，比如redis addr之类都做检查
	if r.cfg == nil {
		return errors.New("config is nil")
	}
	r.queueName = r.cfg.QueueName

	// 初始化 redis client
	cli := redis.NewClient(&redis.Options{
		Addr:        fmt.Sprintf("%s:%s", r.cfg.Addr, r.cfg.Port),
		Password:    r.cfg.Password,
		DB:          r.cfg.Db,
		ReadTimeout: 10 * time.Second, // go-redis默认值为3s, 可能在BRPOP中堵塞超时设置稍微大一点就会有问题，所以调大它
	})
	_, err = cli.Ping(context.Background()).Result()
	if err != nil {
		return fmt.Errorf("init redis client error : %v", err)
	}

	// init job channel
	jobChLen := r.cfg.JobBufferSize
	if jobChLen <= 0 {
		jobChLen = defaultJobChLen
	}
	r.jobCh = make(chan *model.Job, jobChLen)

	r.stopFlag.Store(false)

	r.status = receiver_status.Initialized
	return nil
}

func (r *RedisQueueReceiver) Run() error {
	r.lock.Lock()
	defer r.lock.Unlock()

	// 状态检查
	err := r.checkStatus(receiver_status.Initialized, receiver_status.Stopped)
	if err != nil {
		return err
	}

	// 起一个新协程循环从Redis Queue中读取数据，然后投递到jobCh中
	go func(receiver *RedisQueueReceiver) {
		for !receiver.stopFlag.Load() {
			v, err := receiver.redisCli.BRPop(context.Background(), defaultRedisPopTimeout, receiver.queueName).Result()
			if err != nil {
				// redis.Nil不用管
				if !errors.Is(err, redis.Nil) {
					// 遇到错误了，等一会再试试
					log.Printf("[ERROR] pop message from redis error : %v", err)
					time.Sleep(time.Second)
				}
			}

			job := new(model.Job)
			err = json.Unmarshal([]byte(v[1]), job)
			if err != nil {
				// 坏数据，直接丢弃
				log.Printf("[ERROR] parse redis queue message data error : %v", err)
			}

			r.jobCh <- job
		}
		receiver.stopFlag.Store(false) // 重制标志位
		log.Printf("[RedisQueueReceiver] loop exited")
	}(r)

	r.status = receiver_status.Running
	return nil
}

func (r *RedisQueueReceiver) Stop() error {
	r.lock.Lock()
	defer r.lock.Unlock()

	// 状态检查
	err := r.checkStatus(receiver_status.Running)
	if err != nil {
		return err
	}
	r.status = receiver_status.Stopped

	r.stopFlag.Store(true)
	return nil
}

func (r *RedisQueueReceiver) GetJobChannel() <-chan *model.Job {
	return r.jobCh
}

func (r *RedisQueueReceiver) checkStatus(statusList ...receiver_status.ReceiverStatus) error {
	for _, status := range statusList {
		if r.status == status {
			return nil
		}
	}

	// 构造错误信息
	statusStrList := make([]string, 0, len(statusList))
	for _, status := range statusList {
		statusStrList = append(statusStrList, status.String())
	}
	return fmt.Errorf("receiver status must be in %s, current status: %s", strings.Join(statusStrList, "/"), r.status.String())
}
