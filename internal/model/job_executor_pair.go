package model

import "github.com/yellowb/simple-task-dispatch-demo/internal/app/consumer/iface"

// JobExecutorPair 一个Job实例与之对应Executor实例的组合
type JobExecutorPair struct {
	Job      *Job
	Executor iface.Executor
}
