package iface

import (
	"github.com/yellowb/simple-task-dispatch-demo/internal/global"
	"github.com/yellowb/simple-task-dispatch-demo/internal/model"
)

type Receiver interface {
	// Config 注入配置到Consumer
	Config(cfg *global.ReceiverConfig) Receiver

	//初始化
	Init() error
	//启动
	Run() error
	//关闭
	Shutdown() error

	// 获取得到的消息chan
	GetJobChan() <-chan *model.Job
}
