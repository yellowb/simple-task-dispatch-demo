package iface

import (
	"github.com/yellowb/simple-task-dispatch-demo/internal/global"
)

type Consumer interface {
	/* Consumer配置相关方法 */

	// Config 注入配置到Consumer
	Config(cfg *global.ConsumerConfig) Consumer
	//注入Provider
	Provider(provider Provider) Consumer
	//注入worker
	Worker(worker Worker) Consumer
	//注入TaskMapping
	TaskMapping(taskMapping TaskMapping) Consumer

	// 消费消息
	ConsumeTask()
}
