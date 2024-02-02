package impl

import (
	"github.com/sirupsen/logrus"
	"github.com/yellowb/simple-task-dispatch-demo/internal/app/consumer/iface"
	"github.com/yellowb/simple-task-dispatch-demo/internal/app/consumer/model"
	"github.com/yellowb/simple-task-dispatch-demo/internal/global"
	"sync"
)

type DemoConsumer struct {
	cfg         *global.ConsumerConfig
	worker      iface.Worker
	taskMapping iface.TaskMapping
	provider    iface.Provider
	waitGroup   sync.WaitGroup
	workerPool  chan struct{}
}

// NewConsumer returns a new instance of DemoConsumer
func NewConsumer() iface.Consumer {
	return &DemoConsumer{
		waitGroup: sync.WaitGroup{},
	}
}

func (d *DemoConsumer) Config(cfg *global.ConsumerConfig) iface.Consumer {
	d.cfg = cfg
	return d
}

func (d *DemoConsumer) Worker(worker iface.Worker) iface.Consumer {
	d.worker = worker
	return d
}

func (d *DemoConsumer) TaskMapping(taskMapping iface.TaskMapping) iface.Consumer {
	d.taskMapping = taskMapping
	return d
}

func (d *DemoConsumer) Provider(provider iface.Provider) iface.Consumer {
	d.provider = provider
	return d
}

// 初始化
func (d *DemoConsumer) Init() error {
	//TODO implement me
	panic("implement me")
}

// 启动
func (d *DemoConsumer) Run() error {
	//TODO implement me
	panic("implement me")
}

// 关闭
func (d *DemoConsumer) Shutdown() error {
	//TODO implement me
	panic("implement me")
}

func (d *DemoConsumer) ConsumeTask() {
	for job := range d.provider.GetTaskChan() {
		taskHandler, err := d.taskMapping.GetTaskHandler(d.taskMapping.GetMappingKey(job))
		if err != nil {
			logrus.Errorf("")
			continue
		}
		// 获取一个空结构体，表示占用了一个worker的资源,如果workerPool满了将阻塞在这里
		d.workerPool <- struct{}{}
		go func(job *model.Job, handler *model.TaskHandler) {
			defer func() {
				// 释放一个worker的资源
				<-d.workerPool
				d.waitGroup.Done()
			}()
			d.waitGroup.Add(1)
			//调用worker处理任务
			d.worker.ProcessTask(job, taskHandler)
		}(job, taskHandler)
	}
}
