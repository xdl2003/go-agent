// Package tool
/**
@author: xdl2003
@desc:
@date: 2025/6/17
**/
package tool

import (
	"go-manus/go-manus/model"
)

func GetAvailableTools() map[string]model.BaseTool {
	result := map[string]model.BaseTool{
		"terminate": NewTerminateTool(),
		// "plan":      NewPlanTool(),
	}
	tools, _ := model.GetToolList()
	for _, tool := range tools {
		result[tool.Function.Name] = tool
	}
	return result
}
