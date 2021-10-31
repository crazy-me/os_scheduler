package logic

import (
	"github.com/crazy-me/os_scheduler/common/entity"
	"github.com/crazy-me/os_scheduler/work/conf"
)

var (
	TaskResult *Record
)

// Record 任务结果记录
type Record struct {
	TaskResultChan chan *entity.JobExecuteResult
}

// InitRecord 初始化
func InitRecord() {
	TaskResult = &Record{
		TaskResultChan: make(chan *entity.JobExecuteResult, conf.C.JobEventChan),
	}
	go TaskResultLoop()
	return
}

// PushTaskResult
// TODO 接收调度器的任务结果
func (record *Record) PushTaskResult(jobResult *entity.JobExecuteResult) {
	record.TaskResultChan <- jobResult
}
