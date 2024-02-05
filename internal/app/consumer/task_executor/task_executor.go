package task_executor

import (
	"github.com/yellowb/simple-task-dispatch-demo/internal/app/consumer/iface"
)

var JobExecutorMap map[string]iface.Executor

func init() {
	JobExecutorMap = make(map[string]iface.Executor)
	loadJobExecutor()

}

func loadJobExecutor() {
	JobExecutorMap["task-5s"] = NewSimpleExecutor()
	JobExecutorMap["task-10s"] = NewSimpleExecutor()
}
