package iface

import "github.com/yellowb/simple-task-dispatch-demo/internal/app/consumer/model"

type Worker interface {
	//注入TaskDataStorage
	TaskDataStorage(storage TaskDataStorage) Worker
	//注入任务执行器Executor
	Executor(executor Executor) Worker

	//处理消息
	ProcessTask(job *model.Job, handler *model.TaskHandler)
}
