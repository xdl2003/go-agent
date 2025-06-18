// Package model
/**
@author: xdl2003
@desc:
@date: 2025/6/6
**/
package model

import (
	"context"
	"github.com/mark3labs/mcp-go/mcp"
	mcp2 "go-manus/go-manus/mcp"
)

func MCPTool2Tool(mcpTool *mcp.Tool) *Tool {
	tool := Tool{
		Type: ToolTypeFunction,
		Function: FunctionDefinition{
			Name:        mcpTool.Name,
			Description: mcpTool.Description,
			Parameters:  mcpTool.InputSchema,
		},
	}
	return &tool
}

func GetToolList() ([]*Tool, error) {
	tools, err := mcp2.McpClient.ListTools(context.Background(), mcp.ListToolsRequest{})
	if err != nil {
		return nil, err
	}
	var result []*Tool
	for _, tool := range tools.Tools {
		result = append(result, MCPTool2Tool(&tool))
	}
	return result, nil
}

type Tool struct {
	Type     ToolType
	Function FunctionDefinition
}

func (t *Tool) GetTool() *Tool {
	return t
}

func (t *Tool) Execute(input string, method string) (string, error) {
	return mcp2.Execute(input, method)
}

type FunctionDefinition struct {
	Name        string
	Description string
	Parameters  interface{}
}

type BaseTool interface {
	GetTool() *Tool
	Execute(input string, method string) (string, error)
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
