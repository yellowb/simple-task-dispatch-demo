package iface

type Worker interface {
	//注入TaskDataStorage
	TaskDataStorage(storage TaskDataStorage) Worker
	//注入任务执行器Executor
	Executor(executor Executor) Worker

	//处理消息
	ProcessTask()
}
