//
// All rights reserved
//
// @Author: 'rgc'

package main

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"video_server/conf"
	"video_server/crontabs"
	"video_server/pkg/glog"
	"video_server/pkg/mqtt_service"
	"video_server/pkg/udp_service"
	"video_server/pkg/validator_rewrite"
	"video_server/routers"
	"video_server/video_handler"
)

//
//  @Description: 初始化函数,在main函数前执行
//
func init() {
	conf.SetUp()
	crontabs.SetUp()
	glog.SetUp()
	validator_rewrite.SetUp()
	//udp服务启动
	udp_service.UDPSetUp()
	//mqtt服务启动
	mqtt_service.MqttSetUp()
}

//
//  @Description: 项目入口
//
func main() {
	defer udp_service.UDPClose()
	go video_handler.VideoHandler()
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
