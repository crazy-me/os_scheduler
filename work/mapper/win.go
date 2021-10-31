package mapper

import (
	"encoding/json"
	"errors"
	"github.com/crazy-me/os_scheduler/common/entity"
	"strconv"
	"time"
)

type Win struct {
	Base
}

func (w *Win) Push(r *entity.JobExecuteResult) (err error) {
	// 解析任务结果
	var (
		jobResultParser entity.JobResultDto
		item            entity.AgentPush
	)

	if errs := w.ParserJob(r.Output, &jobResultParser); errs != nil {
		err = errors.New("command " + r.ExecStatus.Job.JobCommand +
			"code: " + strconv.Itoa(jobResultParser.Code) + "msg: " +
			jobResultParser.Msg + "|Base.ParserJob:" + errs.Error())
		return
	}

	// win 采集数据
	winMap := jobResultParser.Data

	// 过滤不需要上报的指标
	w.FilterIndex(winMap, []string{"up_time"})

	for k, v := range winMap {
		if b := w.NoCalculationTsRequired(k); b {
			winMap[k] = w.ConvertToGb(v)
		}
	}

	pushList := make([]entity.AgentPush, 0)
	for queue, value := range winMap {
		item.Ident = r.ExecStatus.Job.JobIdent
		item.Alias = queue
		item.Metric = queue
		item.Time = time.Now().Unix()
		item.Value = value
		pushList = append(pushList, item)
	}

	body, err := json.Marshal(pushList)

	if err != nil {
		err = errors.New("command " + r.ExecStatus.Job.JobCommand + "push windows err: " + err.Error())
		return
	}
	err = w.HttpRequestByPost(body)

	return
}
