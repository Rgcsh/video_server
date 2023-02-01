package utils

import (
	"fmt"
	"video_server/pkg/api_error"
	"video_server/pkg/app"
	"video_server/pkg/glog"
	"runtime/debug"
)

//
// ModelErrorHandler
//  @Description: 错误捕获处理中心
//  @param err: 错误内容,最好传 错误对象的指针
//  @return int: 对应的错误码
//
func ModelErrorHandler(err interface{}) (int, string) {
	switch errCase := err.(type) {
	// 匹配 此自定义错误的指针类型
	case *api_error.ApiError:
		//打印错误栈
		glog.Log.Debug(string(debug.Stack()))
		return errCase.ErrCode, errCase.ErrMessage
	case error:
		if errCase.Error() == "record not found" {
			return 503, ""
		}
		//打印错误栈
		glog.Log.Error(fmt.Sprintf("%v", err))
		stackMsg := string(debug.Stack())
		glog.Log.Error(stackMsg)
		return 500, ""
	default:
		return 200, ""
	}
}

//
//  @Description:
//  @param responseCode:
//  @param ctx:
//
func StatusCodeHandler(ctx *app.Gin, responseCode int, message string) {
	if responseCode != 200 {
		ctx.Fail(responseCode, message)
	} else {
		ctx.Success()
	}
}

//
//  @Description:
//  @param err:
//  @param ctx:
//
func HandlerErrorAndResponse(ctx *app.Gin, err error) {
	code, message := ModelErrorHandler(err)
	StatusCodeHandler(ctx, code, message)
}

//
// ErrorCheck
//  @Description: 错误检查,如果捕捉到错误就直接panic,不再执行后续代码
//  @param err:
//
func ErrorCheck(err error) {
	code, message := ModelErrorHandler(err)
	if code != 200 {
		panic(api_error.New(code, message))
	}
}
