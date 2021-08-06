package main

import (
	"flag"
	"github.com/crazy-me/os_scheduler/work/conf"
	"github.com/crazy-me/os_scheduler/work/data_source/etcd"
	"github.com/crazy-me/os_scheduler/work/data_source/redis"
	"github.com/crazy-me/os_scheduler/work/logger"
	"github.com/crazy-me/os_scheduler/work/logic"
	"log"
	"os"
)

var configFile string

func main() {
	initArgs()
	initLoad()
	redis.InitRedis()
	// 调度器
	if err := logic.InitSchedule(); err != nil {
		log.Println("logic.InitSchedule err:", err)
		os.Exit(-1)
	}

	// 执行器
	if err := logic.InitExecutor(); err != nil {
		log.Println("logic.InitExecutor err:", err)
		os.Exit(-1)
	}

	// Etcd
	if err := etcd.InitEtcd(); err != nil {
		log.Println("etcd.InitEtcd err:", err)
		os.Exit(-1)
	}

	// 监听任务事件
	if err := logic.WatchJobs(); err != nil {
		log.Println("logic.WatchJobs err:", err)
		os.Exit(-1)
	}
	select {}
}

func initArgs() {
	flag.StringVar(&configFile, "c", "etc/work.yaml", "configuration")
	flag.Parse()
}

func initLoad() {
	// 加载配置文件
	err := conf.InitConf(configFile)
	if err != nil {
		log.Println("load configuration file err:", err)
		os.Exit(-1)
	}

	// 初始化日志
	err = logger.InitLogger()
	if err != nil {
		log.Println("logger init err:", err)
		os.Exit(-1)
	}
}
