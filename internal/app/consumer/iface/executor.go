package iface

import "github.com/yellowb/simple-task-dispatch-demo/internal/model"

/*
Executor表示真正执行一个Job的业务逻辑
*/

// Executor 业务逻辑处理器接口
type Executor interface {
	// ExecutorKey 返回当前Executor对应的ExecutorKey
	ExecutorKey() string
	// Execute Job处理业务逻辑入口
	Execute(job *model.Job) (*model.JobExecResult, error)
}
