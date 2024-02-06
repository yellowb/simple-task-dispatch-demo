package main

import (
	"github.com/yellowb/simple-task-dispatch-demo/internal/app/consumer/impl"
	"github.com/yellowb/simple-task-dispatch-demo/internal/global"
	"time"
)

/*
*
消费者启动入口
*/
func main() {
	// receiver
	receiverCfg := &global.ReceiverConfig{
		// Redis地址密码等
		RedisCfg: global.RedisCfg{
			Addr:     "localhost",
			Port:     "6379",
			Password: "",
			Db:       0,
		},
		// Redis Queue的名称
		QueueName:   "mq",
		JobChanSize: 100,
	}

	global.RedisInit(receiverCfg.RedisCfg)
	global.RedSyncLockInit()

	receiver := impl.NewRedisQueueReceiver()
	receiver.Config(receiverCfg)
	_ = receiver.Init()
	_ = receiver.Run()

	//Consumer
	consumerCfg := &global.ConsumerConfig{
		WorkerPoolSize:      5,
		JobExecutorChanSize: 100,
	}
	consumer := impl.NewConsumer()
	consumer.Config(consumerCfg).Receiver(receiver)
	_ = consumer.Init()
	_ = consumer.Run()

	// worker
	taskResultStorage := impl.NewMongoTaskDataStorage()
	for i := 0; i < consumerCfg.WorkerPoolSize; i++ {
		worker := impl.NewWorker(i)
		worker.Consumer(consumer).TaskResultStorage(taskResultStorage)
		_ = worker.Init()
		_ = worker.Run()
	}
	select {
	case <-time.After(24 * time.Hour):

	}
}
