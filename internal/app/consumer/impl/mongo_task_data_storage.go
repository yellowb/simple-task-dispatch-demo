package impl

import (
	"fmt"
	"github.com/yellowb/simple-task-dispatch-demo/internal/model"
)

type MongoTaskDataStorage struct {
	//taskRecordDao
	//taskStatusDao
}

func NewMongoTaskDataStorage() *MongoTaskDataStorage {
	return &MongoTaskDataStorage{}
}

func (m *MongoTaskDataStorage) SaveTaskRecord(jobRecord *model.JobRecord) error {
	fmt.Printf("insert job record success, record = %v\n", jobRecord)
	return nil
}

func (m *MongoTaskDataStorage) UpdateTaskStatus(taskStatus *model.TaskStatus) error {
	fmt.Printf("upsert task status success, record = %v\n", taskStatus)
	return nil
}
