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
	// 启动多少个Worker的go routine用于处理Task
	WorkerPoolSize int
	// 接收到消息放入缓冲队列的大小
	JobExecutorChanSize int
}

// ReceiverConfig 消费者进程启动配置
type ReceiverConfig struct {
	RedisCfg
	// 消息队列名字
	QueueName string
	// 接收到消息放入缓冲队列的大小
	JobChanSize int
}
