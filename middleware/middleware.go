// Package middleware
// @Description: 中间件 可以注册到 全局路由,或者某个路由组 或者单个路由中
//方法如下链接的 注册中间件 部分:https://www.liwenzhou.com/posts/Go/Gin_framework/
package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"video_server/pkg/app"
	"video_server/pkg/e"
	"video_server/pkg/glog"
	"video_server/pkg/utils"
	"time"
)

//
// CostTime
//  @Description: 定义一个 统计接口耗时的中间件
//  @return gin.HandlerFunc:
//
func CostTime() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		//继续执行请求的后续代码
		c.Next()
		// 计算距离当前时间的耗时
		costTime := time.Since(startTime)
		glog.Log.Info(fmt.Sprintf("接口:%v,耗时:%v", c.Request.URL, costTime))
	}
}

//
// PanicRecovery
//  @Description: 此处捕捉到错误后,判断错误是否为 可以处理的错误,如果是就 返回相应的错误码给前端
//  此方法给 CustomRecovery 中间件使用
//  @param c:
//  @param err:
//
func PanicRecovery(c *gin.Context, err interface{}) {
	code, message := utils.ModelErrorHandler(err)
	if message == "" {
		c.AbortWithStatusJSON(200, app.NewResponse(code, e.GetMessage(code), make(map[string]interface{})))
	} else {
		c.AbortWithStatusJSON(200, app.NewResponse(code, message, make(map[string]interface{})))
	}
}
