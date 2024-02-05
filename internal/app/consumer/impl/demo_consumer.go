package impl

import (
	"context"
	"errors"
	"fmt"
	"github.com/yellowb/simple-task-dispatch-demo/internal/app/consumer/consumer_status"
	"github.com/yellowb/simple-task-dispatch-demo/internal/app/consumer/iface"
	"github.com/yellowb/simple-task-dispatch-demo/internal/global"
	"log"
	"strings"
	"sync"
	"time"
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

	// Job执行结果的存储器
	taskResultStorage iface.TaskResultStorage

	// 用于退出Consumer协程
	ctx        context.Context
	cancelFunc context.CancelFunc

	// 其它内部属性
	workerCount   int                            // Worker数量
	workerPool    []iface.Worker                 // Worker池子，好像暂时没用上
	status        consumer_status.ConsumerStatus // Consumer状态
	lock          sync.Mutex
	executorMap   map[string]iface.Executor   // key -> Executor 的映射
	jobExecutorCh chan *iface.JobExecutorPair // 用于往Worker投递消息的Channel
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

func (d *DemoConsumer) TaskResultStorage(storage iface.TaskResultStorage) iface.Consumer {
	d.taskResultStorage = storage
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
	// TODO: 其实还要检查Receiver的状态是否已经在运行，这里没写

	if d.taskResultStorage == nil {
		return errors.New("task result storage is nil")
	}

	// 初始化内部字段
	d.ctx, d.cancelFunc = context.WithCancel(context.Background())
	if d.cfg.WorkerPoolSize < 0 { // TODO: 应该还要加一个右边界判断，这里没写
		d.workerCount = defaultWorkerCount
	} else {
		d.workerCount = d.cfg.WorkerPoolSize
	}
	d.initExecutorMap()
	if d.cfg.JobBufferSize < 0 { // TODO: 应该还要加一个右边界判断，这里没写
		d.jobExecutorCh = make(chan *iface.JobExecutorPair, defaultChanSize)
	} else {
		d.jobExecutorCh = make(chan *iface.JobExecutorPair, d.cfg.JobBufferSize)
	}

	// 初始化Worker池子
	d.workerPool = make([]iface.Worker, d.workerCount)
	for i := 0; i < d.workerCount; i++ {
		worker := NewDemoWorker(i).Consumer(d).TaskResultStorage(d.taskResultStorage)
		err = worker.Init()
		if err != nil {
			return fmt.Errorf("init worker error : %v", err)
		}
		d.workerPool[i] = worker
	}

	d.status = consumer_status.Initialized
	return nil
}

func (d *DemoConsumer) Run() error {
	d.lock.Lock()
	defer d.lock.Unlock()

	// 状态检查
	err := d.checkStatus(consumer_status.Initialized)
	if err != nil {
		return err
	}

	// 起一个协程不断从Receiver获取Job，投递给Worker处理
	ticker := time.NewTicker(time.Minute)
	receiverCh := d.receiver.GetJobChannel()
	go func(consumer *DemoConsumer) {
		log.Printf("[Consumer] consumer started")
		for {
			select {
			case v, ok := <-receiverCh:
				{
					if !ok {
						// receiver的channel被关闭了，一般不会先关闭receiver。正常应该是先关闭consumer
						// 所以正确使用时应该不会走到这里。
						// 把receiverCh设置为nil防止外层select再跑进这个case导致空循环
						receiverCh = nil
						log.Printf("[Consumer] receiver channel closed")
					}
					// 根据job查找对应的Executor
					executor := consumer.getMatchedExecutor(v.ExecutorKey)
					if executor == nil {
						// 找不到匹配的Executor！正常不应该出现这种情况，打log
						log.Printf("[Warn][Consumer] cannot match executor with key : %s", v.ExecutorKey)
					}
					// 把Job和Executor打包投递给Worker
					pair := &iface.JobExecutorPair{
						Job:      v,
						Executor: executor,
					}
					consumer.jobExecutorCh <- pair
				}
			case <-consumer.ctx.Done():
				{
					// consumer 被 shutdown
					log.Printf("[Consumer] consumer shutdown")
					break
				}
			case <-ticker.C:
				{
					// 打印一句日志表示自己还活着
					log.Printf("[Consumer] i am alive")
				}
			}
		}
	}(d)

	// 启动所有Worker协程
	for _, worker := range d.workerPool {
		err = worker.Run()
		if err != nil {
			log.Printf("[Error][Consumer] run worker error : %v", err)
			// 其实这里直接return err也是有问题的，因为上面那个consumer协程还在跑，
			// 但是consumer状态还是Initialized。比较完备的做法应该是这里要停止掉consumer协程，再返回。

			// 不过一般不会出错，可以先这么对付一下=_=
			return err
		}
	}

	d.status = consumer_status.Running
	return nil
}

func (d *DemoConsumer) Shutdown() error {
	d.lock.Lock()
	defer d.lock.Unlock()

	// 状态检查
	err := d.checkStatus(consumer_status.Running)
	if err != nil {
		return err
	}

	// 停止Consumer协程
	d.cancelFunc()
	// 关闭Receiver
	_ = d.receiver.Shutdown()
	// 关闭与Worker相连的channel，通知所有Worker退出
	close(d.jobExecutorCh)

	d.status = consumer_status.Shutdown

	return nil
}

func (d *DemoConsumer) GetJobExecutorPairChannel() <-chan *iface.JobExecutorPair {
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
func (d *DemoConsumer) initExecutorMap() {
	// TODO: 正式的代码可以把这个映射表做成全局单例，初始化的地方也不要在Consumer中，
	// TODO: Consumer直接引用这个全局单例即可

	// 这里搞了两个假的executor
	executorMap := make(map[string]iface.Executor)

	e1 := new(FakeExecutorA)
	executorMap[e1.ExecutorKey()] = e1

	e2 := new(FakeExecutorB)
	executorMap[e2.ExecutorKey()] = e2

	d.executorMap = executorMap
}

func (d *DemoConsumer) getMatchedExecutor(key string) iface.Executor {
	return d.executorMap[key]
}
