// Package log
/**
@author: xdl2003
@desc: 日志模块
@date: 2025/6/5
**/
package log

import "go.uber.org/zap"

var (
	Logger *zap.Logger
)

func InitLogger() {
	Logger = zap.NewExample()
}
