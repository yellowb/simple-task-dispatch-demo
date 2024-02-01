package impl

import (
	"errors"
	"fmt"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/yellowb/simple-task-dispatch-demo/internal/app/dispatcher/error_types"
	"github.com/yellowb/simple-task-dispatch-demo/internal/app/dispatcher/iface"
	"github.com/yellowb/simple-task-dispatch-demo/internal/app/dispatcher/model"
	"github.com/yellowb/simple-task-dispatch-demo/internal/constants"
	"github.com/yellowb/simple-task-dispatch-demo/internal/global"
	"log"
	"strings"
	"sync"
	"time"
)

// DemoDispatcher 一个作为Demo的Dispatcher实现
type DemoDispatcher struct {
	// Dispatcher配置
	cfg *global.DispatcherConfig

	// Dispatcher的依赖组件
	deliverier     iface.Deliverier     // Job投递者
	statusStorage  iface.StatusStorage  // Dispatcher状态存储器
	taskDatasource iface.TaskDatasource // Task数据源

	// 第三方调度器
	scheduler gocron.Scheduler

	// Dispatcher其它内部属性
	status constants.DispatcherStatus // Dispatcher状态
	lock   sync.Mutex
}

func NewDemoDispatcher() *DemoDispatcher {
	return &DemoDispatcher{
		status: constants.New,
		lock:   sync.Mutex{},
	}
}

func (d *DemoDispatcher) Config(cfg *global.DispatcherConfig) iface.Dispatcher {
	d.cfg = cfg
	return d
}

func (d *DemoDispatcher) Deliverier(deliverier iface.Deliverier) iface.Dispatcher {
	d.deliverier = deliverier
	return d
}

func (d *DemoDispatcher) StatusStorage(storage iface.StatusStorage) iface.Dispatcher {
	d.statusStorage = storage
	return d
}

func (d *DemoDispatcher) TaskDatasource(datasource iface.TaskDatasource) iface.Dispatcher {
	d.taskDatasource = datasource
	return d
}

// Init 初始化Dispatcher
func (d *DemoDispatcher) Init() error {
	d.lock.Lock()
	defer d.lock.Unlock()

	// 禁止重复初始化
	err := d.checkStatus(constants.New)
	if err != nil {
		return err
	}

	// 简单检查一下各个依赖是否为空即可
	if d.cfg == nil {
		return errors.New("config is nil")
	}
	if d.deliverier == nil {
		return errors.New("deliverier is nil")
	}
	if d.statusStorage == nil {
		return errors.New("status storage is nil")
	}
	if d.taskDatasource == nil {
		return errors.New("task datasource is nil")
	}

	// 创建一个gocron scheduler
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		return err
	}

	d.scheduler = scheduler

	d.status = constants.Initialized
	return nil
}

// Add 往Dispatcher中增加一个Task
func (d *DemoDispatcher) Add(task *model.Task) error {
	d.lock.Lock()
	defer d.lock.Unlock()

	err := d.checkStatus(constants.Initialized, constants.Running)
	if err != nil {
		return err
	}

	// 往scheduler中添加一个调度任务
	_, err = d.addToScheduler(task)

	return err
}

// Remove 根据taskKey把对应Task移除
func (d *DemoDispatcher) Remove(taskKey string) error {
	d.lock.Lock()
	defer d.lock.Unlock()

	err := d.checkStatus(constants.Initialized, constants.Running)
	if err != nil {
		return err
	}

	return d.removeFromScheduler(taskKey)
}

// Reload 重新加载Dispatcher中的所有Task
func (d *DemoDispatcher) Reload() error {
	d.lock.Lock()
	defer d.lock.Unlock()

	err := d.checkStatus(constants.Initialized, constants.Running)
	if err != nil {
		return err
	}

	// 先从数据源获取所有最新的Task
	tasks, err := d.loadLatestTasks()
	if err != nil {
		return fmt.Errorf("load tasks from datasource error : %v", err)
	}

	// 把scheduler中所有调度任务删除、清空状态存储器
	err = d.clearScheduler()
	if err != nil {
		return fmt.Errorf("clear jobs from scheduler error : %v", err)
	}

	// 把从数据源获取的Task全量加回去scheduler
	for _, task := range tasks {
		_, err = d.addToScheduler(task)
		// 这里其实有点问题，如果前面N个Task成功，第N+1个Task出错，则前面N个也不会回滚，也就是不保证批量增加Job的原子性，但好像也没有别的好办法
		if err != nil {
			return fmt.Errorf("create new job in scheduler error : %v", err)
		}
	}

	return nil
}

