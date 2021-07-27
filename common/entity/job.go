package entity

import (
	"github.com/gorhill/cronexpr"
	"time"
)

var JobResultChan chan *JobExecuteResult

// Job 任务
type Job struct {
	JobId      int    `json:"job_id"`
	JobName    string `json:"job_name"`
	JobType    string `json:"job_type"`
	JobCommand string `json:"job_command"`
	JobExpr    string `json:"job_expr"`
}

// JobSchedulerPlan 任务调度计划
type JobSchedulerPlan struct {
	Job      *Job                 // 任务信息
	Expr     *cronexpr.Expression // cron表达式
	NextTime time.Time            // 下次调用时间
}

// JobExecuteStatus 任务执行状态
type JobExecuteStatus struct {
	Job      *Job
	PlanTime time.Time // 理论计划执行时间
	RealTime time.Time // 实际计划执行时间
}

// JobExecuteResult 任务执行结果
type JobExecuteResult struct {
	ExecStatus *JobExecuteStatus
	Err        error
	Output     []byte
	StartTime  time.Time
	EndTime    time.Time
}

// JobEvent 任务事件
type JobEvent struct {
	EventType int
	Job       *Job
}

func BuildJobEvent(eventType int, job *Job) (jobEvent *JobEvent) {
	jobEvent = &JobEvent{
		EventType: eventType,
		Job:       job,
	}
	return
}

// BuildJobExecStatus 构造任务执行状态
func BuildJobExecStatus(jobPlan *JobSchedulerPlan) (jobExecuteStatus *JobExecuteStatus) {
	jobExecuteStatus = &JobExecuteStatus{
		Job:      jobPlan.Job,
		PlanTime: jobPlan.NextTime,
		RealTime: time.Now(),
	}
	return
}
