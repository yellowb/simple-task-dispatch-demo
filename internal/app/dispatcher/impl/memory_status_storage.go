package impl

import (
	"github.com/yellowb/simple-task-dispatch-demo/internal/app/dispatcher/model"
	"sync"
)

// MemoryStatusStorage 基于内存的Dispatcher内部状态存储器，只能用于单进程
type MemoryStatusStorage struct {
	// 存储Task运行时状态数据，key = taskKey
	taskMap sync.Map
}

func (m *MemoryStatusStorage) PutRunningTaskStatus(taskKey string, taskStatus *model.RunningTaskStatus) error {
	//TODO implement me
	panic("implement me")
}

func (m *MemoryStatusStorage) DeleteRunningTaskStatus(taskKey string) error {
	//TODO implement me
	panic("implement me")
}

func (m *MemoryStatusStorage) GetRunningTaskStatus(taskKey string) (*model.RunningTaskStatus, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MemoryStatusStorage) GetDispatcherStatus() (model.DispatcherStatus, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MemoryStatusStorage) Clear() error {
	//TODO implement me
	panic("implement me")
}

func (m *MemoryStatusStorage) Size() int {
	//TODO implement me
	panic("implement me")
}
