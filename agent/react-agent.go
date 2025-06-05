// Package agent
/**
@author: xudongliu.666@bytedance.com
@desc:
@date: 2025/6/5
**/
package agent

import (
	"go-manus/go-manus/config"
	"go-manus/go-manus/llm"
	"go-manus/go-manus/log"
	"go.uber.org/zap"
	"sync"
)

type ReActAgent struct {
	Name           string
	Description    string
	SystemPrompt   string
	NextStepPrompt string
	LLM            *llm.Client
	// Memory             *schema.Memory
	// State              schema.AgentState
	MaxSteps           int
	CurrentStep        int
	DuplicateThreshold int
	mu                 sync.RWMutex

	// AvailableTools     map[string]tool.ToolIF
	// ToolChoices        tool.ToolChoiceType
	SpecialToolNames []string
	// ToolCalls          []*schema.ToolCall
	CurrentBase64Image string
	MaxObserve         int
}

func NewReActAgent() (*ReActAgent, error) {
	agent := &ReActAgent{}
	client, err := llm.NewClient(config.AllConfig.PrimaryConfig)
	if err != nil {
		log.Logger.Error("fail to new client client, err=%v", zap.Error(err))
		return nil, err
	}
	agent.LLM = client
	return agent, nil
}
