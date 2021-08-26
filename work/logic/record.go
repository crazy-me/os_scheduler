package logic

import (
	"github.com/crazy-me/os_scheduler/common/entity"
	"github.com/crazy-me/os_scheduler/work/conf"
	"github.com/crazy-me/os_scheduler/work/logger"
	"go.uber.org/zap"
)

var (
	TaskResult *Record
)

// Record 任务结果记录
type Record struct {
	TaskResultChan chan *entity.JobExecuteResult
}

// WriteTaskResultLoop 记录任务结果
func (record *Record) WriteTaskResultLoop() {
	var (
		jobResult *entity.JobExecuteResult // 任务结果对象
		jobInfo   *entity.Job              // 任务信息对象
	)

	for {
		select {
		case jobResult = <-record.TaskResultChan:
			jobInfo = jobResult.ExecStatus.Job
			// 任务执行错误记录日志并跳过
			if jobResult.Err != nil {
				logger.L.Error("record.WriteTaskResultLoop", zap.Any(string(jobResult.Output), jobInfo))
				continue
			}
			jobAgent := &Agent{JobInfo: jobResult}
			switch jobInfo.JobType {
			case "network":
				jobAgent.Func = jobAgent.PushNetwork
			case "server":
				jobAgent.Func = jobAgent.PushServer
			case "mysql":
				jobAgent.Func = jobAgent.PushMysql
			case "apply":
				jobAgent.Func = jobAgent.PushApply
			default: // 防止 jobAgent.Func 为nil
				jobAgent.Func = jobAgent.PushEmpty
			}

			pushErr := jobAgent.Func()
			if pushErr != nil {
				logger.L.Error("record.WriteTaskResultLoop", zap.Any("http push err:", pushErr))
			}
		}
	}
}

// InitRecord 初始化
func InitRecord() {
	TaskResult = &Record{
		TaskResultChan: make(chan *entity.JobExecuteResult, conf.C.JobEventChan),
	}
	go TaskResult.WriteTaskResultLoop()
	return
}

// PushTaskResult 投递任务结果记录
func (record *Record) PushTaskResult(jobResult *entity.JobExecuteResult) {
	record.TaskResultChan <- jobResult
}
