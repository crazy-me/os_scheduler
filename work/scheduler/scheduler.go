package scheduler

import (
	"github.com/crazy-me/os_scheduler/common/constants"
	"github.com/crazy-me/os_scheduler/common/entity"
	"github.com/crazy-me/os_scheduler/work/conf"
	"github.com/crazy-me/os_scheduler/work/executor"
	"github.com/gorhill/cronexpr"
	"strconv"
	"time"
)

// Scheduler Job调度
type Scheduler struct {
	jobEventChan      chan *entity.JobEvent               // JobEventQueue
	jobPlanTable      map[string]*entity.JobSchedulerPlan // 任务调度计划表
	jobExecutingTable map[string]*entity.JobExecuteStatus // 任务调度执行表
}

var (
	ScheduleInstance *Scheduler
)

func (scheduler *Scheduler) SchedulerLoop() {
	var (
		jobEvent       *entity.JobEvent
		schedulerAfter time.Duration
		schedulerTimer *time.Timer
		jobResult      *entity.JobExecuteResult
	)

	// 初始化调度
	schedulerAfter = scheduler.schedulerJobPlanTime()
	schedulerTimer = time.NewTimer(schedulerAfter)

	for {
		select {
		case jobEvent = <-scheduler.jobEventChan:
			scheduler.HandleJobEvent(jobEvent)
		case <-schedulerTimer.C: // 最近任务到期
		case jobResult = <-entity.JobResultChan: // 监听任务结束
			scheduler.HandleJobResult(jobResult)
		}
		// 调度一次
		schedulerAfter = scheduler.schedulerJobPlanTime()
		// 重置定时器
		schedulerTimer.Reset(schedulerAfter)
	}
}

// HandleJobEvent 处理任务事件
func (scheduler *Scheduler) HandleJobEvent(jobEvent *entity.JobEvent) {
	var (
		err              error
		jobExists        bool
		jobSchedulerPlan *entity.JobSchedulerPlan
		jobPlanTableKey  string
	)
	switch jobEvent.EventType {
	case constants.JOB_PUT_EVENT: // 保存任务事件
		jobPlanTableKey = jobEvent.Job.JobType + "/" + strconv.Itoa(jobEvent.Job.JobId)
		if jobSchedulerPlan, err = scheduler.buildJobSchedulerPlan(jobEvent.Job); err != nil {
			return
		}
		scheduler.jobPlanTable[jobPlanTableKey] = jobSchedulerPlan
	case constants.JOB_DELETE_EVENT: // 删除任务事件
		if jobSchedulerPlan, jobExists = scheduler.jobPlanTable[jobEvent.Job.JobName]; jobExists {
			delete(scheduler.jobPlanTable, jobEvent.Job.JobName)
		}
	}
}

// HandleJobResult 任务执行结果,从任务执行表中移除本次执行完的任务
func (scheduler *Scheduler) HandleJobResult(jobResult *entity.JobExecuteResult) {
	jobKey := jobResult.ExecStatus.Job.JobType + "/" +
		strconv.Itoa(jobResult.ExecStatus.Job.JobId)
	delete(scheduler.jobExecutingTable, jobKey)
	//fmt.Println("任务名称:", jobResult.ExecStatus.Job.JobName)
	//fmt.Println("任务命令:", jobResult.ExecStatus.Job.JobCommand)
	//fmt.Println("任务执行结果:", string(jobResult.Output), jobResult.Err)
}

// PushJobEvent 推送Job事件
func (scheduler *Scheduler) PushJobEvent(jobEvent *entity.JobEvent) {
	scheduler.jobEventChan <- jobEvent
}

// buildJobSchedulerPlan 构建执行计划
func (scheduler *Scheduler) buildJobSchedulerPlan(job *entity.Job) (jobSchedulerPlan *entity.JobSchedulerPlan, err error) {
	var expr *cronexpr.Expression
	// 解析Cron表达式
	if expr, err = cronexpr.Parse(job.JobExpr); err != nil {
		return
	}

	// 生成调度计划
	jobSchedulerPlan = &entity.JobSchedulerPlan{
		Job:      job,
		Expr:     expr,
		NextTime: expr.Next(time.Now()),
	}
	return
}

func (scheduler *Scheduler) RunSchedulerJob(jobPlan *entity.JobSchedulerPlan) {
	var (
		jobExecute   *entity.JobExecuteStatus
		jobIsExecute bool
		jobKey       string
	)
	// 一个任务可能执行很长时间,如果任务正在执行则跳过本次调度
	jobKey = jobPlan.Job.JobType + "/" + strconv.Itoa(jobPlan.Job.JobId)
	// 执行调度表有此任务则跳过本次调度
	if jobExecute, jobIsExecute = scheduler.jobExecutingTable[jobKey]; jobIsExecute {
		return
	}
	// 构建执行状态
	jobExecute = entity.BuildJobExecStatus(jobPlan)
	// 记录执行状态
	scheduler.jobExecutingTable[jobKey] = jobExecute
	// TODO 调用执行器
	executor.ExecInstance.RunExecuteJob(jobExecute)
}

// 计算任务调度时间
func (scheduler *Scheduler) schedulerJobPlanTime() (schedulerAfter time.Duration) {
	var (
		jobPlan  *entity.JobSchedulerPlan
		nowTime  time.Time
		nearTime *time.Time
	)

	// 如果当前没有任务
	if len(scheduler.jobPlanTable) == 0 {
		schedulerAfter = 1 * time.Second
		return
	}

	nowTime = time.Now()
	// 遍历当前任务
	for _, jobPlan = range scheduler.jobPlanTable {
		if jobPlan.NextTime.Before(nowTime) || jobPlan.NextTime.Equal(nowTime) {
			// TODO  尝试执行任务
			scheduler.RunSchedulerJob(jobPlan)
			// 更新下次执行时间
			jobPlan.NextTime = jobPlan.Expr.Next(nowTime)
		}

		// 统计最近要过期的任务时间
		if nearTime == nil || jobPlan.NextTime.Before(*nearTime) {
			nearTime = &jobPlan.NextTime
		}
	}
	// 下次调度间隔(最近要执行的调度时间-当前时间)
	schedulerAfter = (*nearTime).Sub(nowTime)
	return
}

// InitSchedule 初始化调度
func InitSchedule() (err error) {
	entity.JobResultChan = make(chan *entity.JobExecuteResult, conf.C.JobEventChan)
	ScheduleInstance = &Scheduler{
		jobEventChan:      make(chan *entity.JobEvent, conf.C.JobEventChan),
		jobPlanTable:      make(map[string]*entity.JobSchedulerPlan),
		jobExecutingTable: make(map[string]*entity.JobExecuteStatus, conf.C.JobEventChan),
	}
	// 拉取调度协程
	go ScheduleInstance.SchedulerLoop()
	return
}
