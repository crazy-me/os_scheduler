package mapper

import "github.com/crazy-me/os_scheduler/common/entity"

type Win struct {
	Base
}

func (w *Win) Push(r *entity.JobExecuteResult) (err error) {
	return
}
