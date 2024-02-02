package model

// DispatcherStatus Dispatcher状态
type DispatcherStatus struct {
	// Dispatcher当前调度的Task总数
	RunningTaskCount int
	// 所有调度中的Task的运行时状态
	RunningTaskStatuses map[string]*RunningTaskStatus
}
