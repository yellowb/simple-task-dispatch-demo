package model

import (
	"github.com/go-co-op/gocron/v2"
	"github.com/yellowb/simple-task-dispatch-demo/internal/constants"
	"sync/atomic"
	"time"
)

// Task 代表配置好的任务（静态）
type Task struct {
	Name     string // Task名
	Key      string // Task唯一标识
	TaskType string // Task类型（固定间隔/每日固定时刻）

	// 下面3个字段分别对应1种Task类型，有且只有1个不为空
	DurationDef *DurationDefinition
	DailyDef    *DailyDefinition
}

// DurationDefinition 固定间隔相关配置
type DurationDefinition struct {
	Duration time.Duration
}

// DailyDefinition 每日固定时刻相关配置, 每日x时x分x秒
type DailyDefinition struct {
	Hour   uint
	Minute uint
	Second uint
}

// GenerateJob 根据Task产生一个Job实例
func (t *Task) GenerateJob(jobUid string, dispatchedTime, seq int64) *Job {
	return &Job{
		JobUid:         jobUid,
		TaskName:       t.Name,
		TaskKey:        t.Key,
		ExecutorKey:    t.Key, // 现在简单的等于taskKey, 但是对于云厂同步任务来说, 不能直接取taskKey, 因为多个不同的TaskKey可能对应同一种逻辑Handler(云厂+资源类型 -> 一种Handler)
		DispatchedTime: dispatchedTime,
		Seq:            seq,
		Args: map[string]interface{}{
			// fake
			"arg1": "YB",
			"arg2": 999,
		},
	}
}

// ToGocronJobDefinition 根据Task转换成gocron的JobDefinition
func (t *Task) ToGocronJobDefinition() gocron.JobDefinition {
	if t.TaskType == constants.TaskTypeDuration {
		// 固定间隔任务
		return gocron.DurationJob(t.DurationDef.Duration)
	} else {
		// 每日定时任务
		return gocron.DailyJob(
			1,
			gocron.NewAtTimes(
				gocron.NewAtTime(t.DailyDef.Hour, t.DailyDef.Minute, t.DailyDef.Second),
			),
		)
	}
}

// ToRunningTaskStatus 根据Task转换成RunningTaskStatus对象
func (t *Task) ToRunningTaskStatus() *RunningTaskStatus {
	return &RunningTaskStatus{
		TaskName: t.Name,
		TaskKey:  t.Key,
		Seq:      atomic.Int64{}, // 一个新的RunningTaskStatus的Seq从0开始
	}
}
