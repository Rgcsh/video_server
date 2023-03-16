//
// All rights reserved
// create time '2022/12/6 15:33'
//
// Usage:
// 一个协程接收UDP的视屏帧 数据
// 一个协程 处理视频帧数据,通过opencv判断是否需要保存,及发送数据到mqtt通知设备休眠

package video_handler

import (
	"fmt"
	"gocv.io/x/gocv"
	"image"
	"image/color"
	"time"
	"video_server/pkg/glog"
	"video_server/pkg/gtime"
	"video_server/pkg/mqtt_service"
	"video_server/pkg/udp_service"
	"video_server/pkg/utils"
)

const MqttTopic string = "camera_frq"
const TimeCycle time.Duration = 60
const NoDiffTimesCount int = 3

//	 @Description: 处理图片的协程,功能如下:
//			用opencv进行动态变化比对,将变化的视频帧存入视频文件;
//			根据某时间段内 视频帧有无变化,动态变更esp32-cam发送视频帧频率
//	 @param ImgChannel:从摄像头的UDP客户端接收到的待处理的图片切片channel,容量为100,超过处理不了就丢弃
//	 @param MqttMessagesChannel:将消息发送到mqtt服务器的channel,容量为2
func Convert2IMG(ImgChannel <-chan []byte, MqttMessagesChannel chan<- []string) {
	redColor := color.RGBA{255, 0, 0, 0}
	es := gocv.GetStructuringElement(gocv.MorphEllipse, image.Point{9, 4})

	//摄像头发送图片频率 0:高 1:低
	sendImgFrequency := utils.SendImgFrequencyHigh
	// 3个固定时间段内 图片都相同,表示环境无变化,通过mqtt协议通知cam降低发送帧率(节约资源[电力,带宽]), 如果 后续 检测到图片变化,则再次通知其恢复频率;
	noDiffTimes := 0
	//一直死循环
	for {
		var writer *gocv.VideoWriter
		//定义一个运行的时间周期
		timeCycle := gtime.TimeCycle(TimeCycle)
		stop := false
		//基准图片,作用为 检测视频帧是否变化的基准图片
		basicImg := gocv.NewMat()
		//存入到磁盘的视频文件writer
		fileName := fmt.Sprintf(utils.VideoFileFmt, gtime.GetCurrentTime())
		//视频帧之间是否不同(检测监控环境是否有变化)
		isDiffImg := false

		//此for循环保证timeCycle的channel一直运行
		for {
			select {
			case _ = <-timeCycle:
				glog.Log.Debug("执行时间周期耗尽")
				stop = true
			default:
				//开始处理 视频帧
				tempImg, ok := <-ImgChannel
				if !ok {
					glog.Log.Info("ImgChannel关闭,停止循环")
					break
				}

				//从数组中读取 前n个有效长度的数据,也就是 一张图片不会超过10万个byte,而n的值就是 读取图片的byte个数
				img, err := gocv.IMDecode(tempImg, gocv.IMReadColor)
				if err != nil {
					glog.Log.Error(fmt.Sprintf("解析图片数据失败,err:%v", err))
					break
				}
				//当basicImg没有赋值时,对其进行赋值
				if basicImg.Empty() {
					basicImg = img
					//图片转为黑白色
					gocv.CvtColor(basicImg, &basicImg, gocv.ColorBGRToGray)
					// 高斯模糊,image.Point的值 越大,则越模糊,注意 值除以2余数必须为1
					gocv.GaussianBlur(basicImg, &basicImg, image.Point{21, 21}, 0, 0, gocv.BorderDefault)
					writer, err = gocv.VideoWriterFile(fileName, "MJPG", 10, basicImg.Cols(), basicImg.Rows(), true)
					if err != nil {
						fmt.Printf("error opening video writer device: %v\n", fileName)
						return
					}
					continue
				}
				//进行图片处理
				//每隔固定时间 获取基准帧数据
				//将现有图片进行比对,不同就存为mp4
				//若 3个固定时间段内 图片都相同,表示环境无变化,通过mqtt协议通知cam降低发送帧率(节约资源[电力,带宽]), 如果 后续 检测到图片变化,则再次通知其恢复频率;

				//	gocv比对图片
				//图片转为黑白色
				gray_frame := gocv.NewMat()
				gocv.CvtColor(img, &gray_frame, gocv.ColorBGRToGray)
				// 高斯模糊,image.Point的值 越大,则越模糊,注意 值除以2余数必须为1
				gocv.GaussianBlur(gray_frame, &gray_frame, image.Point{21, 21}, 0, 0, gocv.BorderDefault)

				//对比2个图片不同点
				diffImg := gocv.NewMat()
				gocv.AbsDiff(basicImg, gray_frame, &diffImg)
				gocv.Threshold(diffImg, &diffImg, 25, 255, gocv.ThresholdBinary)
				gocv.Dilate(diffImg, &diffImg, es)
				cnts := gocv.FindContours(diffImg, gocv.RetrievalExternal, gocv.ChainApproxSimple)

				for i := 0; i < cnts.Size(); i++ {
					timeCycle := cnts.At(i)
					if gocv.ContourArea(timeCycle) < 1500 {
						continue
					}
					isDiffImg = true
					gocv.Rectangle(&img, gocv.BoundingRect(timeCycle), redColor, 1)
				}
				//跳过
				if !isDiffImg {
					glog.Log.Info("没有不同点,跳过此视屏帧")
					continue
				}
				//检测到不同点
				if isDiffImg && sendImgFrequency != utils.SendImgFrequencyHigh {
					glog.Log.Debug("检测到不同点,恢复发送图片频率")
					sendImgFrequency = utils.SendImgFrequencyHigh
					SendMqttMessageToChannel(MqttMessagesChannel, "0")
				}
				gocv.PutText(&img, gtime.GetCurrentTime(), image.Point{10, 20}, gocv.FontHersheyComplexSmall, 1, redColor, 1)
				//存储到mp4中
				err = writer.Write(img)
				if err != nil {
					panic(fmt.Sprintf("视频写入磁盘失败 %v", err))
				}
				glog.Log.Debug("图片解析成功,开始展示")
			}
			//时间周期已到,收尾工作
			if stop {
				glog.Log.Debug("停止执行任务,时间耗尽")
				//视频writer收尾,保证视频文件保存正常
				err := writer.Close()
				if err != nil {
					glog.Log.Error("关闭视频失败....")
				}
				glog.Log.Debug("已经关闭视频")
				break
			}
		}
		if !isDiffImg {
			noDiffTimes += 1
		}

		//降低摄像头发送频率
		if noDiffTimes > NoDiffTimesCount {
			if sendImgFrequency != utils.SendImgFrequencyLow {
				glog.Log.Debug("检测到环境最近无变化,调低摄像头发送频率")
				sendImgFrequency = utils.SendImgFrequencyLow
				SendMqttMessageToChannel(MqttMessagesChannel, "1")
			}
			noDiffTimes = 0
		}
	}
}

