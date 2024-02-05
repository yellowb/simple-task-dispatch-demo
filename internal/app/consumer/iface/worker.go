package iface

/**
Worker表示一个工作线程
*/

// Worker 工作线程接口
type Worker interface {
	/* Worker启动前的配置相关接口 */

	// Consumer 注入父亲Consumer，可以通过父亲获取Job、把Job运行结果提交回给父亲等。具体实现比较自由。
	Consumer(father Consumer) Worker
	// TaskResultStorage 注入TaskResultStorage，用于被worker调用存储job执行结果
	TaskResultStorage(storage TaskResultStorage) Worker

	/* 改变Worker运行状态的接口 */

	// Init 初始化Worker
	Init() error
	// Run 启动Worker
	Run() error
	// Stop 停止Worker
	Stop() error
}
