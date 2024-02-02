package impl

import (
	"fmt"
	"github.com/yellowb/simple-task-dispatch-demo/internal/app/consumer/iface"
	"github.com/yellowb/simple-task-dispatch-demo/internal/global"
	"sync"
)

type DemoConsumer struct {
	cfg         *global.ConsumerConfig
	worker      iface.Worker
	taskMapping iface.TaskMapping
	provider    iface.Provider
	waitGroup   sync.WaitGroup
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

func (d *DemoConsumer) ConsumeTask() {
	for i := 0; i < d.cfg.WorkerPoolSize; i++ {
		d.waitGroup.Add(1)
		go func() {
			defer d.waitGroup.Done()
			d.worker.ProcessTask()
		}()
	}

	for task := range d.provider.GetTaskChan() {
		// Process the task based on the task mapping
		fmt.Println("Processing task:", task)
	}
}
