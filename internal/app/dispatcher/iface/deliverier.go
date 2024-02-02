package iface

import (
	"github.com/yellowb/simple-task-dispatch-demo/internal/model"
)

// Deliverier Job投递器接口
type Deliverier interface {
	// Deliver 往某个渠道投递一个Job
	Deliver(job *model.Job) error
	// Len 返回渠道中当前堆积了多少Job
	Len() (int64, error)
}
