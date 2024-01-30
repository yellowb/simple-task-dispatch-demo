package dao

import "github.com/yellowb/simple-task-dispatch-demo/internal/app/dispatcher/model"

// FakeTaskDao 一个假的Task DAO
type FakeTaskDao struct {
}

var fakeTasks = []*model.Task{
	&model.Task{
		Name: "task-10s",
		Key:  "task-10s",
	},
}
