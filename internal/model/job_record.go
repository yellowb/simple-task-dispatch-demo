package model

import "time"

type JobRecord struct {
	Job
	JobResult
	StartTime  time.Time
	Duration   int64
	StartDelay int64
	Ctime      time.Time
}

type TaskStatus struct {
	TaskKey   string
	Enabled   bool
	IsSuccess bool
	Ctime     time.Time
	Mtime     time.Time
}
