package iface

import "github.com/yellowb/simple-task-dispatch-demo/internal/app/dispatcher/model"

// Deliverier Job投递器接口
type Deliverier interface {
	// Deliver 往某个渠道投递一个Job
	Deliver(job *model.Job) error
}
