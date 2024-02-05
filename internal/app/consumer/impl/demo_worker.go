package impl

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/yellowb/simple-task-dispatch-demo/internal/app/consumer/iface"
	"github.com/yellowb/simple-task-dispatch-demo/internal/app/consumer/status"
	"github.com/yellowb/simple-task-dispatch-demo/internal/model"
	"runtime"
	"sync"
	"time"
)

type DemoWorker struct {
	taskResultStorage iface.TaskResultStorage
	executor          iface.Executor
	consumer          iface.Consumer
	ctx               context.Context
	cancelFunc        context.CancelFunc
	id                int
	status            status.Status
	lock              sync.Mutex
}

func NewWorker(id int) iface.Worker {
	return &DemoWorker{
		id:     id,
		status: status.New,
		lock:   sync.Mutex{},
	}
}

func (d *DemoWorker) TaskResultStorage(taskResultStorage iface.TaskResultStorage) iface.Worker {
	d.taskResultStorage = taskResultStorage
	return d
}

func (d *DemoWorker) Consumer(consumer iface.Consumer) iface.Worker {
	d.consumer = consumer
	return d
}

func (d *DemoWorker) Init() error {
	d.lock.Lock()
	defer d.lock.Unlock()
	err := status.CheckStatus(d.status, status.New)
	if err != nil {
		return err
	}

	if d.consumer == nil {
		return errors.New("consumer is nil")
	}
	if d.taskResultStorage == nil {
		return errors.New("taskResultStorage is nil")
	}

	d.status = status.Initialized
	return nil
}

func (d *DemoWorker) Run() error {
	d.lock.Lock()
	defer d.lock.Unlock()
	err := status.CheckStatus(d.status, status.Initialized)
	if err != nil {
		return err
	}
	d.ctx, d.cancelFunc = context.WithCancel(context.Background())
	jobCh := d.consumer.GetJobExecutor()
	//开启协程
	go func(d *DemoWorker) {
		for {
			select {
			case job, ok := <-jobCh:
				if !ok {
					//
					logrus.Printf("jobExecutor chan is closed")
					break
				}
				d.processJob(job)
			case <-d.ctx.Done():
				logrus.Printf("worker = %d is stopped", d.id)
				break
			}

		}
	}(d)
	d.status = status.Running
	return nil
}

func (d *DemoWorker) Shutdown() error {
	d.lock.Lock()
	defer d.lock.Unlock()
	err := status.CheckStatus(d.status, status.Running)
	if err != nil {
		return err
	}
	// 停止worker
	d.cancelFunc()
	//将状态置为stopped
	d.status = status.Stopped
	return nil
}

func (d *DemoWorker) processJob(jobExecutor *iface.JobExecutor) {
	defer func() {
		if r := recover(); r != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			logrus.Errorf("running job = %v panic: %v\n%s", jobExecutor.Job, r, buf)
		}
	}()
	//todo 获取该任务的分布式锁

	jobRecord := &model.JobRecord{
		StartTime:  time.Now(),
		StartDelay: time.Now().Unix() - jobExecutor.DispatchedTime,
	}
	jobResult, err := jobExecutor.Executor.Execute(jobExecutor.Args)
	//计算耗时
	jobRecord.Duration = time.Now().Unix() - jobRecord.StartTime.Unix()
	if err != nil {
		jobResult.IsSuccess = false
		logrus.Errorf("[DemoWorker]job execute error, taskKey = %s, duration = %v", jobExecutor.TaskKey, jobRecord.Duration)
	} else {
		jobResult.IsSuccess = true
		logrus.Printf("[DemoWorker]job execute success, taskKey = %s, duration = %v", jobExecutor.TaskKey, jobRecord.Duration)
	}
	jobRecord.JobResult = *jobResult
	//生成任务状态
	taskStatus := &model.TaskStatus{
		TaskKey:   jobExecutor.TaskKey,
		Enabled:   true,
		IsSuccess: jobRecord.IsSuccess,
	}
	//任务执行记录保存
	err = d.taskResultStorage.SaveTaskRecord(jobRecord)
	if err != nil {
		logrus.Errorf("insert job record error, err = %s", err.Error())
	}
	//任务状态更新
	err = d.taskResultStorage.UpdateTaskStatus(taskStatus)
	if err != nil {
		logrus.Errorf("update task status error, err = %s", err.Error())
	}
}
