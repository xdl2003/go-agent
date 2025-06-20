// Package tool
/**
@author: xdl2003
@desc:
@date: 2025/6/18
**/
package tool

import (
	"encoding/json"
	"fmt"
	"github.com/sashabaranov/go-openai/jsonschema"
	"go-manus/go-manus/model"
	"strings"
)

type AskHumanTool struct {
	*model.Tool
}

func NewAskHumanTool() *AskHumanTool {
	t, _ := jsonschema.GenerateSchemaForType(model.AskHumanInput{})
	return &AskHumanTool{
		Tool: &model.Tool{
			Type: model.ToolTypeFunction,
			Function: model.FunctionDefinition{
				Name:        "ask_human",
				Description: "Ask human for help when the request is met OR if the assistant cannot proceed further with the task.\nWhen you have finished all the tasks, call this tool to end the work.",
				Parameters:  t,
			},
		},
	}
}

func (t *AskHumanTool) GetTool() *model.Tool {
	return t.Tool
}

func (t *AskHumanTool) Execute(input string, method string) (string, error) {
	askHumanInput := model.AskHumanInput{}
	err := json.Unmarshal([]byte(input), &askHumanInput)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Bot: %s\n\nYou: \n\n", askHumanInput.Inquire)
	var response string
	fmt.Scanln(&response)
	return strings.TrimSpace(response), nil
}
