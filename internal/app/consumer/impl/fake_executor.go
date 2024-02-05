package impl

import (
	"github.com/yellowb/simple-task-dispatch-demo/internal/model"
	"time"
)

/**
假的业务Executor
*/

type FakeExecutorA struct {
}

func (f *FakeExecutorA) ExecutorKey() string {
	return "task-a"
}

func (f *FakeExecutorA) Execute(job *model.Job) (*model.JobExecResult, error) {
	time.Sleep(2 * time.Second)
	return &model.JobExecResult{
		Succeed: true,
		ErrMsg:  "",
		Logs: []string{
			"task-a-log-1",
		},
		Output: map[string]interface{}{
			"job_uid": job.JobUid,
		},
	}, nil
}

type FakeExecutorB struct {
}

func (f *FakeExecutorB) ExecutorKey() string {
	return "task-b"
}

func (f *FakeExecutorB) Execute(job *model.Job) (*model.JobExecResult, error) {
	time.Sleep(3 * time.Second)
	return &model.JobExecResult{
		Succeed: true,
		ErrMsg:  "",
		Logs: []string{
			"task-b-log-1",
		},
		Output: map[string]interface{}{
			"job_uid": job.JobUid,
		},
	}, nil
}
