package model

type JobResult struct {
	IsSuccess bool
	Logs      []string
	Result    interface{}
}
