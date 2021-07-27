package main

import (
	"flag"
	"github.com/crazy-me/os_scheduler/work/conf"
	"github.com/crazy-me/os_scheduler/work/etcd"
	"github.com/crazy-me/os_scheduler/work/executor"
	"github.com/crazy-me/os_scheduler/work/logger"
	"github.com/crazy-me/os_scheduler/work/scheduler"
	"log"
	"os"
	"runtime"
)

var configFile string

func main() {
	initEnv()
	initArgs()
	initLoad()
	// 启动调度程序
	if err := scheduler.InitSchedule(); err != nil {
		log.Println("scheduler start err:", err)
		os.Exit(-1)
	}

	// 启动执行器
	if err := executor.InitExecutor(); err != nil {
		log.Println("executor.InitExecutor err:", err)
		os.Exit(-1)
	}

	if err := etcd.InitEtcd(); err != nil {
		log.Println("etcd client.New err:", err)
		os.Exit(-1)
	}

	select {}
}

// 系统环境
func initEnv() {
	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)
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
