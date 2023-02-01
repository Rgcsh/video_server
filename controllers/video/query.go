package video_action

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
	"video_server/pkg/api_error"
	"video_server/pkg/app"
	"video_server/pkg/gtime"
	"video_server/pkg/utils"
)

type QueryFields struct {
	StartTime string `json:"start_time" binding:"required"`
	EndTime   string `json:"end_time" binding:"required"`
}

//
//  @Description: 查询 关注数据 或检查是否添加关注
//  四个参数都传 可以检查 对应数据是否添加关注,返回的切片长度为0 表示未关注
//  此接口可以灵活的查询不同范围的数据
//  @param c:
//  /video/query
//
func QueryVideoHandler(c *gin.Context) {
	// 用户入参
	var params QueryFields
	// 生成上下文环境及入参赋值
	ctx := app.NewGin(c, &params)

	startTime := params.StartTime
	endTime := params.EndTime

	//获取文件列表
	files := utils.VideoFileHandler(utils.VideoFileDir)

	//读取 > StartTime时间的后PageSize条数据
	resultFiles := GetGtStartTimeFiles(files, startTime, endTime)
	// 结果数据结构
	ctx.Success(resultFiles)
}

//
//  @Description: 获取 大于 开始时间的文件
//  @param FileNames:
//  @param StartTime:
//  @param PageSize:
//  @return *[]string:
//
func GetGtStartTimeFiles(FileNames *[]string, StartTime, EndTime string) *[]string {
	var FilterFiles []string
	//类型转换
	StartTimeObj, err := gtime.StringToTime(StartTime, gtime.DateTimeFormat)
	if err != nil {
		panic(fmt.Sprintf("时间字符串类型转为Time类型失败,StartTime:%v,err:%v", StartTime, err))
	}
	EndTimeObj, err := gtime.StringToTime(EndTime, gtime.DateTimeFormat)
	if err != nil {
		panic(fmt.Sprintf("时间字符串类型转为Time类型失败,StartTime:%v,err:%v", StartTime, err))
	}
	//校验入参StartTime必须<EndTime
	if EndTimeObj.Before(*StartTimeObj) {
		err := api_error.New(504, "开始日期必须小于结束日期")
		panic(err)
	}
	//循环遍历视频文件名称
	for i := range *FileNames {
		fileName := (*FileNames)[i]
		//	切分获取文件前缀(字符串格式的时间)
		fileTime := strings.Split(fileName, ".avi")[0]
		fileTimeObj, err := gtime.StringToTime(fileTime, gtime.DateTimeFormat)
		if err != nil {
			panic(fmt.Sprintf("时间字符串类型转为Time类型失败,fileTime:%v,err:%v", fileTime, err))
		}
		//取 StartTime之后 EndTime之前 的数据
		if StartTimeObj.Before(*fileTimeObj) && fileTimeObj.Before(*EndTimeObj) {
			FilterFiles = append(FilterFiles, fileName)
		}
	}
	return &FilterFiles
}
