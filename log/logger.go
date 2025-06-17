// Package log
/**
@author: xdl2003
@desc: 日志模块
@date: 2025/6/5
**/
package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var (
	Logger  *zap.Logger
	LogFile *os.File
)

func InitLogger() {
	logFile, err := os.OpenFile("./app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	LogFile = logFile
	if err != nil {
		panic(err)
	}
	// 配置编码器
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// 创建核心（输出到文件，JSON格式，INFO级别）
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(logFile),
		zap.InfoLevel,
	)
	Logger = zap.New(core)
}

func CloseLogger() {
	err := LogFile.Close()
	if err != nil {
		panic(err)
	}
	err = Logger.Sync()
	if err != nil {
		return
	}
}
