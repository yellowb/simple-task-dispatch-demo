package model

import "context"

// 任务执行函数
type TaskExecFunc func(ctx context.Context)

// 任务执行结构体
type TaskHandler struct {
	Args     map[string]interface{} //执行参数
	ExecFunc TaskExecFunc           //执行函数
	Result   interface{}            //执行结果
	Logs     []string               //执行输出日志
}
