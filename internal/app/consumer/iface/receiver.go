package iface

import (
	"github.com/yellowb/simple-task-dispatch-demo/internal/global"
	"github.com/yellowb/simple-task-dispatch-demo/internal/model"
)

/**
Receiver代表Job的接收器，可以从任意地方接收Job
*/

// Receiver Job接收器接口
type Receiver interface {
	// Config 把配置信息注入Receiver
	Config(cfg *global.ReceiverConfig) Receiver
	// Init 初始化Receiver
	Init() error
	// Run 启动Receiver
	Run() error
	// Stop 停止Receiver
	Stop() error
	// Shutdown 关闭Receiver
	Shutdown() error
	// GetJobChannel 返回一个channel，后续使用者可以从这个channel中拿到Receiver接收到的Job
	GetJobChannel() <-chan *model.Job
}
