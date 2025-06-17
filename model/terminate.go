// Package model
/**
@author: xdl2003
@desc:
@date: 2025/6/17
**/
package model

type TerminateInput struct {
	Status string `json:"status" description:"The finish status of the interaction." enum:"成功,失败"`
}
