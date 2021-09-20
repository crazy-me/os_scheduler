package mapper

import "github.com/crazy-me/os_scheduler/common/entity"

type Linux struct {
	Base
}

func (l *Linux) Push(r *entity.JobExecuteResult) (err error) {
	return
}
