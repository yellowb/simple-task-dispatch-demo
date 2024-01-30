package iface

import (
	"github.com/yellowb/simple-task-dispatch-demo/internal/app/dispatcher/model"
	"github.com/yellowb/simple-task-dispatch-demo/internal/global"
)

// Dispatcher 工作分派者接口
type Dispatcher interface {
	/* 下面这些是Dispatcher启动前的配置接口 */

	// Config 注入配置到Dispatcher
	Config(cfg *global.DispatcherConfig) Dispatcher
	// Deliverier 注入Job投递器到Dispatcher
	Deliverier(deliverier *Deliverier) Dispatcher
	// StatusStorage 注入内部状态存储器到Dispatcher
	StatusStorage(storage *StatusStorage) Dispatcher
	// TaskDatasource 注入Task数据源到Dispatcher
	TaskProvider(provider *TaskDatasource) Dispatcher

	/* 下面这些是改变Dispatcher中被调度Task的接口 */

	// Add 往Dispatcher中增加Task
	Add(tasks ...*model.Task) error
	// Remove 移除Dispatcher中已有的Task
	Remove(taskKeys ...string) error
	// Refresh 刷新一遍所有要调度的Task
	Refresh() error

	/* 下面这些是影响Dispatcher运行状态的接口 */

	// Run 启动Dispatcher
	Run() error
	// Stop 停止Dispatcher
	Stop() error
}
