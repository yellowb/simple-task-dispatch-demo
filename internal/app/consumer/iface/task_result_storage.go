package iface

import "github.com/yellowb/simple-task-dispatch-demo/internal/model"

/**
TaskResultStorage 表示任务执行结果存储器接口
*/

type TaskResultStorage interface {
	// SaveJobExecRecord 保存一条Job运行记录
	SaveJobExecRecord(jobExecRecord *model.JobExecRecord)
}
