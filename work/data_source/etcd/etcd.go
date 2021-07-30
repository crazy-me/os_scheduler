package etcd

import (
	"context"
	"github.com/crazy-me/os_scheduler/work/conf"
	"github.com/crazy-me/os_scheduler/work/lock"
	"github.com/crazy-me/os_scheduler/work/logger"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"os"
	"time"
)

var (
	Cli *Etcd // Etcd 全局连接实例
)

// Etcd 实例结构
type Etcd struct {
	conn  *clientv3.Client
	Kv    clientv3.KV
	Less  clientv3.Lease
	Watch clientv3.Watcher
}

// InitEtcd 初始化Etcd连接实例
func InitEtcd() (err error) {
	var (
		config clientv3.Config  // Etcd 配置
		conn   *clientv3.Client // Etcd 连接实例
		kv     clientv3.KV
		less   clientv3.Lease
		watch  clientv3.Watcher
	)
	// 初始化Etcd配置项
	config = clientv3.Config{
		Endpoints:   conf.C.Etcd.Endpoints,
		DialTimeout: time.Duration(conf.C.Etcd.Timeout) * time.Second,
	}
	// 连接Etcd
	if conn, err = clientv3.New(config); err != nil {
		logger.L.Error("etcd client.New err:", zap.Any("connect", err))
		os.Exit(-1)
	}
	// 3.X版本使用新的均衡器client.New不会抛出链接错误,使用Status函数来检测当前的连接状态
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err = conn.Status(ctx, config.Endpoints[0])
	if err != nil { // 连接超时
		logger.L.Error("etcd connection", zap.Any("err", err))
		os.Exit(-1)
	}

	kv = clientv3.NewKV(conn)         // k-v 操作
	less = clientv3.NewLease(conn)    // 租约操作
	watch = clientv3.NewWatcher(conn) // 监听操作

	// 构建Etcd操作对象
	Cli = &Etcd{
		conn:  conn,
		Kv:    kv,
		Less:  less,
		Watch: watch,
	}
	return
}

func (e *Etcd) CreateJobLock(jobName string) (jobLock *lock.JobLock) {
	jobLock = lock.InitJobLock(jobName, e.Kv, e.Less)
	return
}
