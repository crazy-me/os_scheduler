package lock

import (
	"context"
	"errors"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	jobLock *JobLock
)

type JobLock struct {
	kv        clientv3.KV
	less      clientv3.Lease
	jobName   string
	canelFunc context.CancelFunc //终止自动租约
	leaseId   clientv3.LeaseID   // 租约ID
	isLock    bool
}

// Lock 尝试上一把etcd锁
func (lock *JobLock) Lock() (err error) {
	var (
		lessResp      *clientv3.LeaseGrantResponse
		ctx           context.Context
		ctxFunc       context.CancelFunc
		keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
		txn           clientv3.Txn
		txnResp       *clientv3.TxnResponse
	)
	// 创建一个etcd租约5秒过期
	if lessResp, err = lock.less.Grant(context.TODO(), 5); err != nil {
		return
	}

	// 如果抢锁成功则需要续租 context用于取消自动续租
	ctx, ctxFunc = context.WithCancel(context.TODO())
	keepAliveChan, err = lock.less.KeepAlive(ctx, lessResp.ID)
	if err != nil {
		goto FAIL
	}

	// 处理续租应答
	go func() {
		for {
			select {
			case keepResp := <-keepAliveChan:
				if keepResp == nil {
					goto END
				}
			}
		}
	END:
	}()

	// 创建一个事物
	txn = lock.kv.Txn(context.TODO())
	txn.If(clientv3.Compare(clientv3.CreateRevision(lock.jobName), "=", 0)).
		Then(clientv3.OpPut(lock.jobName, "", clientv3.WithLease(lessResp.ID))).
		Else(clientv3.OpGet(lock.jobName))

	// 提交事物
	if txnResp, err = txn.Commit(); err != nil {
		goto FAIL
	}

	if !txnResp.Succeeded { // 锁被占用
		err = errors.New("锁被占用")
		goto FAIL
	}

	// 抢锁成功
	lock.leaseId = lessResp.ID
	lock.canelFunc = ctxFunc
	lock.isLock = true
	return

FAIL:
	// 取消自动续租，释放租约
	ctxFunc()
	_, _ = lock.less.Revoke(context.TODO(), lessResp.ID)
	return
}

// Unlock 释放锁
func (lock *JobLock) Unlock() {
	if lock.isLock {
		lock.canelFunc()                                      //取消自动续租的协程
		_, _ = lock.less.Revoke(context.TODO(), lock.leaseId) // 释放租约
	}
}

// InitJobLock 初始化锁
func InitJobLock(jobName string, kv clientv3.KV, less clientv3.Lease) (jobLock *JobLock) {
	jobLock = &JobLock{
		kv:      kv,
		less:    less,
		jobName: jobName,
	}
	return
}
