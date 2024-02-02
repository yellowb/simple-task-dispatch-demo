package impl

import (
	"github.com/yellowb/simple-task-dispatch-demo/internal/constants"
	"github.com/yellowb/simple-task-dispatch-demo/internal/model"
	"time"
)

// FakeTaskDatasource 一个假的Task数据源实现，总是返回写死的数据
type FakeTaskDatasource struct {
	// 正常的实现应该要包含DAO从DB中获取数据
}

func NewFakeTaskDatasource() *FakeTaskDatasource {
	return &FakeTaskDatasource{}
}

func (p *FakeTaskDatasource) GetAllTasks() ([]*model.Task, error) {
	return []*model.Task{
		&model.Task{
			Name:        "task-5s",
			Key:         "task-a",
			TaskType:    constants.TaskTypeDuration,
			DurationDef: &model.DurationDefinition{Duration: 5 * time.Second},
		},
		&model.Task{
			Name:        "task-10s",
			Key:         "task-b",
			TaskType:    constants.TaskTypeDuration,
			DurationDef: &model.DurationDefinition{Duration: 10 * time.Second},
		},
	}, nil
}
