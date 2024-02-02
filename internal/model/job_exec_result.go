package model

// JobExecResult 一个Job的执行结果，一般由Handler返回
type JobExecResult struct {
	Succeed bool        `json:"succeed"` // 是否成功
	ErrMsg  string      `json:"err_msg"` // 错误信息（成功时为空）
	Logs    []string    `json:"logs"`    // Job执行中间输出的日志
	Output  interface{} `json:"output"`  // Job的结果
}
