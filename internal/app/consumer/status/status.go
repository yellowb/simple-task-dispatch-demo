package status

type Status uint8

const (
	New Status = iota
	Initialized
	Running
	Stopped
)

func (s Status) String() string {
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
