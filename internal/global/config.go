package global

type Config struct {
	// 每一个Consumer节点启动多少个go routine用于处理Task
	ConsumerPoolSize int
}
