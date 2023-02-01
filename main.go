// Copyright 2019 Wu Dong
// All rights reserved
//
// @Author: 'Wu Dong <wudong@eastwu.cn>'
// @Time: '2021/9/30 5:03 下午'

package main

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"video_server/conf"
	"video_server/crontabs"
	"video_server/pkg/glog"
	"video_server/pkg/validator_rewrite"
	"video_server/routers"
	"video_server/video_handler"
)

func init() {
	conf.SetUp()
	crontabs.SetUp()
	glog.SetUp()
	validator_rewrite.SetUp()
}

func main() {
	go video_handler.VideoUDPServiceSetUp()
	gin.SetMode(gin.ReleaseMode)
	// 设置未声明的参数无法传参
	gin.EnableJsonDecoderDisallowUnknownFields()
	r := routers.InitRouter()

	glog.Log.Info("http server will start at", zap.Int("port", 7069))
	if err := r.Run(":7069"); err != nil {
		glog.Log.Error("http server start failed", zap.Error(err))
		return
	}
}
