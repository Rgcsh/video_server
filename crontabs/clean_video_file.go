package crontabs

import (
	"fmt"
	"strings"
	"video_server/pkg/gtime"
	"video_server/pkg/utils"
)

func CleanVideoFile() {
	fmt.Println("开始删除7天前的视频文件")
	//  获取7天前的时间
	StartTimeObj := gtime.GetPastDate(-7)
	//获取文件列表
	FileNames := utils.VideoFileHandler(utils.VideoFileDir)
	for _, fileName := range *FileNames {
		//	切分获取文件前缀(字符串格式的时间)
		fileTime := strings.Split(fileName, ".avi")[0]
		fileTimeObj, err := gtime.StringToTime(fileTime, gtime.DateTimeFormat)
		if err != nil {
			continue
		}
		//判断文件是否在7天前
		if fileTimeObj.Before(*StartTimeObj) {
			fmt.Printf("开始删除7天前视频文件:%v \n", fileName)
			utils.DelFile(fmt.Sprintf(utils.VideoFileFmt, fileTime))
		}
	}
}
