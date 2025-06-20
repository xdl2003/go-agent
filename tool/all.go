// Package tool
/**
@author: xdl2003
@desc:
@date: 2025/6/17
**/
package tool

import (
	"go-manus/go-manus/mcp"
	"go-manus/go-manus/model"
)

func GetAvailableTools() map[string]model.BaseTool {
	result := map[string]model.BaseTool{
		"terminate": NewTerminateTool(),
		// "plan":      NewPlanTool(),
		"ask_human": NewAskHumanTool(),
	}
	for _, tool := range mcp.AllTools {
		result[tool.Name] = model.MCPTool2Tool(tool)
	}
	return result
}
