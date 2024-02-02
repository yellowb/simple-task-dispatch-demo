package impl

import (
	"context"
	"github.com/yellowb/simple-task-dispatch-demo/internal/app/consumer/model"
)

type TaskExecutor struct{}

func (TaskExecutor) Execute(ctx context.Context, taskHandler *model.TaskHandler) (bool, error) {
	taskHandler.ExecFunc(ctx)
	if ctx.Err() != nil {
		return false, ctx.Err()
	}
	return true, nil
}
