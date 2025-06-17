// Package tool
/**
@author: xdl2003
@desc:
@date: 2025/6/17
**/
package tool

import (
	"encoding/json"
	"fmt"
	"github.com/sashabaranov/go-openai/jsonschema"
	"go-manus/go-manus/model"
)

type TerminateTool struct {
	*model.Tool
}

func NewTerminateTool() *TerminateTool {
	t, _ := jsonschema.GenerateSchemaForType(model.TerminateInput{})
	return &TerminateTool{
		Tool: &model.Tool{
			Type: model.ToolTypeFunction,
			Function: model.FunctionDefinition{
				Name:        "terminate",
				Description: "Terminate the interaction when the request is met OR if the assistant cannot proceed further with the task.\nWhen you have finished all the tasks, call this tool to end the work.",
				Parameters:  t,
			},
		},
	}
}

func (t *TerminateTool) GetTool() *model.Tool {
	return t.Tool
}

func (t *TerminateTool) Execute(input string) (string, error) {
	terminateInput := model.TerminateInput{}
	err := json.Unmarshal([]byte(input), &terminateInput)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("交互已完成，状态为：%v", terminateInput.Status), nil
}
