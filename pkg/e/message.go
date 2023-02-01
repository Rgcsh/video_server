// Copyright 2019 Wu Dong
// All rights reserved
//
// @Author: 'Wu Dong <wudong@eastwu.cn>'
// @Time: '2021/10/9 1:13 下午'

package e

var Message = map[int]string {
	200: "Success",
	500: "Server Internal error",
	504: "入参解析失败",
}

func GetMessage(code int) string {
	message, ok := Message[code]
	if ok {
		return message
	}
	return Message[500]
}
