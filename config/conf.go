// Package config
/**
@author: xdl2003
@desc:
@date: 2025/6/5
**/
package config

import (
	"github.com/spf13/viper"
	"go-manus/go-manus/log"
	"go.uber.org/zap"
)

var (
	AllConfig *Settings
)

// Settings
//
//	@Description: 总配置结构体
type Settings struct {

	//  PrimaryConfig
	//  @Description: 计划模型配置
	PrimaryConfig *ModelConfig `yaml:"primary_config"`

	//  ExecutorConfig
	//  @Description: 执行模型配置
	ExecutorConfig *ModelConfig `yaml:"executor_config"`
}

// ModelConfig
//
//	@Description: 模型配置结构体
type ModelConfig struct {
	ModelSource string `yaml:"model_source"`
	ModelName   string `yaml:"model_name"`
	ApiKey      string `yaml:"api_key"`
}

// InitConfig
//
//	@Description: 初始化配置, 通过viper在config.yaml中读取配置
func InitConfig() {
	v := viper.New()
	v.SetConfigName("config") // 设置配置文件名 (不带后缀)
	v.SetConfigType("yaml")
	v.AddConfigPath(".") // 第一个搜索路径
	v.AddConfigPath("./config")
	v.AddConfigPath("./../config")
	v.AddConfigPath("./../../config")

	err := v.ReadInConfig() // 读取配置数据
	if err != nil {
		log.Logger.Error("fail to read config, err=%v", zap.Error(err))
		panic(err)
	}
	e := v.Unmarshal(&AllConfig) // 将配置信息绑定到结构体上
	if e != nil {
		log.Logger.Info("fail to unmarshal, e=%v", zap.Any("error", e))
	}
}
