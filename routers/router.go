// Copyright 2019 Wu Dong
// All rights reserved
//
// @Author: 'Wu Dong <wudong@eastwu.cn>'
// @Time: '2021/10/8 5:09 下午'

package routers

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"video_server/conf"
	"video_server/controllers/video"
	"video_server/crontabs"
	"video_server/middleware"
)

//
// InitRouter
//  @Description: 注册 路由/定时任务
//  @return *gin.Engine:
//
func InitRouter() *gin.Engine {
	r := gin.New()
	// 是否启动 性能分析功能 true:启动
	if conf.Conf.App.EnablePProf {
		pprof.Register(r)
	}

	// 路由注册
	apiVideo := r.Group("/video", middleware.CostTime(), gin.RecoveryWithWriter(nil, middleware.PanicRecovery))
	apiVideo.GET("/download", video_action.DownloadVideoHandler)
	apiVideo.POST("/query", video_action.QueryVideoHandler)
	//定时任务注册
	_, _ = crontabs.Crontab.AddFunc("* * 1 * * *", crontabs.Job1)
	return r
}
