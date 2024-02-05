package iface

type Worker interface {
	//注入TaskDataStorage
	TaskResultStorage(storage TaskResultStorage) Worker
	//注入Consumer
	Consumer(consumer Consumer) Worker

	//初始化
	Init() error
	//启动
	Run() error
	//关闭
	Shutdown() error
}
