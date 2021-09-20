package logic

import (
	"github.com/crazy-me/os_scheduler/common/entity"
	"github.com/crazy-me/os_scheduler/work/logger"
	"github.com/crazy-me/os_scheduler/work/mapper"
	"go.uber.org/zap"
)

var resourceTypeMap map[string]ResourcePublish

// ResourcePublish 资源数据发布接口
type ResourcePublish interface {
	Push(r *entity.JobExecuteResult) error
}

func init() {
	resourceTypeMap = make(map[string]ResourcePublish)
	resourceTypeMap["network"] = &mapper.Network{}
	resourceTypeMap["linux"] = &mapper.Linux{}
	resourceTypeMap["windows"] = &mapper.Win{}
}

func TaskResultLoop() {
	for {
		select {
		case resource := <-TaskResult.TaskResultChan:
			t := resource.ExecStatus.Job.JobType
			err := resourceTypeMap[t].Push(resource)
			if err != nil {
				// TODO 日志
				logger.L.Error("taskResultLoop", zap.Any("push", err))
			}
		}
	}
}
