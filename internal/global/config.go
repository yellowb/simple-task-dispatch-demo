package global

// RedisCfg Redis配置信息
type RedisCfg struct {
	Addr     string
	Port     string
	Password string
	Db       int
}

// DispatcherConfig 工作分派者进程启动配置
type DispatcherConfig struct {
	RedisCfg
	// 消息队列名字
	QueueName string
}

// ConsumerConfig 消费者进程启动配置
type ConsumerConfig struct {
	RedisCfg
	// 消息队列名字
	QueueName string
	// 每一个Worker节点启动多少个go routine用于处理Task
	WorkerPoolSize int
}
