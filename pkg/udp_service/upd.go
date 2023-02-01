// (C) Guangcai Ren <rgc@bvrft.com>
// All rights reserved
// create time '2022/12/9 16:56'
//
// Usage:
//

package udp_service

import (
	"fmt"
	"net"
)

var UDPListen *net.UDPConn

func UPDSetUp() {
	var err error
	UDPListen, err = net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 9090,
	})
	if err != nil {
		panic(fmt.Sprintf("监听UDP服务失败,请修改,err:", err))
	}
	fmt.Println("启动UDP服务成功,端口为:9090......")
	return
}

func UPDClose() {
	err := UDPListen.Close()
	if err != nil {
		panic(fmt.Sprintf("关闭UDPListen失败,err:%v", err))
	}
}

var SendToClient bool
var SendToClientAddr *net.UDPAddr

//
//  @Description: 接收UDP数据,并发送给channel,当channel缓存满时,丢弃数据
//  @param ImgChannel:
//
func UDPReceive(ImgChannel chan<- []byte, ImgSendChannel chan<- []byte) {
	SendToClient = false
	for {
		//定义一个长度为10万的切片
		sliceData := make([]byte, 65500)
		//将读取的数据存到数组中
		n, addr, err := UDPListen.ReadFromUDP(sliceData)
		fmt.Println(addr)
		if err != nil {
			fmt.Println("读取UDP数据失败,err:", err)
			continue
		}
		if n == 0 || n > 65500 {
			continue
		}
		//截取前10个字符串
		messageType := string(sliceData[:10])
		imgData := sliceData[10:]
		if messageType == "clientPlay" {
			fmt.Println("接收到客户端的直播请求")
			SendToClient = true
			SendToClientAddr = addr
		} else if messageType == "clientStop" {
			SendToClient = false
		} else if messageType == "cameraSend" {
			fmt.Println("接收到图片数据")
			select {
			case ImgChannel <- imgData:
				fmt.Println("发送数据到图片处理channel")
			default:
				fmt.Println("发送数据到图片处理channel阻塞,不发送")
			}
		}

		if SendToClient && messageType == "cameraSend" {
			select {
			case ImgSendChannel <- imgData:
				fmt.Println("发送数据到转发channel")
			default:
				fmt.Println("发送数据到转发channel阻塞,不发送")
			}
		}
	}
}

//
//  @Description: 发送UDP数据给客户端,并发送给channel,当channel缓存满时,丢弃数据
//  @param ImgChannel:
//
func UDPSend(ImgSendChannel <-chan []byte) {
	for {
		Img, ok := <-ImgSendChannel
		if !ok {
			fmt.Println("发送数据到图片处理channel关闭")
			break
		}
		fmt.Println("发送直播数据到客户端")
		_, err := UDPListen.WriteToUDP(Img, SendToClientAddr)
		if err != nil {
			panic(fmt.Sprintf("转发UDP数据失败 err:%v", err))
		}
	}
}
