// Package util
/**
@author: xdl2003
@desc:
@date: 2025/6/16
**/
package util

import (
	"encoding/json"
	"go-manus/go-manus/log"
	"go.uber.org/zap"
	"runtime/debug"
)

func MustJson(data interface{}) (j string) {
	if data == nil {
		return "{}"
	}
	defer func() {
		if err := recover(); err != nil {
			log.Logger.Error("MustJson error: %s, stack: %s", zap.Any("err", err), zap.Any("stack", string(debug.Stack())))
		}
	}()
	r, e := json.Marshal(data)
	if e != nil {
		panic(e)
	}
	return string(r)
}
