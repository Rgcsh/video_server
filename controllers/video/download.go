package video_action

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"path/filepath"
	"strings"
	"video_server/pkg/app"
	"video_server/pkg/glog"
	"video_server/pkg/utils"
)

type DownloadFields struct {
	FileName string `json:"file_name" binding:"required"`
}

//
//  @Description: 下载视频文件
// /video/download
//  @param c:
//
func DownloadVideoHandler(c *gin.Context) {
	// 生成上下文环境及入参赋值
	fileName := c.Query("FileName")
	ctx := app.NewGin(c, nil)
	
	// 校验文件后缀 及 . 的个数 判断文件名是否合法
	fmt.Println(strings.Split(fileName, "."), len(strings.Split(fileName, ".")), fileName)
	if !strings.HasSuffix(fileName, ".avi") || len(strings.Split(fileName, ".")) > 2 {
		utils.StatusCodeHandler(ctx, 500, "文件名错误")
		return
	}
	// 判断文件是否存在
	// 拼接文件绝对路径
	pwd, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("获取当前路径失败 err:%v", err))
	}
	fileAbsPath := filepath.Join(pwd, utils.VideoFileDir, fileName)
	if _, err := os.Stat(fileAbsPath); os.IsNotExist(err) {
		utils.StatusCodeHandler(ctx, 500, "文件不存在")
		return
	}
	glog.Log.Info("文件检测结束")
	// 返回文件
	c.File(fileAbsPath)
}
