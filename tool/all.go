// Package tool
/**
@author: xdl2003
@desc:
@date: 2025/6/17
**/
package tool

import "go-manus/go-manus/model"

func GetAvailableTools() map[string]model.BaseTool {
	return map[string]model.BaseTool{
		"terminate": NewTerminateTool(),
		"plan":      NewPlanTool(),
	}
}
