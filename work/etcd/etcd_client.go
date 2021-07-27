package etcd

import (
	"context"
	"github.com/crazy-me/os_scheduler/common/constants"
	"github.com/crazy-me/os_scheduler/common/entity"
	"github.com/crazy-me/os_scheduler/common/utils"
	"github.com/crazy-me/os_scheduler/work/conf"
	"github.com/crazy-me/os_scheduler/work/logger"
	"github.com/crazy-me/os_scheduler/work/scheduler"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"time"
)

var (
	Client *clientEtcd
)

type clientEtcd struct {
	conn  *clientv3.Client
	kv    clientv3.KV
	less  clientv3.Lease
	watch clientv3.Watcher
}

func InitEtcd() (err error) {
	var (
		config clientv3.Config
		conn   *clientv3.Client
		kv     clientv3.KV
		less   clientv3.Lease
		watch  clientv3.Watcher
	)
	config = clientv3.Config{
		Endpoints:   conf.C.Etcd.Endpoints,
		DialTimeout: time.Duration(conf.C.Etcd.Timeout) * time.Second,
	}

	if conn, err = clientv3.New(config); err != nil {
		logger.L.Error("etcd clientv3.New err:", zap.Any("connect", err))
		return
	}

	kv = clientv3.NewKV(conn)
	less = clientv3.NewLease(conn)
	watch = clientv3.NewWatcher(conn)

	Client = &clientEtcd{
		conn:  conn,
		kv:    kv,
		less:  less,
		watch: watch,
	}

	if err = Client.watchJobs(); err != nil {
		logger.L.Error("InitEtcd", zap.Any("Client.watchJobs", err))
	}
	return
}

// 初始化加载任务，并监听任务的变化
func (e *clientEtcd) watchJobs() (err error) {
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

	// 获取所有的任务
	if getResp, err = e.kv.Get(context.TODO(), constants.JOB_SAVE_DIR, clientv3.WithPrefix()); err != nil {
		return
	}

	// 遍历任务反序列化得到Job对象并投递给调度程序
	for _, kvPair = range getResp.Kvs {
		if job, err = utils.UnpackJob(kvPair.Value); err == nil {
			// TODO 构造JobEvent事件投递到调度程序
			jobEvent = &entity.JobEvent{
				EventType: constants.JOB_PUT_EVENT,
				Job:       job,
			}
			scheduler.ScheduleInstance.PushJobEvent(jobEvent)
		}
	}

	// 监听任务的变化来更新调度程序
	go func() {
		// 从revision向后监听事件变化
		watchStartRevision = getResp.Header.Revision + 1
		watchChan = e.watch.Watch(context.TODO(), constants.JOB_SAVE_DIR, clientv3.WithRev(watchStartRevision), clientv3.WithPrefix())
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
					job = &entity.Job{JobName: utils.ExtractJobKey(string(watchEvent.Kv.Key))}
					jobEvent = entity.BuildJobEvent(constants.JOB_DELETE_EVENT, job)
				}
				// TODO 将事件投递给调度程序
				scheduler.ScheduleInstance.PushJobEvent(jobEvent)
			}
		}
	}()

	return
}
