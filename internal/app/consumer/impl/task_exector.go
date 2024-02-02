package impl

import (
	"context"
	"github.com/yellowb/simple-task-dispatch-demo/internal/app/consumer/model"
)

type TaskExecutor struct{}

func (TaskExecutor) Execute(ctx context.Context, taskHandler *model.TaskHandler) error {
	//TODO implement me
	panic("implement me")
}
