// Package llm
/**
@author: xdl2003
@desc: 大模型调用模块
@date: 2025/6/5
**/
package llm

import (
	"context"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	doubao_m "github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"go-manus/go-manus/config"
	"go-manus/go-manus/log"
	"go-manus/go-manus/model"
	"go.uber.org/zap"
)

type ModelSource string

const (
	ModelSourceDoubao = "doubao"
)

type Client struct {
	ModelConfig  *config.ModelConfig
	DoubaoClient *arkruntime.Client
}

func NewClient(config *config.ModelConfig) (*Client, error) {
	client := &Client{}
	client.ModelConfig = config
	client.DoubaoClient = arkruntime.NewClientWithApiKey(config.ApiKey)
	log.Logger.Info("llm client init success")
	return client, nil
}

type AskRequest struct {
	Messages []*model.Message
	Tools    []*model.Tool
}

func ConvertMessagesToDoubaoMessages(messages []*model.Message) []*doubao_m.ChatCompletionMessage {
	doubaoMessages := make([]*doubao_m.ChatCompletionMessage, 0)
	for _, m := range messages {
		doubaoMessage := &doubao_m.ChatCompletionMessage{
			Role:       string(m.Role),
			Content:    &doubao_m.ChatCompletionMessageContent{StringValue: &m.Content},
			ToolCalls:  nil,
			ToolCallID: m.ToolCallID,
		}
		if len(m.ReasonContent) > 0 {
			doubaoMessage.ReasoningContent = &m.ReasonContent
		}
		if len(m.ToolCalls) > 0 {
			doubaoMessage.ToolCalls = convertToolCall2Doubao(m.ToolCalls)
		}
		doubaoMessages = append(doubaoMessages, doubaoMessage)
	}
	return doubaoMessages
}

func convertToolCall2Doubao(toolCalls []*model.ToolCall) []*doubao_m.ToolCall {
	results := make([]*doubao_m.ToolCall, 0)
	for _, tl := range toolCalls {
		result := &doubao_m.ToolCall{
			ID:   tl.ID,
			Type: doubao_m.ToolType(tl.Type),
			Function: doubao_m.FunctionCall{
				Name:      tl.Function.Name,
				Arguments: tl.Function.Arguments,
			},
			Index: tl.Index,
		}
		results = append(results, result)
	}
	return results
}

func ConvertDoubaoMessageToMessage(doubaoMessage *doubao_m.ChatCompletionMessage) *model.Message {
	result := model.Message{
		Role:          model.RoleType(doubaoMessage.Role),
		Content:       "",
		ReasonContent: "",
		ToolCalls:     nil,
		ToolCallID:    doubaoMessage.ToolCallID,
	}
	if doubaoMessage.Content != nil && doubaoMessage.Content.StringValue != nil {
		result.Content = *doubaoMessage.Content.StringValue
	}
	if doubaoMessage.ReasoningContent != nil {
		result.ReasonContent = *doubaoMessage.ReasoningContent
	}
	if len(doubaoMessage.ToolCalls) > 0 {
		result.ToolCalls = convertDoubaoToolCall2Type(doubaoMessage.ToolCalls)
	}
	return &result
}

func convertDoubaoToolCall2Type(toolCalls []*doubao_m.ToolCall) []*model.ToolCall {
	results := make([]*model.ToolCall, 0)
	for _, tl := range toolCalls {
		result := &model.ToolCall{
			ID:   tl.ID,
			Type: model.ToolType(tl.Type),
			Function: model.FunctionCall{
				Name:      tl.Function.Name,
				Arguments: tl.Function.Arguments,
			},
			Index: tl.Index,
		}
		results = append(results, result)
	}
	return results
}

type ChatResp struct {
	Message *model.Message `json:"message"`
	Error   error          `json:"error"`
}

func convertTool2Doubao(tools []*model.Tool) []*doubao_m.Tool {
	results := make([]*doubao_m.Tool, 0)
	for _, tool := range tools {
		result := &doubao_m.Tool{
			Type: doubao_m.ToolType(tool.Type),
			Function: &doubao_m.FunctionDefinition{
				Name:        tool.Function.Name,
				Description: tool.Function.Description,
				Parameters:  tool.Function.Parameters,
			},
		}
		results = append(results, result)
	}
	return results
}

func (c *Client) AskTool(req *AskRequest) *ChatResp {
	log.Logger.Info("ask tool, req=%v", zap.Any("req", req))
	doubaoReq := doubao_m.CreateChatCompletionRequest{
		Messages:   ConvertMessagesToDoubaoMessages(req.Messages),
		Tools:      convertTool2Doubao(req.Tools),
		Model:      "doubao-1.5-pro-32k-250115",
		ToolChoice: doubao_m.ToolChoiceStringTypeAuto,
	}
	completion, err := c.DoubaoClient.CreateChatCompletion(context.Background(), doubaoReq)
	resp := &ChatResp{}
	if err != nil {
		log.Logger.Error("fail to ask tool, err=%v", zap.Error(err))
		resp.Error = err
	} else {
		resp.Message = ConvertDoubaoMessageToMessage(&completion.Choices[0].Message)
	}
	return resp
}
