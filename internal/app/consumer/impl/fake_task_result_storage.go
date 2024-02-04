package impl

import (
	"github.com/yellowb/simple-task-dispatch-demo/internal/model"
	"log"
)

// FakeTaskResultStorage 假的TaskResultStorage
type FakeTaskResultStorage struct {
	// 应该要有一些DAO什么的...
}

func (f *FakeTaskResultStorage) SaveJobExecRecord(jobExecRecord *model.JobExecRecord) {
	log.Printf("job exec record saved : %s : %s", jobExecRecord.Job.TaskKey, jobExecRecord.Job.JobUid)
}
