package iface

import "github.com/yellowb/simple-task-dispatch-demo/internal/model"

type Executor interface {
	//任务执行
	Execute(args map[string]interface{}) (*model.JobResult, error)
}
