package impl

import (
	"github.com/yellowb/simple-task-dispatch-demo/internal/model"
	"log"
)

// FakeTaskResultStorage 假的TaskResultStorage
type FakeTaskResultStorage struct {
	// 应该要有一些DAO什么的...
}

func (f *FakeTaskResultStorage) SaveJobExecRecord(who string, jobExecRecord *model.JobExecRecord) error {
	log.Printf("[%s] job exec record saved : %s : %s", who, jobExecRecord.Job.TaskKey, jobExecRecord.Job.JobUid)
	return nil
}
