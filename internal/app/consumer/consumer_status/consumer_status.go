package consumer_status

// ConsumerStatus Consumer状态
type ConsumerStatus uint8

const (
	New         ConsumerStatus = iota // 刚创建，还没有初始化好，不能运行
	Initialized                       // 已经初始化好，但还没开始运行
	Running                           // 已经开始运行
	Shutdown                          // 已经关闭，不能再运行
)

func (s ConsumerStatus) String() string {
	switch s {
	case New:
		return "New"
	case Initialized:
		return "Initialized"
	case Running:
		return "Running"
	case Shutdown:
		return "Shutdown"
	default:
		return "Unknown"
	}
}
