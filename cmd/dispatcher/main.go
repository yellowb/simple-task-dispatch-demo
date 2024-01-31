package main

import (
	"github.com/yellowb/simple-task-dispatch-demo/internal/app/dispatcher/impl"
	"github.com/yellowb/simple-task-dispatch-demo/internal/global"
	"log"
	"time"
)

/**
分派者启动入口
*/

func main() {
	// 配置信息，生产环境应该从环境变量读取，这里写死
	cfg := &global.DispatcherConfig{
		// Redis地址密码等
		RedisCfg: global.RedisCfg{
			Addr:     "localhost",
			Port:     "6379",
			Password: "",
			Db:       0,
		},
		// Redis Queue的名称
		QueueName: "mq",
	}

	/* 先创建依赖的组件 */
	// Job投递者
	deliverier, err := impl.NewRedisQueueDeliverier(cfg.Addr, cfg.Port, cfg.Password, cfg.Db, cfg.QueueName)
	if err != nil {
		log.Fatalf("create redis queue deliverier error : %v", err)
	}
	// Dispatcher状态存储器
	statusStorage := impl.NewMemoryStatusStorage()
	// Task数据源
	taskDatasource := impl.NewFakeTaskDatasource()

	// 初始化Dispatcher
	dispatcher := impl.NewDemoDispatcher().Config(cfg).Deliverier(deliverier).StatusStorage(statusStorage).TaskDatasource(taskDatasource)
	err = dispatcher.Init()
	if err != nil {
		log.Fatalf("init dispatcher error : %v", err)
	}

	// 装载要调度的Task
	err = dispatcher.Reload()
	if err != nil {
		log.Fatalf("dispatcher reloads task error : %v", err)
	}

	// 启动Dispatcher
	_ = dispatcher.Run()

	// 简单不让进程立即退出...
	select {
	case <-time.After(24 * time.Hour):
	}

	// TODO: 还需要监听os事件调用dispatcher.Shutdown()来退出
}
