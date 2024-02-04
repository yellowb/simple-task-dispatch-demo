package impl

import (
	"context"
	"errors"
	"fmt"
	"github.com/yellowb/simple-task-dispatch-demo/internal/app/consumer/iface"
	"github.com/yellowb/simple-task-dispatch-demo/internal/app/consumer/worker_status"
	"strings"
	"sync"
)

const (
	namePrefix = "DemoWorker-"
)

// DemoWorker 一个Demo的Worker实现
type DemoWorker struct {
	// 父Consumer
	father iface.Consumer

	// 用于停止自身Worker协程
	ctx        context.Context
	cancelFunc context.CancelFunc

	// 其它Worker自身属性
	id     int                        // Worker编号
	name   string                     // Worker名字，当前只是用于打日志用
	status worker_status.WorkerStatus // Worker状态
	lock   sync.Mutex
}

func NewDemoWorker(workerId int) *DemoWorker {
	return &DemoWorker{
		id:     workerId,
		name:   fmt.Sprintf("%s%d", namePrefix, workerId), // e.g. "DemoWorker-1"
		status: worker_status.New,
		lock:   sync.Mutex{},
	}
}

func (d *DemoWorker) Consumer(father iface.Consumer) iface.Worker {
	d.father = father
	return d
}

func (d *DemoWorker) Init() error {
	d.lock.Lock()
	defer d.lock.Unlock()

	// 禁止重复初始化
	err := d.checkStatus(worker_status.New)
	if err != nil {
		return err
	}

	if d.father == nil {
		return errors.New("father consumer is nil")
	}

	// 暂时没别的事情做，正常来说可以做一些DAO初始化之类的

	d.status = worker_status.Initialized
	return nil
}

func (d *DemoWorker) Run() error {
	d.lock.Lock()
	defer d.lock.Unlock()

	// 状态检查
	err := d.checkStatus(worker_status.Initialized, worker_status.Stopped)
	if err != nil {
		return err
	}

	// 每次运行都产生新的ctx，因为ctx没法重复cancel
	d.ctx, d.cancelFunc = context.WithCancel(context.Background())

	// TODO: 起一个协程，包装一个处理函数

	d.status = worker_status.Running
	return nil
}

func (d *DemoWorker) Stop() error {
	//TODO implement me
	panic("implement me")
}

// 检查当前Dispatcher的状态是否在给定的statusList集合中
func (d *DemoWorker) checkStatus(statusList ...worker_status.WorkerStatus) error {
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
	return fmt.Errorf("dispatcher status must be in %s, current status: %s", strings.Join(statusStrList, "/"), d.status.String())
}
