package global

// RedisCfg Redis配置信息
type RedisCfg struct {
	Addr     string
	Port     string
	Password string
	Db       int
	Protocol int
}

// DispatcherConfig 工作分派者进程启动配置
type DispatcherConfig struct {
	RedisCfg
}

// WorkerConfig 消费者进程启动配置
type WorkerConfig struct {
	RedisCfg
	// 每一个Worker节点启动多少个go routine用于处理Task
	WorkerPoolSize int
}
