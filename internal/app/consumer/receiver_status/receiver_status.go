package receiver_status

// ReceiverStatus Receiver状态
type ReceiverStatus uint8

const (
	New         ReceiverStatus = iota // 刚创建，还没有初始化好，不能运行
	Initialized                       // 已经初始化好，但还没开始运行
	Running                           // 已经开始运行
	Stopped                           // 已经停止
	Shutdown                          // 已经关闭
)

func (s ReceiverStatus) String() string {
	switch s {
	case New:
		return "New"
	case Initialized:
		return "Initialized"
	case Running:
		return "Running"
	case Stopped:
		return "Stopped"
	case Shutdown:
		return "Shutdown"
	default:
		return "Unknown"
	}
}
