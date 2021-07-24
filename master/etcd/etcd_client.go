package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/crazy-me/os_scheduler/common/constants"
	"github.com/crazy-me/os_scheduler/common/logger"
	"github.com/crazy-me/os_scheduler/master/conf"
	"github.com/crazy-me/os_scheduler/master/entity"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"time"
)

var (
	Client *clientEtcd
)

type clientEtcd struct {
	conn *clientv3.Client
	kv   clientv3.KV
	less clientv3.Lease
}

func InitEtcd() (err error) {
	var (
		config clientv3.Config
		conn   *clientv3.Client
		kv     clientv3.KV
		less   clientv3.Lease
	)
	config = clientv3.Config{
		Endpoints:   conf.C.Etcd.Endpoints,
		DialTimeout: time.Duration(conf.C.Etcd.Timeout) * time.Millisecond,
	}

	if conn, err = clientv3.New(config); err != nil {
		logger.L.Error("etcd clientv3.New err:", zap.Any("connect", err))
		return
	}
	kv = clientv3.NewKV(conn)
	less = clientv3.NewLease(conn)

	Client = &clientEtcd{
		conn: conn,
		kv:   kv,
		less: less,
	}
	return
}

// SaveJob 保存任务(如果是更新操作则返回旧值给客户端)
func (e *clientEtcd) SaveJob(job *entity.Job) (oldJob *entity.Job, err error) {
	var (
		jobKey   string
		jobValue []byte
		resp     *clientv3.PutResponse
	)

	jobKey = fmt.Sprintf(constants.JOB_SAVE_DIR+"%s/%d", job.JobType, job.JobId)
	fmt.Println(jobKey)
	if jobValue, err = json.Marshal(job); err != nil {
		return
	}
	// put时返回旧值，带clientv3.WithPrevKV() 选项
	if resp, err = e.kv.Put(context.TODO(), jobKey, string(jobValue), clientv3.WithPrevKV()); err != nil {
		return
	}
	// 如果是更新则返回旧的值
	if resp.PrevKv != nil {
		_ = json.Unmarshal(resp.PrevKv.Value, &oldJob)
	}
	return
}

// DeleteJob 删除Job 删除成功则返回被删除的Job
func (e *clientEtcd) DeleteJob(job *entity.Job) (oldJob *entity.Job, err error) {
	var (
		jobKey string
		resp   *clientv3.DeleteResponse
	)
	fmt.Println(jobKey)
	jobKey = fmt.Sprintf(constants.JOB_SAVE_DIR+"%s/%d", job.JobType, job.JobId)
	if resp, err = e.kv.Delete(context.TODO(), jobKey, clientv3.WithPrevKV()); err != nil {
		return
	}

	// 如果删除成功则取出被删除的值
	if len(resp.PrevKvs) != 0 {
		_ = json.Unmarshal(resp.PrevKvs[0].Value, &oldJob)
	}
	return
}

// ListJob 获取所有job
func (e *clientEtcd) ListJob() (listJob []*entity.Job, err error) {
	var (
		jobKey   string
		respList *clientv3.GetResponse
		tmpJob   *entity.Job
		kvPair   *mvccpb.KeyValue
	)

	jobKey = constants.JOB_SAVE_DIR
	if respList, err = e.kv.Get(context.TODO(), jobKey, clientv3.WithPrefix()); err != nil {
		return
	}

	// 初始化listJob
	listJob = make([]*entity.Job, 0)
	for _, kvPair = range respList.Kvs {
		tmpJob = &entity.Job{}
		_ = json.Unmarshal(kvPair.Value, tmpJob)
		listJob = append(listJob, tmpJob)
	}
	return
}
