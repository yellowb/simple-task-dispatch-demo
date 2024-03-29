package impl

import (
	"github.com/yellowb/simple-task-dispatch-demo/internal/app/dispatcher/model"
	"github.com/yellowb/simple-task-dispatch-demo/internal/constants"
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
			Key:         "task-5s",
			TaskType:    constants.TaskTypeDuration,
			DurationDef: &model.DurationDefinition{Duration: 5 * time.Second},
		},
		&model.Task{
			Name:        "task-10s",
			Key:         "task-10s",
			TaskType:    constants.TaskTypeDuration,
			DurationDef: &model.DurationDefinition{Duration: 10 * time.Second},
		},
	}, nil
}
