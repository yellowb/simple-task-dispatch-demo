package task_executor

import (
	"fmt"
	"github.com/yellowb/simple-task-dispatch-demo/internal/model"
)

type SimpleExecutor struct {
}

func NewSimpleExecutor() *SimpleExecutor {
	return &SimpleExecutor{}
}

func (s *SimpleExecutor) Execute(args map[string]interface{}) (*model.JobResult, error) {
	//定时任务具体业务逻辑执行
	fmt.Println("simple task execute")
	logs := []string{"simple task execute"}
	return &model.JobResult{
		IsSuccess: true,
		Logs:      logs,
	}, nil
}
