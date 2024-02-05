package impl

import (
	"context"
	"errors"
	"fmt"
	"github.com/yellowb/simple-task-dispatch-demo/internal/app/consumer/consumer_status"
	"github.com/yellowb/simple-task-dispatch-demo/internal/app/consumer/iface"
	"github.com/yellowb/simple-task-dispatch-demo/internal/global"
	"github.com/yellowb/simple-task-dispatch-demo/internal/model"
	"strings"
	"sync"
)

const (
	defaultWorkerCount = 10
	defaultChanSize    = 10 * 1024
)

// DemoConsumer 一个Demo的Consumer实现
type DemoConsumer struct {
	// 配置信息
	cfg *global.ConsumerConfig

	// Receiver
	receiver iface.Receiver

	// 用于退出Consumer协程
	ctx        context.Context
	cancelFunc context.CancelFunc

	// 其它内部属性
	workerCount   int                            // Worker数量
	status        consumer_status.ConsumerStatus // Consumer状态
	lock          sync.Mutex
	executorMap   map[string]iface.Executor   // key -> Executor 的映射
	jobExecutorCh chan *model.JobExecutorPair // 用于往Worker投递消息的Channel
}

func NewDemoConsumer() *DemoConsumer {
	return &DemoConsumer{
		status: consumer_status.New,
		lock:   sync.Mutex{},
	}
}

func (d *DemoConsumer) Config(cfg *global.ConsumerConfig) iface.Consumer {
	d.cfg = cfg
	return d
}

func (d *DemoConsumer) Receiver(receiver iface.Receiver) iface.Consumer {
	d.receiver = receiver
	return d
}

func (d *DemoConsumer) Init() error {
	d.lock.Lock()
	defer d.lock.Unlock()

	// 禁止重复初始化
	err := d.checkStatus(consumer_status.New)
	if err != nil {
		return err
	}

	if d.cfg == nil {
		return errors.New("config is nil")
	}

	if d.receiver == nil {
		return errors.New("receiver is nil")
	}

	// 初始化内部字段
	d.ctx, d.cancelFunc = context.WithCancel(context.Background())
	if d.cfg.WorkerPoolSize < 0 { // TODO: 应该还要加一个右边界判断，这里没写
		d.workerCount = defaultWorkerCount
	} else {
		d.workerCount = d.cfg.WorkerPoolSize
	}
	d.executorMap = d.initExecutorMap()
	if d.cfg.JobBufferSize < 0 { // TODO: 应该还要加一个右边界判断，这里没写
		d.jobExecutorCh = make(chan *model.JobExecutorPair, defaultChanSize)
	} else {
		d.jobExecutorCh = make(chan *model.JobExecutorPair, d.cfg.JobBufferSize)
	}

	d.status = consumer_status.Initialized
	return nil
}

func (d *DemoConsumer) Run() error {
	//TODO implement me
	panic("implement me")
}

func (d *DemoConsumer) Shutdown() error {
	//TODO implement me
	panic("implement me")
}

func (d *DemoConsumer) GetJobExecutorPairChannel() <-chan *model.JobExecutorPair {
	return d.jobExecutorCh
}

// 检查当前Consumer的状态是否在给定的statusList集合中
func (d *DemoConsumer) checkStatus(statusList ...consumer_status.ConsumerStatus) error {
	for _, status := range statusList {
		if d.status == status {
			return nil
		}
	}

	// 构造错误信息
	statusStrList := make([]string, 0, len(statusList))
	for _, status := range statusList {
		statusStrList = append(statusStrList, status.String())
	}
	return fmt.Errorf("consumer status must be in %s, current status: %s", strings.Join(statusStrList, "/"), d.status.String())
}

// 初始化 key -> Executor 的映射
func (d *DemoConsumer) initExecutorMap() map[string]iface.Executor {
	// TODO: 正式的代码可以把这个映射表做成全局单例，初始化的地方也不要在Consumer中，
	// TODO: Consumer直接引用这个全局单例即可

	// 这里搞了两个假的executor
	executorMap := make(map[string]iface.Executor)

	e1 := new(FakeExecutorA)
	executorMap[e1.ExecutorKey()] = e1

	e2 := new(FakeExecutorA)
	executorMap[e2.ExecutorKey()] = e2

	return executorMap
}
