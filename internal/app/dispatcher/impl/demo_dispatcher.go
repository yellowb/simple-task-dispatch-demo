package impl

import (
	"github.com/go-co-op/gocron/v2"
	"github.com/yellowb/simple-task-dispatch-demo/internal/app/dispatcher/iface"
	"github.com/yellowb/simple-task-dispatch-demo/internal/app/dispatcher/model"
	"github.com/yellowb/simple-task-dispatch-demo/internal/global"
)

// DemoDispatcher 一个作为Demo的Dispatcher实现
type DemoDispatcher struct {
	// Dispatcher配置
	cfg *global.DispatcherConfig

	// Dispatcher的依赖组建
	deliverier     *iface.Deliverier
	statusStorage  *iface.StatusStorage
	taskDatasource *iface.TaskDatasource

	// 第三方定时调度库
	scheduler *gocron.Scheduler
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

func (d *DemoDispatcher) TaskDatasource(datasource *iface.TaskDatasource) iface.Dispatcher {
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
