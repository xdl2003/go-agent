// Package model
/**
@author: xdl2003
@desc:
@date: 2025/6/16
**/
package model

// AgentState 表示Agent的状态
type AgentState int

const (
	// AgentStateIDLE 表示Agent处于空闲状态
	AgentStateIDLE AgentState = iota
	// AgentStateRUNNING 表示Agent正在运行
	AgentStateRUNNING
	// AgentStateFINISHED 表示Agent已完成执行
	AgentStateFINISHED
	// AgentStateERROR 表示Agent遇到错误
	AgentStateERROR
)

// String 将AgentState转换为字符串
func (s AgentState) String() string {
	switch s {
	case AgentStateIDLE:
		return "IDLE"
	case AgentStateRUNNING:
		return "RUNNING"
	case AgentStateFINISHED:
		return "FINISHED"
	case AgentStateERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}
