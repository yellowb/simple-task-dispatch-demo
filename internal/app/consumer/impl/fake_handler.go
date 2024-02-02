package impl

import (
	"github.com/yellowb/simple-task-dispatch-demo/internal/model"
	"time"
)

/**
假的业务Handler
*/

type FakeHandlerA struct {
}

func (f *FakeHandlerA) HandlerKey() string {
	return "task-a"
}

func (f *FakeHandlerA) Handle(job *model.Job) (*model.JobExecResult, error) {
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

type FakeHandlerB struct {
}

func (f FakeHandlerB) HandlerKey() string {
	return "task-b"
}

func (f FakeHandlerB) Handle(job *model.Job) (*model.JobExecResult, error) {
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
