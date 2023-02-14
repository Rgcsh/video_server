// 
// All rights reserved
// create time '2022/12/9 21:34'
//
// Usage:
//

package utils

//视频存储文件夹路径
const VideoFileDir string = "video_file"
const VideoFileFmt string = "./video_file/%v.avi"

//摄像头发送图片频率 0:高 1:低
const (
	SendImgFrequencyHigh int = iota
	SendImgFrequencyLow
)
