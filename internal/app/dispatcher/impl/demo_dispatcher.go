package impl

import (
	"github.com/yellowb/simple-task-dispatch-demo/internal/app/dispatcher/iface"
	"github.com/yellowb/simple-task-dispatch-demo/internal/app/dispatcher/model"
	"github.com/yellowb/simple-task-dispatch-demo/internal/global"
)

// DemoDispatcher 一个作为Demo的Dispatcher实现
type DemoDispatcher struct {
	// TODO：一些字段
}

func (d *DemoDispatcher) Config(cfg *global.DispatcherConfig) iface.Dispatcher {
	//TODO implement me
	panic("implement me")
}

func (d *DemoDispatcher) Deliverier(deliverier *iface.Deliverier) iface.Dispatcher {
	//TODO implement me
	panic("implement me")
}

func (d *DemoDispatcher) StatusStorage(storage *iface.StatusStorage) iface.Dispatcher {
	//TODO implement me
	panic("implement me")
}

func (d *DemoDispatcher) TaskProvider(provider *iface.TaskDatasource) iface.Dispatcher {
	//TODO implement me
	panic("implement me")
}

func (d *DemoDispatcher) Add(tasks ...*model.Task) error {
	//TODO implement me
	panic("implement me")
}

func (d *DemoDispatcher) Remove(taskKeys ...string) error {
	//TODO implement me
	panic("implement me")
}

func (d *DemoDispatcher) Refresh() error {
	//TODO implement me
	panic("implement me")
}

func (d *DemoDispatcher) Run() error {
	//TODO implement me
	panic("implement me")
}

func (d *DemoDispatcher) Stop() error {
	//TODO implement me
	panic("implement me")
}
