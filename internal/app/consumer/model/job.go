package model

// Job 从队列获取到的任务
type Job struct {
	JobUid string

	// Task相关信息
	TaskName string
	TaskKey  string // Consumer端靠这个来路由到真正的处理逻辑

	// Job相关信息
	DispatchedTime int64                  // Job被分派的时间
	Seq            int64                  // 这个Task从Dispatcher进程启动开始累计被调度多少次（从1开始）
	Args           map[string]interface{} // Job运行需要的参数
}
