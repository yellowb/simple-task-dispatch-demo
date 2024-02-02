package impl

import "github.com/yellowb/simple-task-dispatch-demo/internal/app/consumer/model"

type MongoTaskDataStorage struct {
	//taskRecordDao
	//taskStatusDao
}

func (m *MongoTaskDataStorage) SaveTaskRecord(job *model.Job, isSuccess bool, result interface{}, logs []string) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoTaskDataStorage) UpdateTaskStatus(taskKey string, isSuccess bool) {
	//TODO implement me
	panic("implement me")
}