// Run 让Dispatcher开始调度任务
func (d *DemoDispatcher) Run() error {
	d.lock.Lock()
	defer d.lock.Unlock()

	err := d.checkStatus(constants.Initialized)
	if err != nil {
		return err
	}

	d.scheduler.Start()
	d.status = constants.Running
	return nil
}

// Shutdown 关闭Dispatcher
func (d *DemoDispatcher) Shutdown() error {
	d.lock.Lock()
	defer d.lock.Unlock()

	err := d.checkStatus(constants.Running)
	if err != nil {
		return err
	}

	d.status = constants.Shutdown
	return d.scheduler.Shutdown()
}

// 检查当前Dispatcher的状态是否在给定的statusList集合中
func (d *DemoDispatcher) checkStatus(statusList ...constants.DispatcherStatus) error {
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

// 从Task数据源获取所有最新的Task
func (d *DemoDispatcher) loadLatestTasks() ([]*model.Task, error) {
	return d.taskDatasource.GetAllTasks()
}

// 把一个Task作为调度任务添加进scheduler
func (d *DemoDispatcher) addToScheduler(task *model.Task) (gocron.Job, error) {
	// 不能重复添加TaskKey相同的Task
	ok, err := d.statusStorage.ExistRunningTaskStatus(task.Key)
	if err != nil {
		return nil, fmt.Errorf("check exist running task error : %v", err)
	}
	if ok {
		return nil, error_types.ErrTaskAlreadyExist
	}

	gocronJob, err := d.scheduler.NewJob(
		// 调度类型
		task.ToGocronJobDefinition(),
		// 都是同一个处理函数，只是闭包中的参数不一样
		gocron.NewTask(
			d.taskFunc(task),
		),
		// 把Task的key作为tag附加在这个调度任务上，方便后续根据TaskKey删除调度任务
		gocron.WithTags(task.Key),
	)
	if err != nil {
		return nil, err
	}

	// 添加这个Task进状态存储器进行跟踪
	err = d.statusStorage.PutRunningTaskStatus(task.Key, task.ToRunningTaskStatus())
	if err != nil {
		return nil, err
	}

	return gocronJob, nil
}

// 根据TaskKey从scheduler删除一个调度任务
func (d *DemoDispatcher) removeFromScheduler(taskKey string) error {
	// 由于addToScheduler方法中已经把TaskKey作为调度任务的tag，所以可以简单通过RemoveByTags函数删除
	d.scheduler.RemoveByTags(taskKey)
	// 从状态存储器删除这个Task的跟踪状态
	return d.statusStorage.DeleteRunningTaskStatus(taskKey)
}

// 清空scheduler中所有调度任务、清空状态存储器
func (d *DemoDispatcher) clearScheduler() error {
	// 删除scheduler中所有调度任务
	gocronJobs := d.scheduler.Jobs()
	// 迭代每一个调度任务的UUID删除
	for _, gocronJob := range gocronJobs {
		err := d.scheduler.RemoveJob(gocronJob.ID())
		if err != nil {
			return err
		}
	}

	// 清空状态存储器
	err := d.statusStorage.Clear()
	return err
}

// 更新taskKey对应的Task的运行时状态。本方法会返回最新的RunningTaskStatus对象。
func (d *DemoDispatcher) updateRunningTaskStatus(taskKey string) (*model.RunningTaskStatus, error) {
	taskStatus, err := d.statusStorage.GetRunningTaskStatus(taskKey)
	if err != nil {
		return nil, err
	}
	if taskStatus == nil {
		return nil, fmt.Errorf("task [%s] is not existed", taskKey)
	}

	// 下面是更新逻辑，当前只需要自增Seq。有更复杂的逻辑可以加到下面。
	taskStatus.Seq++
	return taskStatus, nil
}

// 返回一个用于被scheduler调度的闭包函数，由于Dispatcher的需求都是往Queue投递消息，所以只有一种被调度的函数
func (d *DemoDispatcher) taskFunc(task *model.Task) func() {
	return func() {
		// 每次分派一个Task，都更新状态存储器中对应的RunningTaskStatus条目
		taskStatus, err := d.updateRunningTaskStatus(task.Key)
		if err != nil {
			log.Printf("update running task status error : %v", err)
		}

		// 根据Task内容生成Job对象
		job := task.GenerateJob(uuid.NewString(), time.Now().Unix(), taskStatus.Seq)

		// 通过deliverier投递job
		err = d.deliverier.Deliver(job)
		if err != nil {
			log.Printf("deliver job error : %v", err)
		}
		log.Printf("delivered job : [%s]", task.Key)
	}
}
