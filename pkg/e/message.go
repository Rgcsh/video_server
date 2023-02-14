// 
// All rights reserved
//
// @Author: 'rgc'

package e

var Message = map[int]string{
	200: "Success",
	500: "Server Internal error",
	504: "入参解析失败",
}

//
//  @Description: 根据状态码获取对应的消息
//  @param code:
//  @return string:
//
func GetMessage(code int) string {
	message, ok := Message[code]
	if ok {
		return message
	}
	return Message[500]
}
