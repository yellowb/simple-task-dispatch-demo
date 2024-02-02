package impl

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/yellowb/simple-task-dispatch-demo/internal/app/consumer/iface"
	"github.com/yellowb/simple-task-dispatch-demo/internal/app/consumer/model"
)

type DemoWorker struct {
	taskDataStorage iface.TaskDataStorage
	executor        iface.Executor
}

func NewWorker() iface.Worker {
	return &DemoWorker{}
}

func (d *DemoWorker) TaskDataStorage(taskDataStorage iface.TaskDataStorage) iface.Worker {
	d.taskDataStorage = taskDataStorage
	return d
}

func (d *DemoWorker) Executor(executor iface.Executor) iface.Worker {
	d.executor = executor
	return d
}

func (d *DemoWorker) ProcessTask(job *model.Job, handler *model.TaskHandler) {
	ctx := context.Background()
	handler.Args = job.Args
	isSuccess, err := d.executor.Execute(ctx, handler)
	if err != nil {
		logrus.Errorf(err.Error())
	}
	//
	result := ctx.Value("result")
	logsVal := ctx.Value("logs")
	logs, ok := logsVal.([]string)
	if !ok {
		logrus.Errorf("")
		d.taskDataStorage.SaveTaskRecord(job, isSuccess, result, nil)
	}
	d.taskDataStorage.SaveTaskRecord(job, isSuccess, result, logs)
}
