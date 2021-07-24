package main

import (
	"flag"
	"fmt"
	"github.com/crazy-me/os_scheduler/common/logger"
	"github.com/crazy-me/os_scheduler/master/conf"
	"github.com/crazy-me/os_scheduler/master/etcd"
	"github.com/crazy-me/os_scheduler/master/router"
	"log"
	"os"
	"runtime"
)

var configFile string

func main() {
	initEnv()
	initArgs()
	initLoad()
	if err := etcd.InitEtcd(); err != nil {
		fmt.Println("etcd clientv3.New err:", err)
		os.Exit(-1)
	}

	router.RunServer()
}

// initArgs 命令行参数
func initArgs() {
	flag.StringVar(&configFile, "c", "etc/scheduler.yaml", "configuration")
	flag.Parse()
}

// initEnv 系统环境
func initEnv() {
	cpuNum := runtime.NumCPU()
	fmt.Printf("System Cpu Kernel number:%d\n", cpuNum)
	runtime.GOMAXPROCS(cpuNum)
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
