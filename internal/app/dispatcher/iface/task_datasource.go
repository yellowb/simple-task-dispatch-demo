package iface

import (
	"github.com/yellowb/simple-task-dispatch-demo/internal/model"
)

// TaskDatasource 提供Task数据的源头
type TaskDatasource interface {
	// GetAllTasks 获取所有Task
	GetAllTasks() ([]*model.Task, error)
}