// @Description:mqtt消息发送到channel中
// @param MqttMessagesChannel:
// @param message:
func SendMqttMessageToChannel(MqttMessagesChannel chan<- []string, message string) {
	select {
	case MqttMessagesChannel <- []string{MqttTopic, message}:
		glog.Log.Info("发送消息到mqtt成功")
	default:
		glog.Log.Info("mqtt channel缓存已满,丢弃此消息")
	}
}

// @Description: 处理视频帧
func VideoHandler() {
	//从摄像头的UDP客户端接收到的待处理的图片切片channel,容量为100,超过处理不了就丢弃
	ImgReceiveChannel := make(chan []byte, 100)
	//将图片数据转发给用户客户端进行直播的 图片切片channel,容量为20,超过处理不了就丢弃
	ImgSendChannel := make(chan []byte, 20)
	// 将消息发送到mqtt服务器的channel,容量为2
	MqttMessagesChannel := make(chan []string, 2)

	//udp服务接收及发送数据
	go udp_service.UDPReceive(ImgReceiveChannel, ImgSendChannel)
	go udp_service.UDPSend(ImgSendChannel)

	//发布消息到mqtt协议服务的broker(emqx)
	go mqtt_service.MqttPublish(MqttMessagesChannel)

	//处理视频帧数据
	Convert2IMG(ImgReceiveChannel, MqttMessagesChannel)
}
