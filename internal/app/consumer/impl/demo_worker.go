package impl

import (
	"context"
	"errors"
	"fmt"
	"github.com/yellowb/simple-task-dispatch-demo/internal/app/consumer/iface"
	"github.com/yellowb/simple-task-dispatch-demo/internal/app/consumer/worker_status"
	"github.com/yellowb/simple-task-dispatch-demo/internal/model"
	"log"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

const (
	namePrefix = "DemoWorker-"
)

// DemoWorker 一个Demo的Worker实现
type DemoWorker struct {
	// 父Consumer
	father iface.Consumer

	// Job执行结果的存储器
	taskResultStorage iface.TaskResultStorage

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

func (d *DemoWorker) TaskResultStorage(storage iface.TaskResultStorage) iface.Worker {
	d.taskResultStorage = storage
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

	// 起一个协程，不断从父Consumer的channel获取Job执行
	ch := d.father.GetJobExecutorPairChannel()
	go func(worker *DemoWorker) {
		for {
			select {
			case v, ok := <-ch:
				{
					if !ok {
						ch = nil
						log.Printf("[%s] father channel closed, worker exited", d.name)
						break
					}
					d.work(v)
				}
			case <-d.ctx.Done():
				log.Printf("[%s] work stopped", d.name)
				break
			}
		}
	}(d)

	d.status = worker_status.Running
	return nil
}

func (d *DemoWorker) Stop() error {
	d.lock.Lock()
	defer d.lock.Unlock()

	// 状态检查
	err := d.checkStatus(worker_status.Running)
	if err != nil {
		return err
	}

	// 停止worker协程
	d.cancelFunc()

	// TODO：实际上可能还需要一些清理工作，这里没有就不写

	d.status = worker_status.Stopped
	return nil
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

// 包装一个方法用于worker协程执行executor逻辑，把通用的前置和后置逻辑抽到这个方法中
func (d *DemoWorker) work(pair *model.JobExecutorPair) {
	job := pair.Job
	executor := pair.Executor

	// 前置工作
	// TODO: 正式的代码应该还需要获取一个分布式锁，Demo忽略了
	var (
		record    *model.JobExecRecord = new(model.JobExecRecord) // Job执行记录
		result    *model.JobExecResult = nil                      // Executor的执行结果
		err       error                = nil
		startTime time.Time            = time.Now() // 开始时间
	)

	// 后置工作, 主要是把执行结果记录写入db
	defer func() {
		// 在Executor逻辑发生panic时要捕获，防止worker挂了，同时也要把panic信息作为result存下来
		if panicErr := recover(); panicErr != nil {
			stack := string(debug.Stack()) // 获取goroutine调用栈信息，帮助后续排查问题
			result = &model.JobExecResult{
				Succeed: false,
				ErrMsg:  fmt.Sprintf("err=%v, stack=%s", panicErr, stack),
			}
			log.Printf("[Panic][%s] worker panic error : %v, stack : %v", d.name, panicErr, stack)
		}

		// 计算耗时等元信息，包装成JobExecRecord对象存储到db
		record.Result = result
		record.Job = job
		record.StartTime = startTime.Unix()
		record.EndTime = time.Now().Unix()
		record.TimeCost = record.EndTime - record.StartTime
		record.Delay = record.StartTime - job.DispatchedTime

		// 内层Defer，防止下面保存JobExecRecord的逻辑也panic了
		defer func() {
			if panicErr := recover(); panicErr != nil {
				// 太惨了，除了打条日志之外无能为力
				log.Printf("[Panic][%s] worker panic error : %v", d.name, panicErr)
			}
		}()

		storeErr := d.taskResultStorage.SaveJobExecRecord(record)
		if storeErr != nil {
			log.Printf("[Error][%s] save job exec record error: %v", d.name, storeErr)
			return
		}
	}()

	// 调用Executor逻辑处理Job
	result, err = executor.Execute(job) // 注意有些傻逼写的Executor可能会panic
	if err != nil {
		// 如果Executor返回的result为nil, 也要手工创建一个用于后续存入数据库
		if result == nil {
			result = &model.JobExecResult{
				Succeed: false,
			}
		}
		// 优先用Executor返回的result中的错误信息，如果没填，则用err替代之
		if result.ErrMsg == "" {
			result.ErrMsg = err.Error()
		}
	}
}
