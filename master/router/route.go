package router

import (
	"fmt"
	"github.com/crazy-me/os_scheduler/master/api"
	"github.com/crazy-me/os_scheduler/master/conf"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func RunServer() {
	var route = gin.Default()
	if conf.C.System.Env == "app" {
		gin.SetMode(gin.ReleaseMode)
	}
	// 加载静态文件及模版
	route.StaticFS("/static", http.Dir("./resource"))
	route.LoadHTMLGlob("resource/view/*")
	route.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Job Manage",
		})
	})
	route.POST("/job/save", api.SaveJob)
	route.POST("/job/kill", api.KillJob)
	route.POST("/job/del", api.DelJob)
	route.POST("/job/list", api.ListJob)
	address := ":" + string(conf.C.System.Addr)
	serve := endless.NewServer(address, route)
	serve.ReadHeaderTimeout = 10 * time.Millisecond
	serve.WriteTimeout = 10 * time.Second
	serve.MaxHeaderBytes = 1 << 20
	fmt.Println("server run success:http://0.0.0.0" + address)
	log.Println(serve.ListenAndServe().Error())
}
