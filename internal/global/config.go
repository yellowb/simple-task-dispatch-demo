package global

// DispatcherConfig 工作分派者进程启动配置
type DispatcherConfig struct {
	// TODO：redis configs
}

// ConsumerConfig 消费者进程启动配置
type ConsumerConfig struct {
	// 每一个Consumer节点启动多少个go routine用于处理Task
	ConsumerPoolSize int
}
