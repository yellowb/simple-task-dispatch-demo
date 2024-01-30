package model

// Job Dispatcher分派给Consumer的一项工作（动态）
type Job struct {
	JobUid string

	// Task相关信息
	TaskName string
	TaskKey  string // Consumer端靠这个来路由到真正的处理逻辑

	// Job相关信息
	DispatchedTime int64                  // Job被分派的时间
	Args           map[string]interface{} // Job运行需要的参数
}
