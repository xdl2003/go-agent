// Package model
/**
@author: xdl2003
@desc:
@date: 2025/6/6
**/
package model

type Tool struct {
	Type     ToolType
	Function FunctionDefinition
}

type FunctionDefinition struct {
	Name        string
	Description string
	Parameters  interface{}
}

type BaseTool interface {
	GetTool() *Tool
	Execute(input string) (string, error)
}

type ToolType string

const (
	ToolTypeFunction ToolType = "function"
)

type ToolCall struct {
	ID       string       `json:"id"`
	Type     ToolType     `json:"type"`
	Function FunctionCall `json:"function"`
	Index    *int         `json:"index,omitempty"`
}

type FunctionCall struct {
	Name      string `json:"name,omitempty"`
	Arguments string `json:"arguments,omitempty"`
}
