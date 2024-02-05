package iface

import (
	"github.com/yellowb/simple-task-dispatch-demo/internal/global"
	"github.com/yellowb/simple-task-dispatch-demo/internal/model"
)

/**
Consumer表示Job的消费者
*/

// Consumer Job消费者接口
type Consumer interface {
	/* 下面这些是Consumer启动前的配置接口 */

	// Config 注入配置到Consumer
	Config(cfg *global.ConsumerConfig) Consumer
	// Receiver 注入Job的接收器到Consumer
	Receiver(receiver Receiver) Consumer
	// TaskResultStorage 注入TaskResultStorage，用于被worker调用存储job执行结果
	TaskResultStorage(storage TaskResultStorage) Consumer

	/* 下面这些是影响Consumer运行状态的接口 */

	// Init 初始化Consumer
	Init() error
	// Run 启动Consumer
	Run() error
	// Shutdown 终止Consumer
	Shutdown() error

	/* 下面这些是Consumer暴露给Worker使用的接口 */

	// GetJobExecutorPairChannel 返回一个channel，后续Worker可以从这个channel中拿到Consumer派发的Job和对应的Executor组合
	GetJobExecutorPairChannel() <-chan *model.JobExecutorPair
}
