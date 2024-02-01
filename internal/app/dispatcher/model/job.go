package model

// Job Dispatcher分派给Consumer的一项工作（动态）
type Job struct {
	JobUid string `json:"job_uid"`

	// Task相关信息
	TaskName   string `json:"task_name"`
	TaskKey    string `json:"task_key"`
	HandlerKey string `json:"handler_key"` // Consumer端靠这个来路由到真正的处理逻辑

	// Job相关信息
	DispatchedTime int64                  `json:"dispatched_time"` // Job被分派的时间
	Seq            int64                  `json:"seq"`             // 这个Task从Dispatcher进程启动开始累计被调度多少次（从1开始）
	Args           map[string]interface{} `json:"args"`            // Job运行需要的参数
}
