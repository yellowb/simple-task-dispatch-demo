package iface

import "github.com/yellowb/simple-task-dispatch-demo/internal/app/consumer/model"

type Provider interface {
	// 监听queue，获取消息
	WatchTask()

	// 获取得到的消息chan
	GetTaskChan() chan *model.Job
}
