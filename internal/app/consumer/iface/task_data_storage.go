package iface

import "github.com/yellowb/simple-task-dispatch-demo/internal/app/consumer/model"

type TaskDataStorage interface {
	//任务执行记录
	SaveTaskRecord(job *model.Job, isSuccess bool, result interface{}, logs []string)
	//更新任务执行结果
	UpdateTaskStatus(taskKey string, isSuccess bool)
}
