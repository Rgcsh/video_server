//
// All rights reserved
// create time '2022/12/9 16:56'
//
// Usage:
//

package udp_service

import (
	"fmt"
	"net"
	"video_server/conf"
	"video_server/pkg/glog"
)

var UDPListen *net.UDPConn

//
//  @Description: UDP服务端初始化
//
func UDPSetUp() {
	var err error
	UDPListen, err = net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: conf.Conf.UDPConf.Port,
	})
	if err != nil {
		panic(fmt.Sprintf("监听UDP服务失败,请修改,err:", err))
	}
	glog.Log.Info(fmt.Sprintf("启动UDP服务成功,端口为:%v", conf.Conf.UDPConf.Port))
	return
}

//
//  @Description: 关闭UDP服务
//
func UDPClose() {
	err := UDPListen.Close()
	if err != nil {
		panic(fmt.Sprintf("关闭UDPListen失败,err:%v", err))
	}
}

var SendToClient bool
var SendToClientAddr *net.UDPAddr

//
//  @Description: 接收UDP数据,并发送给channel,当channel缓存满时,丢弃数据
//  @param ImgReceiveChannel: 从摄像头的UDP客户端接收到的待处理的图片切片channel,容量为100,超过处理不了就丢弃
//  @param ImgSendChannel: 将图片数据转发给用户客户端进行直播的 图片切片channel,容量为20,超过处理不了就丢弃
//
func UDPReceive(ImgReceiveChannel chan<- []byte, ImgSendChannel chan<- []byte) {
	SendToClient = false
	for {
		//定义一个切片
		sliceData := make([]byte, 65500)
		//将读取的数据存到数组中
		n, addr, err := UDPListen.ReadFromUDP(sliceData)
		if err != nil {
			glog.Log.Error(fmt.Sprintf("读取UDP数据失败,err:%v", err))
			continue
		}
		if n == 0 || n > 65500 {
			continue
		}
		//截取前10个字符串
		messageType := string(sliceData[:10])
		imgData := sliceData[10:]
		if messageType == "clientPlay" {
			glog.Log.Info("接收到客户端的直播请求")
			SendToClient = true
			SendToClientAddr = addr
		} else if messageType == "clientStop" {
			SendToClient = false
		} else if messageType == "cameraSend" {
			glog.Log.Info("接收到图片数据")
			select {
			case ImgReceiveChannel <- imgData:
				glog.Log.Info("发送数据到图片处理channel")
			default:
				glog.Log.Info("发送数据到图片处理channel阻塞,不发送")
			}
		}

		if SendToClient && messageType == "cameraSend" {
			select {
			case ImgSendChannel <- imgData:
				glog.Log.Info("发送数据到转发channel")
			default:
				glog.Log.Info("发送数据到转发channel阻塞,不发送")
			}
		}
	}
}

//
//  @Description: 发送UDP数据给客户端(直播使用),当channel缓存满时,丢弃数据
//  @param ImgSendChannel:
//
func UDPSend(ImgSendChannel <-chan []byte) {
	for {
		Img, ok := <-ImgSendChannel
		if !ok {
			glog.Log.Info("发送数据到图片处理channel关闭")
			break
		}
		glog.Log.Info("发送直播数据到客户端")
		_, err := UDPListen.WriteToUDP(Img, SendToClientAddr)
		if err != nil {
			panic(fmt.Sprintf("转发UDP数据失败 err:%v", err))
		}
	}
}
