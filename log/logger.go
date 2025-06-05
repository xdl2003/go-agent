// Package log
/**
@author: xudongliu.666@bytedance.com
@desc:
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
