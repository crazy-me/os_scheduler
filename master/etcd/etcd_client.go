package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/crazy-me/os_scheduler/common/constants"
	"github.com/crazy-me/os_scheduler/common/entity"
	"github.com/crazy-me/os_scheduler/master/conf"
	"github.com/crazy-me/os_scheduler/master/logger"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"os"
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
		DialTimeout: time.Duration(conf.C.Etcd.Timeout) * time.Second,
	}

	if conn, err = clientv3.New(config); err != nil {
		logger.L.Error("etcd clientv3.New err:", zap.Any("connect", err))
		return
	}

	// 3.X版本使用新的均衡器clientv3.New不会抛出链接错误,使用Status函数来检测当前的连接状态
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err = conn.Status(ctx, config.Endpoints[0])
	if err != nil {
		logger.L.Error("etcd connection", zap.Any("err", err))
		os.Exit(-1)
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

// KillJob 强杀Job
func (e *clientEtcd) KillJob(job *entity.Job) (err error) {
	var (
		killJobKey string
		lessGrant  *clientv3.LeaseGrantResponse
	)
	// 获取一个带有效期的租约
	if lessGrant, err = e.less.Grant(context.TODO(), 1); err != nil {
		return
	}

	killJobKey = fmt.Sprintf(constants.JOB_KILLER_DIR+"%s/%d", job.JobType, job.JobId)
	// 向etcd投递事件来通知Work杀掉当前任务
	if _, err = e.kv.Put(context.TODO(), killJobKey, "", clientv3.WithLease(lessGrant.ID)); err != nil {
		return
	}
	return
}
