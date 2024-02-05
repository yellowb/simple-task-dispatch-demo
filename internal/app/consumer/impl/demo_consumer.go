package impl

import (
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/yellowb/simple-task-dispatch-demo/internal/app/consumer/iface"
	"github.com/yellowb/simple-task-dispatch-demo/internal/app/consumer/status"
	"github.com/yellowb/simple-task-dispatch-demo/internal/app/consumer/task_executor"
	"github.com/yellowb/simple-task-dispatch-demo/internal/global"
	"github.com/yellowb/simple-task-dispatch-demo/internal/model"
	"sync"
	"sync/atomic"
)

type DemoConsumer struct {
	cfg             *global.ConsumerConfig
	receiver        iface.Receiver
	jobExecutorChan chan *iface.JobExecutor
	status          status.Status
	lock            sync.Mutex
	isStop          atomic.Bool
}

// NewConsumer returns a new instance of DemoConsumer
func NewConsumer() iface.Consumer {
	return &DemoConsumer{
		status: status.New,
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

// 初始化
func (d *DemoConsumer) Init() error {
	d.lock.Lock()
	defer d.lock.Unlock()
	err := status.CheckStatus(d.status, status.New)
	if err != nil {
		return err
	}
	if d.cfg == nil {
		return errors.New("cfg is nil")
	}
	if d.receiver == nil {
		return errors.New("receiver is nil")
	}
	d.jobExecutorChan = make(chan *iface.JobExecutor, d.cfg.JobExecutorChanSize)
	d.status = status.Initialized
	return nil
}

// 启动
func (d *DemoConsumer) Run() error {
	d.lock.Lock()
	defer d.lock.Unlock()
	err := status.CheckStatus(d.status, status.Initialized)
	if err != nil {
		return err
	}
	jobChan := d.receiver.GetJobChan()
	go func(d *DemoConsumer) {
		for !d.isStop.Load() {
			select {
			case job, ok := <-jobChan:
				if !ok {
					logrus.Printf("receiver job chan is closed")
					break
				}
				jobExecutor := d.GetMappingExecutor(job)
				if jobExecutor == nil {
					continue
				}
				d.jobExecutorChan <- jobExecutor
			}
		}
		d.isStop.Store(false)
	}(d)
	d.status = status.Running
	return nil
}

// 关闭
func (d *DemoConsumer) Shutdown() error {
	d.lock.Lock()
	defer d.lock.Unlock()
	err := status.CheckStatus(d.status, status.Running)
	if err != nil {
		return err
	}
	d.status = status.Stopped
	d.isStop.Store(true)
	return nil
}

func (d *DemoConsumer) GetJobExecutor() <-chan *iface.JobExecutor {
	return d.jobExecutorChan
}

func (d *DemoConsumer) GetMappingExecutor(job *model.Job) *iface.JobExecutor {
	mappingKey := job.TaskKey
	executor, ok := task_executor.JobExecutorMap[mappingKey]
	if !ok {
		logrus.Errorf("task_key = %s, executor no found", job.TaskKey)
		return nil
	}
	return &iface.JobExecutor{
		Job:      *job,
		Executor: executor,
	}
}
