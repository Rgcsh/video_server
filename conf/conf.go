// Copyright 2019 Wu Dong
// All rights reserved
//
// @Author: 'Wu Dong <wudong@eastwu.cn>'
// @Time: '2021/10/8 4:30 下午'

package conf

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

type App struct {
	//是否启动 性能分析功能 true:启动
	EnablePProf bool
}

type Config struct {
	LogConf LogConf `yaml:"LogConf"`
	App     App     `yaml:"App"`
}

type LogConf struct {
	FilePath string `yaml:"FilePath"`
}

var Conf = &Config{}

// SetUp 读取配置信息
func SetUp() {
	path := os.Getenv("LOC_CFG")
	if path == "" {
		panic("Invalid config path to load please use `LOC_CFG` set to os environment!")
	}

	if f, err := os.Open(path); err != nil {
		panic(fmt.Sprintf("Load config from path %s failed: %s", path, err.Error()))
	} else {
		_ = yaml.NewDecoder(f).Decode(Conf)
	}
}
