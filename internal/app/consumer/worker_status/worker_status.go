package worker_status

// WorkerStatus Consumer状态
type WorkerStatus uint8

const (
	New         WorkerStatus = iota // 刚创建，还没有初始化好，不能运行
	Initialized                     // 已经初始化好，但还没开始运行
	Running                         // 已经开始运行
	Stopped                         // 已经停止
)

func (s WorkerStatus) String() string {
	switch s {
	case New:
		return "New"
	case Initialized:
		return "Initialized"
	case Running:
		return "Running"
	case Stopped:
		return "Stopped"
	default:
		return "Unknown"
	}
}
