package executor

import (
	"context"
	"github.com/crazy-me/os_scheduler/common/entity"
	"os/exec"
	"time"
)

var (
	ExecInstance *Executor
)

type Executor struct {
}

func (executor *Executor) RunExecuteJob(jobExecutor *entity.JobExecuteStatus) {
	// 拉取协程执行任务
	go func() {
		var (
			err    error
			output []byte
			cmd    *exec.Cmd
			result *entity.JobExecuteResult
		)

		result = &entity.JobExecuteResult{
			ExecStatus: jobExecutor,
			Output:     make([]byte, 0),
		}
		// 开始时间
		result.StartTime = time.Now()
		cmd = exec.CommandContext(context.TODO(), "/bin/bash", "-c", jobExecutor.Job.JobCommand)
		// 执行捕获输出
		output, err = cmd.CombinedOutput()
		// 结束时间
		result.EndTime = time.Now()
		result.Output = output
		result.Err = err

		// 将执行结果返回到调度程序，调度程序移除本次执行记录
		executor.PushJobResult(result)
	}()
}

// PushJobResult 投递任务结果
func (executor *Executor) PushJobResult(jobResult *entity.JobExecuteResult) {
	entity.JobResultChan <- jobResult
}

func InitExecutor() (err error) {
	ExecInstance = &Executor{}
	return
}
