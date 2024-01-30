package iface

import "github.com/yellowb/simple-task-dispatch-demo/internal/app/dispatcher/model"

// TaskProvider 提供Task数据的源头
type TaskProvider interface {
	// GetAllTasks 获取所有Task
	GetAllTasks() ([]*model.Task, error)
}
