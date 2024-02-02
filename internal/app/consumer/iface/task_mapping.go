package iface

import "github.com/yellowb/simple-task-dispatch-demo/internal/app/consumer/model"

type TaskMapping interface {
	LoadTaskHandler()
	GetMappingKey(job *model.Job) string
	GetTaskHandler(key string) (*model.TaskHandler, error)
}
