package iface

import "github.com/yellowb/simple-task-dispatch-demo/internal/global"

/**
Consumer表示Job的消费者
*/

// Consumer Job消费者接口
type Consumer interface {
	/* 下面这些是Consumer启动前的配置接口 */

	// Config 注入配置到Consumer
	Config(cfg *global.ConsumerConfig) Consumer
	// Receiver 注入Job的接收器到Consumer
	Receiver(receiver Receiver)

	/* 下面这些是影响Consumer运行状态的接口 */

	// Init 初始化Consumer
	Init() error
	// Run 启动Consumer
	Run() error
	// Shutdown 终止Consumer
	Shutdown() error
}
