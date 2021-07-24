package entity

// Job 任务
type Job struct {
	JobId      int    `json:"job_id"`
	JobName    string `json:"job_name"`
	JobType    string `json:"job_type"`
	JobCommand string `json:"job_command"`
	JobExpr    int    `json:"job_expr"`
}
