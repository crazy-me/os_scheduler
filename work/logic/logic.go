package logic

import (
	"context"
	"github.com/crazy-me/os_scheduler/common/constants"
	"github.com/crazy-me/os_scheduler/common/entity"
	"github.com/crazy-me/os_scheduler/common/utils"
	etcd2 "github.com/crazy-me/os_scheduler/work/data_source/etcd"
	"github.com/crazy-me/os_scheduler/work/logger"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"os"
)

// WatchJobs 监听Etcd中任务的行为实时更新任务到内存
func WatchJobs() (err error) {
	// 初始化变量
	var (
		getResp            *clientv3.GetResponse
		kvPair             *mvccpb.KeyValue
		job                *entity.Job
		watchStartRevision int64
		watchChan          clientv3.WatchChan
		watchResp          clientv3.WatchResponse
		watchEvent         *clientv3.Event
		jobEvent           *entity.JobEvent
	)

	// 初始化所有任务并加载到内存
	if getResp, err = etcd2.Cli.Kv.Get(context.TODO(), constants.JOB_SAVE_DIR, clientv3.WithPrefix()); err != nil {
		logger.L.Error("logic-WatchJobs", zap.Any("etcd.Cli.Kv.Get", err))
		os.Exit(-1)
	}

	// 遍历任务反序列化得到Job对象并投递给调度器
	for _, kvPair = range getResp.Kvs {
		if job, err = utils.UnpackJob(kvPair.Value); err == nil {
			// TODO 构造JobEvent事件对象
			jobEvent = &entity.JobEvent{
				EventType: constants.JOB_PUT_EVENT,
				Job:       job,
			}
			// TODO 将事件对象投递到调度器
			ScheduleInstance.PushJobEvent(jobEvent)
		}
	}

	// 开启协程监听Etcd任务变化，实时更新内存中的任务，
	//保证加载到内存中的任务为最新
	go func() {
		// 从revision向后监听事件变化
		watchStartRevision = getResp.Header.Revision + 1
		watchChan = etcd2.Cli.Watch.Watch(context.TODO(), constants.JOB_SAVE_DIR, clientv3.WithRev(watchStartRevision), clientv3.WithPrefix())
		// 处理监听事件
		for watchResp = range watchChan {
			for _, watchEvent = range watchResp.Events {
				switch watchEvent.Type {
				case mvccpb.PUT: // 任务保存事件
					// 反序列化Job
					if job, err = utils.UnpackJob(watchEvent.Kv.Value); err != nil {
						continue
					}
					// TODO 构造一个更新事件
					jobEvent = entity.BuildJobEvent(constants.JOB_PUT_EVENT, job)
				case mvccpb.DELETE: // 任务删除事件
					// TODO 构造一个删除事件
					job = &entity.Job{JobName: utils.ExtractJobKey(constants.JOB_SAVE_DIR, string(watchEvent.Kv.Key))}
					jobEvent = entity.BuildJobEvent(constants.JOB_DELETE_EVENT, job)
				}
				// TODO 将事件投递给调度器
				ScheduleInstance.PushJobEvent(jobEvent)
			}
		}
	}()

	// TODO 监听任务killer事件，强制杀死任务
	go func() {
		var (
			job                 *entity.Job
			jobEvent            *entity.JobEvent
			jobKillerChan       clientv3.WatchChan
			jobKillerWatchResp  clientv3.WatchResponse
			jobKillerWatchEvent *clientv3.Event
		)
		jobKillerChan = etcd2.Cli.Watch.Watch(context.TODO(), constants.JOB_KILLER_DIR, clientv3.WithPrefix())

		// 处理监听事件
		for jobKillerWatchResp = range jobKillerChan {
			for _, jobKillerWatchEvent = range jobKillerWatchResp.Events {
				switch jobKillerWatchEvent.Type {
				case mvccpb.PUT: // Job Killer Event
					job = &entity.Job{JobName: utils.ExtractJobKey(constants.JOB_KILLER_DIR, string(jobKillerWatchEvent.Kv.Key))}
					jobEvent = entity.BuildJobEvent(constants.JOB_KILLER_EVENT, job)
					ScheduleInstance.PushJobEvent(jobEvent)
				case mvccpb.DELETE: // killer 过期自动删除
				}
			}
		}
	}()

	return
}
