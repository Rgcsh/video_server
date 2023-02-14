// 
// All rights reserved
//
// @Author: 'rgc'

package app

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"reflect"
	"time"
	"video_server/pkg/api_error"
	"video_server/pkg/e"
	"video_server/pkg/glog"
	"video_server/pkg/validator_rewrite"
)

//
//  Response
//  @Description: 响应结构体
//
type Response struct {
	Timestamp int64       `json:"timestamp"`
	RespCode  int         `json:"respCode"`
	RespMsg   string      `json:"respMsg"`
	Result    interface{} `json:"result"`
}

//
// NewResponse
//  @Description: 新建一个响应体对象
//  @param code:
//  @param message:
//  @param body:
//  @return res:
//
func NewResponse(code int, message string, body interface{}) (res *Response) {
	return &Response{
		Timestamp: time.Now().Unix(),
		RespCode:  code,
		RespMsg:   message,
		Result:    body,
	}
}

type Gin struct {
	C      *gin.Context
	Params interface{}
}

//
//  @Description: 新建Gin对象,且将入参绑定到Params参数中
//  @param c:
//  @param params:
//  @return *Gin:
//
func NewGin(c *gin.Context, params interface{}) *Gin {
	ins := new(Gin)
	ins.C = c
	ins.Params = params

	if ins.Params != nil {
		err := ins.Request(ins.Params)
		if err != nil {
			panic(api_error.New(500))
		}
	}
	return ins
}

//
//  @Description: 读取请求入参,并绑定到结构体中
//  @receiver g
//  @param obj:
//  @return error:
//
func (g *Gin) Request(obj interface{}) error {
	//将入参绑定到结构体中
	err := g.C.ShouldBindJSON(&obj)
	if err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			errMessage, _ := json.Marshal(errs.Translate(validator_rewrite.Trans))
			panic(api_error.New(504, string(errMessage)))
		}
		panic(api_error.New(504))
	}
	//结构体进行序列化为 byte切片
	objByte, err := json.Marshal(obj)
	if err != nil {
		g.Response(590, err.Error(), make(map[string]string))
		return err
	}
	//打印入参日志
	glog.Log.Info(fmt.Sprintf("request url: %v ; params: %v", g.C.Request.URL, string(objByte)))
	return nil
}

//
//  @Description: 设置接口响应
//  @receiver g
//  @param code:
//  @param message:
//  @param body:
//
func (g *Gin) Response(code int, message string, body interface{}) {
	res := NewResponse(code, message, body)
	g.C.JSON(200, res)
	resBytes, err := json.Marshal(res)
	if err != nil {
		glog.Log.Error(fmt.Sprintf("%v", err))
		return
	}
	glog.Log.Info(fmt.Sprintf("request url: %v ; response: %v ;", g.C.Request.URL, string(resBytes)))
}

//
// Success
//  @Description:成功响应
//  @receiver g
//  @param body: 实现 入参非必传
//
func (g *Gin) Success(body ...interface{}) {
	if body == nil || reflect.ValueOf(body).IsNil() {
		g.Response(200, e.GetMessage(200), make(map[string]string))
	} else {
		g.Response(200, e.GetMessage(200), body[0])
	}
}

//
//  @Description: 失败响应
//  @receiver g
//  @param code:
//  @param message:
//
func (g *Gin) Fail(code int, message string) {
	if message == "" {
		g.Response(code, e.GetMessage(code), make(map[string]string))
	} else {
		g.Response(code, message, make(map[string]string))
	}
}
