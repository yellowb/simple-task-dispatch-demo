package impl

import "github.com/yellowb/simple-task-dispatch-demo/internal/app/consumer/iface"

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

func (d *DemoWorker) ProcessTask() {
	//TODO implement me
	panic("implement me")
}
