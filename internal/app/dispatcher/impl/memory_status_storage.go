package impl

import (
	model2 "github.com/yellowb/simple-task-dispatch-demo/internal/model"
	"sync"
)

// MemoryStatusStorage 基于内存的Dispatcher内部状态存储器，只能用于单进程
type MemoryStatusStorage struct {
	// 存储Task运行时状态数据，key = taskKey
	taskMap sync.Map
}

func NewMemoryStatusStorage() *MemoryStatusStorage {
	return &MemoryStatusStorage{
		taskMap: sync.Map{},
	}
}

func (m *MemoryStatusStorage) PutRunningTaskStatus(taskKey string, taskStatus *model2.RunningTaskStatus) error {
	m.taskMap.Store(taskKey, taskStatus)
	return nil
}

func (m *MemoryStatusStorage) DeleteRunningTaskStatus(taskKey string) error {
	m.taskMap.Delete(taskKey)
	return nil
}

func (m *MemoryStatusStorage) GetRunningTaskStatus(taskKey string) (*model2.RunningTaskStatus, error) {
	value, ok := m.taskMap.Load(taskKey)
	if !ok {
		return nil, nil
	}
	return value.(*model2.RunningTaskStatus), nil
}

func (m *MemoryStatusStorage) ExistRunningTaskStatus(taskKey string) (bool, error) {
	_, ok := m.taskMap.Load(taskKey)
	return ok, nil
}

func (m *MemoryStatusStorage) GetDispatcherStatus() (*model2.DispatcherStatus, error) {
	data := make(map[string]*model2.RunningTaskStatus)
	m.taskMap.Range(func(key, value any) bool {
		data[key.(string)] = value.(*model2.RunningTaskStatus)
		return true
	})
	return &model2.DispatcherStatus{
		RunningTaskCount:    len(data),
		RunningTaskStatuses: data,
	}, nil
}

func (m *MemoryStatusStorage) Clear() error {
	m.taskMap.Range(func(key, value any) bool {
		m.taskMap.Delete(key)
		return true
	})
	return nil
}

func (m *MemoryStatusStorage) Size() int {
	counter := 0
	// 因为sync.Map没有接口可以获取当前Size，这里用了一种效率较低的方式循环迭代每一个key来统计个数
	// 一开始考虑使用一个外置的atomic.Int32作为Size counter，但是发现这种做法在Put/Delete的场景下会有并发问题，于是舍弃
	m.taskMap.Range(func(key, value any) bool {
		counter++
		return true
	})
	return counter
}
