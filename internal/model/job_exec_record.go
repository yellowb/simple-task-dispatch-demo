package model

// JobExecRecord 一个Job的执行记录
type JobExecRecord struct {
	// Handler返回的执行结果
	Result *JobExecResult `json:"result"`
	// 被执行的Job信息
	Job *Job `json:"job"`
	// 其它统计信息
	StartTime int64 `json:"start_time"` // Job真正被执行的时间
	EndTime   int64 `json:"end_time"`   // Job结束时间
	TimeCost  int64 `json:"time_cost"`  // Job执行耗时：EndTime - StartTime
	Delay     int64 `json:"delay"`      // Job开始执行的延迟：StartTime - DispatchedTime
}
