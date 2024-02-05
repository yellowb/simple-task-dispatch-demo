package iface

import (
	"github.com/yellowb/simple-task-dispatch-demo/internal/model"
)

// JobExecutorPair 一个Job实例与之对应Executor实例的组合
type JobExecutorPair struct {
	Job      *model.Job
	Executor Executor
}
