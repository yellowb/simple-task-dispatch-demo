package model

// RunningTaskStatus 一个被Dispatcher调度的Task的运行时状态
type RunningTaskStatus struct {
	TaskName string
	TaskKey  string
	Seq      int64 // 这个Task从Dispatcher进程启动开始累计被调度多少次（从1开始）
}
