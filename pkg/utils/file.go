// (C) Guangcai Ren <rgc@bvrft.com>
// All rights reserved
// create time '2022/12/9 21:48'
//
// Usage:
//

package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"video_server/pkg/api_error"
)

//
//  @Description: 获取文件夹中的所有文件
//  @param dir:
//  @return *[]string:
//
func FileLS(dir string) *[]string {
	var FileSlice []string
	pwd, err := os.Getwd()
	if err != nil {
		panic(api_error.New(504, "获取当前路径失败 err:%v", err))
	}
	fmt.Println(pwd)
	fileInfos, err := ioutil.ReadDir(filepath.Join(pwd, dir))

	if err != nil {
		panic(api_error.New(504, fmt.Sprintf("获取当前路径失败 err:%v", err)))
	}

	for i := range fileInfos {
		FileSlice = append(FileSlice, fileInfos[i].Name())
	}
	return &FileSlice
}

//
//  @Description: 过滤 掉 非avi结尾的文件
//  @param FileNames:
//  @return *[]string:
//
func FileFilter(FileNames *[]string) *[]string {
	var FilterFileNames []string
	usefulFileSuffix := "avi"
	for i := range *FileNames {
		if strings.HasSuffix((*FileNames)[i], usefulFileSuffix) {
			FilterFileNames = append(FilterFileNames, (*FileNames)[i])
		}
	}
	return &FilterFileNames
}

//
//  @Description: 对文件名按照时间 从过去到现在排序
//  @param FileNames:
//  @return *[]string:
//
func FileSort(FileNames *[]string) *[]string {
	sort.Strings(*FileNames)
	return FileNames
}

//
//  @Description: 视频文件处理  从读取到过滤到排序
//  @param dir:
//  @return FileNames:
//
func VideoFileHandler(dir string) (FileNames *[]string) {
	return FileSort(FileFilter(FileLS(dir)))
}
