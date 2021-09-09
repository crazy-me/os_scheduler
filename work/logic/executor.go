package logic

import (
	"github.com/crazy-me/os_scheduler/common/constants"
	"github.com/crazy-me/os_scheduler/common/entity"
	"github.com/crazy-me/os_scheduler/work/data_source/etcd"
	"github.com/crazy-me/os_scheduler/work/lock"
	"os/exec"
	"time"
)

var (
	ExecInstance *Executor
)

// Executor 执行器结构
type Executor struct {
}

// RunJob 执行任务
func (executor *Executor) RunJob(jobExecutor *entity.JobExecuteStatus) {
	go func() {
		var (
			err     error
			output  []byte
			cmd     *exec.Cmd
			result  *entity.JobExecuteResult
			jobKey  string
			jobLock *lock.JobLock
		)
		// 构建任务结果
		result = &entity.JobExecuteResult{
			ExecStatus: jobExecutor,
			Output:     make([]byte, 0),
		}

		// TODO 获取分布式锁
		jobKey = constants.JOB_LOCK_DIR + jobExecutor.Job.JobType +
			"/" + jobExecutor.Job.JobId
		jobLock = etcd.Cli.CreateJobLock(jobKey)
		// 开始时间
		result.StartTime = time.Now()
		// TODO 上锁
		err = jobLock.Lock()
		defer jobLock.Unlock()
		if err != nil { // 抢锁失败
			result.Err = err
			result.EndTime = time.Now()
		} else {
			// 重置开始时间
			result.StartTime = time.Now()
			cmd = exec.CommandContext(jobExecutor.CancelCtx, "/bin/bash", "-c", jobExecutor.Job.JobCommand)
			// 执行捕获输出
			output, err = cmd.CombinedOutput()
			// 结束时间
			result.EndTime = time.Now()
			result.Output = output
			result.Err = err
		}

		// TODO 将执行结果投递到调度器
		ScheduleInstance.PushJobResult(result)
	}()

}

// InitExecutor 初始化执行器
func InitExecutor() (err error) {
	ExecInstance = &Executor{}
	return
}
