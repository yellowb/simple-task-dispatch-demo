package global

// DispatcherConfig 工作分派者进程启动配置
type DispatcherConfig struct {
	// TODO：redis configs
}

// WorkerConfig 消费者进程启动配置
type WorkerConfig struct {
	// 每一个Worker节点启动多少个go routine用于处理Task
	WorkerPoolSize int
}
