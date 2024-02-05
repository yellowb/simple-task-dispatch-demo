package main

import (
	"github.com/yellowb/simple-task-dispatch-demo/internal/app/consumer/impl"
	"github.com/yellowb/simple-task-dispatch-demo/internal/global"
	"log"
	"time"
)

/**
消费者启动入口
*/

func main() {
	// 配置信息
	receiverCfg := &global.ReceiverConfig{
		// Redis地址密码等
		RedisCfg: global.RedisCfg{
			Addr:     "localhost",
			Port:     "6379",
			Password: "",
			Db:       0,
		},
		// Redis Queue的名称
		QueueName:     "mq",
		JobBufferSize: 10 * 1024,
	}
	consumerCfg := &global.ConsumerConfig{
		WorkerPoolSize: 2,
		JobBufferSize:  10 * 1024,
	}

	/* 先创建依赖的组件 */
	// Job Receiver
	receiver := impl.NewRedisQueueReceiver().Config(receiverCfg)
	err := receiver.Init()
	if err != nil {
		log.Fatalf("init receiver error : %v", err)
	}
	err = receiver.Run()
	if err != nil {
		log.Fatalf("start receiver error : %v", err)
	}
	// Task Result Storage
	storage := new(impl.FakeTaskResultStorage)

	// 初始化Consumer
	consumer := impl.NewDemoConsumer().Config(consumerCfg).Receiver(receiver).TaskResultStorage(storage)
	err = consumer.Init()
	if err != nil {
		log.Fatalf("init consumer error : %v", err)
	}

	// 启动Consumer
	err = consumer.Run()
	if err != nil {
		log.Fatalf("start consumer error : %v", err)
	}

	// 简单不让进程立即退出...
	select {
	case <-time.After(24 * time.Hour):
	}

	// TODO: 还需要监听os事件调用dispatcher.Shutdown()来退出, 还有断开数据库链接之类。。。
}
