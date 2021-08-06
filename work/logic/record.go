package logic

import (
	"github.com/crazy-me/os_scheduler/common/entity"
	"github.com/crazy-me/os_scheduler/work/conf"
	"github.com/crazy-me/os_scheduler/work/data_source/redis"
	"github.com/crazy-me/os_scheduler/work/logger"
	"go.uber.org/zap"
	"strconv"
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
		jobResult *entity.JobExecuteResult
		isWrite   bool
	)
	redisClient := redis.Pool()
	defer redisClient.ReleaseRedisClient()
	for {
		select {
		case jobResult = <-record.TaskResultChan:
			isWrite = redisClient.Hset(jobResult.ExecStatus.Job.JobType,
				strconv.Itoa(jobResult.ExecStatus.Job.JobId),
				string(jobResult.Output))
			if !isWrite { // 记录失败
				logger.L.Info("WriteTaskResultLoop:"+jobResult.ExecStatus.Job.JobCommand, zap.Any("info", "job write fail"))
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
