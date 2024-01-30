package model

import "time"

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
