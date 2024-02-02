package iface

import (
	"context"
	"github.com/yellowb/simple-task-dispatch-demo/internal/app/consumer/model"
)

type Executor interface {
	//任务执行
	Execute(ctx context.Context, taskHandler *model.TaskHandler) error
}
