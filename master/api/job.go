package api

import (
	"github.com/crazy-me/os_scheduler/common/logger"
	"github.com/crazy-me/os_scheduler/master/entity"
	"github.com/crazy-me/os_scheduler/master/etcd"
	"github.com/crazy-me/os_scheduler/master/resut"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

// SaveJob 保存Job任务
func SaveJob(c *gin.Context) {
	var (
		err     error
		job     entity.Job
		saveJob *entity.Job
	)
	if err = c.ShouldBindJSON(&job); err != nil {
		logger.L.Error("api-SaveJob err:", zap.Any("c.ShouldBindJSON", err))
		c.JSON(http.StatusOK, resut.FAIL())
		return
	}
	// 保存Job
	if saveJob, err = etcd.Client.SaveJob(&job); err != nil {
		logger.L.Error("api-SaveJob err:", zap.Any("etcd.Client.SaveJob", err))
		c.JSON(http.StatusOK, resut.FAIL())
		return
	}
	c.JSON(http.StatusOK, resut.DATA(saveJob))
}

// DelJob 删除Job任务
func DelJob(c *gin.Context) {
	var (
		err    error
		job    entity.Job
		delJob *entity.Job
	)
	if err = c.ShouldBindJSON(&job); err != nil {
		logger.L.Error("api-DelJob err:", zap.Any("c.ShouldBindJSON", err))
		c.JSON(http.StatusOK, resut.FAIL())
		return
	}

	// 删除Job
	if delJob, err = etcd.Client.DeleteJob(&job); err != nil {
		logger.L.Error("api-DelJob err:", zap.Any("etcd.Client.DeleteJob", err))
		c.JSON(http.StatusOK, resut.FAIL())
		return
	}
	c.JSON(http.StatusOK, resut.DATA(delJob))
}

// ListJob 获取所有Job
func ListJob(c *gin.Context) {
	var (
		err     error
		listJob []*entity.Job
	)

	if listJob, err = etcd.Client.ListJob(); err != nil {
		logger.L.Error("api-ListJob err:", zap.Any("etcd.Client.ListJob", err))
		c.JSON(http.StatusOK, resut.FAIL())
		return
	}
	c.JSON(http.StatusOK, resut.DATA(listJob))
}

// KillJob 强杀Job
func KillJob(c *gin.Context) {
	var job entity.Job
	if err := c.ShouldBindJSON(&job); err != nil {
		logger.L.Error("api-KillJob err:", zap.Any("c.ShouldBindJSON", err))
		c.JSON(http.StatusOK, resut.FAIL())
		return
	}
	if err := etcd.Client.KillJob(&job); err != nil {
		logger.L.Error("api-KillJob err:", zap.Any("etcd.Client.KillJob", err))
		c.JSON(http.StatusOK, resut.FAIL())
		return
	}
	c.JSON(http.StatusOK, resut.SUCCESS())
}
