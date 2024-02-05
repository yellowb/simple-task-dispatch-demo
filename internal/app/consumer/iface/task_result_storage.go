package iface

import "github.com/yellowb/simple-task-dispatch-demo/internal/model"

type TaskResultStorage interface {
	//任务执行记录
	SaveTaskRecord(record *model.JobRecord) error
	//更新任务执行结果
	UpdateTaskStatus(taskStatus *model.TaskStatus) error
}
