// Package agent
/**
@author: xdl2003
@desc: 基础的交互agent
@date: 2025/6/5
**/
package agent

import (
	"go-manus/go-manus/config"
	"go-manus/go-manus/llm"
	"go-manus/go-manus/log"
	"go-manus/go-manus/model"
	"go.uber.org/zap"
	"sync"
)

type ReActAgent struct {
	Name           string
	Description    string
	SystemPrompt   string
	NextStepPrompt string
	LLM            *llm.Client
	Memory         *model.Memory
	// State              schema.AgentState
	MaxSteps           int
	CurrentStep        int
	DuplicateThreshold int
	mu                 sync.RWMutex
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
	agent.NextStepPrompt = model.GetNextStepPrompt()
	agent.SystemPrompt = model.GetSystemPrompt()
	agent.MaxSteps = 100
	agent.Memory = model.NewMemory()
	agent.MaxObserve = 100000
	return agent, nil
}
