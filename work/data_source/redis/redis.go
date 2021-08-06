package redis

import (
	"errors"
	"fmt"
	"github.com/crazy-me/os_scheduler/work/conf"
	"github.com/crazy-me/os_scheduler/work/logger"
	"github.com/gomodule/redigo/redis"
	"go.uber.org/zap"
	"time"
)

var redisPool *redis.Pool

type redisClient struct {
	client redis.Conn
}

func InitRedis() {
	fmt.Println("redis pool connection init")
	redisPool = redisClientPool()
}

// Pool 获取连接
func Pool() *redisClient {
	maxRetryTimes := conf.C.Redis.MaxRetryTimes
	var oneConn redis.Conn
	for i := 1; i <= maxRetryTimes; i++ {
		oneConn = redisPool.Get()
		if oneConn.Err() != nil {
			logger.L.Info("[utils-redis-Pool]", zap.Any("error", oneConn.Err()))
			if i == maxRetryTimes {
				return nil
			}
			time.Sleep(1 * time.Second)
		} else {
			break
		}
	}
	return &redisClient{oneConn}
}

// 连接池
func redisClientPool() *redis.Pool {
	redisPool = &redis.Pool{
		MaxIdle:     conf.C.Redis.MaxIdle,
		MaxActive:   conf.C.Redis.MaxActive,
		IdleTimeout: 60 * time.Second,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", conf.C.Redis.Addr)
			if err != nil {
				panic(err)
			}
			auth := conf.C.Redis.Password
			if len(auth) >= 1 {
				if _, err := conn.Do("AUTH", auth); err != nil {
					fmt.Println(err)
					_ = conn.Close()
				}
			}
			_, _ = conn.Do("select", conf.C.Redis.DB)
			return conn, err
		},
	}
	return redisPool
}

// ReleaseRedisClient 释放连接到连接池
func (r *redisClient) ReleaseRedisClient() {
	_ = r.client.Close()
}

// Set 缓存字符串
func (r *redisClient) Set(k string, v interface{}) bool {
	_, err := r.client.Do("SET", k, v)
	if err != nil {
		return false
	}
	return true
}

// SetExpire 带过期时间缓存字符串
func (r *redisClient) SetExpire(k string, v interface{}, expiration int64) bool {
	_, err := r.client.Do("SET", k, v, "EX", expiration)
	if err != nil {
		return false
	}
	return true
}

// Get 获取缓存字符串
func (r *redisClient) Get(k string) (string, error) {
	reply, err := r.client.Do("GET", k)
	if err != nil {
		return "", err
	}
	return redis.String(reply, err)
}

// ExistsKey 缓存标识是否存在
func (r *redisClient) ExistsKey(k string) bool {
	reply, _ := r.client.Do("EXISTS", k)
	if reply == int64(0) {
		return false
	}
	return true
}

// DelKey 删除缓存标识
func (r *redisClient) DelKey(k string) bool {
	reply, _ := r.client.Do("DEL", k)
	if reply == int64(0) {
		return false
	}
	return true
}

// ExpireKey 设置缓存标识的过期时间
func (r *redisClient) ExpireKey(k string, expiration int64) bool {
	reply, _ := r.client.Do("EXPIRE", k, expiration)
	if reply == int64(0) {
		return false
	}
	return true
}

// Hset 设置哈希
func (r *redisClient) Hset(k, f string, v interface{}) bool {
	_, err := r.client.Do("HSET", k, f, v)
	if err != nil {
		return false
	}
	return true
}

// Hget 获取哈希中的字段值
func (r *redisClient) Hget(k, f string) (string, error) {
	if ok := r.Hexists(k, f); !ok {
		return "", errors.New(f + " is not found")
	}
	reply, err := r.client.Do("hget", k, f)
	if err != nil {
		return "", err
	}
	return redis.String(reply, err)
}

// Hdel 删除哈希中的字段值
func (r *redisClient) Hdel(k, f string) bool {
	_, err := r.client.Do("HDEL", k, f)
	if err != nil {
		return false
	}
	return true
}

// Hexists 判断哈希中是否有某个字段
func (r *redisClient) Hexists(k, f string) bool {
	reply, _ := r.client.Do("hexists", k, f)
	if reply == int64(0) {
		return false
	}
	return true
}
