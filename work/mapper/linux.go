package mapper

import (
	"encoding/json"
	"errors"
	"github.com/crazy-me/os_scheduler/common/entity"
	"strconv"
	"time"
)

type Linux struct {
	Base
}

func (l *Linux) Push(r *entity.JobExecuteResult) (err error) {
	// 解析任务结果
	var (
		jobResultParser entity.JobResultDto
		item            entity.AgentPush
	)

	if errs := l.ParserJob(r.Output, &jobResultParser); errs != nil {
		err = errors.New("command " + r.ExecStatus.Job.JobCommand +
			"code: " + strconv.Itoa(jobResultParser.Code) + "msg: " +
			jobResultParser.Msg + "|Base.ParserJob:" + errs.Error())
		return
	}

	// linux 采集数据
	linuxMap := jobResultParser.Data

	// 过滤不需要上报的指标
	l.FilterIndex(linuxMap, []string{"up_time"})

	for k, v := range linuxMap {
		if b := l.NoCalculationTsRequired(k); b {
			linuxMap[k] = l.ConvertToGb(v)
		}
	}

	pushList := make([]entity.AgentPush, 0)
	for queue, value := range linuxMap {
		item.Ident = r.ExecStatus.Job.JobIdent
		item.Alias = queue
		item.Metric = queue
		item.Time = time.Now().Unix()
		item.Value = value
		pushList = append(pushList, item)
	}

	body, err := json.Marshal(pushList)
	if err != nil {
		err = errors.New("command " + r.ExecStatus.Job.JobCommand + "push linux err: " + err.Error())
		return
	}
	err = l.HttpRequestByPost(body)

	return
}
