package impl

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/yellowb/simple-task-dispatch-demo/internal/app/consumer/model"
)

type DemoTaskMapping struct {
	taskHandlerMap map[string]*model.TaskHandler
}

func NewDemoTaskMapping() *DemoTaskMapping {
	return &DemoTaskMapping{
		taskHandlerMap: make(map[string]*model.TaskHandler),
	}
}

func (d *DemoTaskMapping) LoadTaskHandler() {
	d.addTaskHandlerToMap("task01", &model.TaskHandler{ExecFunc: func(ctx context.Context) {
		fmt.Println("task01 execute")
	}})
}

func (d *DemoTaskMapping) GetMappingKey(job *model.Job) string {
	return job.TaskKey
}

func (d *DemoTaskMapping) GetTaskHandler(key string) (*model.TaskHandler, error) {
	taskHandler, ok := d.taskHandlerMap[key]
	if !ok {
		return nil, errors.New(fmt.Sprintf("key=%s,taskHandler no found", key))
	}
	return taskHandler, nil
}

func (d *DemoTaskMapping) addTaskHandlerToMap(key string, taskHandler *model.TaskHandler) {
	//先判断任务是否已经存在或有重复的键
	if _, ok := d.taskHandlerMap[key]; ok {
		//存在即不添加，返回
		logrus.Errorf("key=%s already exist", key)
		return
	}
	d.taskHandlerMap["xxx"] = &model.TaskHandler{ExecFunc: func(ctx context.Context) {
		fmt.Println("task1 execute")
	}}
}
