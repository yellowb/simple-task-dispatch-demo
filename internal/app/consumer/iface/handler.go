package iface

import "github.com/yellowb/simple-task-dispatch-demo/internal/model"

/*
Handler表示真正执行一个Job的业务逻辑
*/

// Handler 业务逻辑处理器接口
type Handler interface {
	// HandlerKey 返回当前Handler对应的HandlerKey
	HandlerKey() string
	// Handle Job处理业务逻辑入口
	Handle(job *model.Job) (*model.JobExecResult, error)
}
