package iface

import (
	"github.com/yellowb/simple-task-dispatch-demo/internal/global"
	"github.com/yellowb/simple-task-dispatch-demo/internal/model"
)

type Consumer interface {
	/* Consumer配置相关方法 */

	// Config 注入配置到Consumer
	Config(cfg *global.ConsumerConfig) Consumer
	//Receiver
	Receiver(receiver Receiver) Consumer

	//初始化
	Init() error
	//启动
	Run() error
	//关闭
	Shutdown() error

	// 任务执行器
	GetJobExecutor() <-chan *JobExecutor

	//任务映射
	GetMappingExecutor(job *model.Job) *JobExecutor
}

// JobExecutor 从队列获取的任务和执行器
type JobExecutor struct {
	model.Job
	Executor Executor
}
