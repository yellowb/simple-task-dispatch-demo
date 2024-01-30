package iface

import "github.com/yellowb/simple-task-dispatch-demo/internal/app/dispatcher/model"

// StatusStorage Dispatcher内部状态存储器接口
type StatusStorage interface {
	// PutRunningTaskStatus 更新某一个Task的运行时状态
	PutRunningTaskStatus(taskKey string, taskStatus *model.RunningTaskStatus) error
	// DeleteRunningTaskStatus 删除某一个Task的运行时状态
	DeleteRunningTaskStatus(taskKey string) error
	// GetRunningTaskStatus 获取某一个Task的运行时状态
	GetRunningTaskStatus(taskKey string) (*model.RunningTaskStatus, error)
	// GetDispatcherStatus 获取Dispatcher完整运行时状态
	GetDispatcherStatus() (model.DispatcherStatus, error)
	// Clear 清空Dispatcher内部状态数据
	Clear() error
	// Size 当前Task数量
	Size() int
}
