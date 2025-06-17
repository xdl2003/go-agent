// Package model
/**
@author: xdl2003
@desc:
@date: 2025/6/9
**/
package model

// Message
//
//	@Description: 发送给大模型的消息结构
type Message struct {

	//  Role
	//  @Description: 消息角色
	Role RoleType `json:"role"`

	//  Content
	//  @Description: 消息内容
	Content string `json:"content"`

	//  ReasonContent
	//  @Description:
	ReasonContent string `json:"reason_content"`

	//  ToolCalls
	//  @Description: 工具调用
	ToolCalls []*ToolCall `json:"tool_calls"`

	//  Name
	//  @Description:
	Name string `json:"name"`

	//  ToolCallID
	//  @Description: 工具调用ID
	ToolCallID string `json:"tool_call_id"`

	//  Base64Image
	//  @Description: 图片base64编码
	Base64Image string `json:"base64_image"`
}

type RoleType string

const (
	RoleUser      = RoleType("user")
	RoleSystem    = RoleType("system")
	RoleAssistant = RoleType("assistant")
	RoleTool      = RoleType("tool")
)

// NewUserMessage 创建用户消息
func NewUserMessage(content string, base64Image string) *Message {
	return &Message{
		Role:        RoleUser,
		Content:     content,
		Base64Image: base64Image,
	}
}

// NewSystemMessage 创建系统消息
func NewSystemMessage(content string) *Message {
	return &Message{
		Role:    RoleSystem,
		Content: content,
	}
}

// NewAssistantMessage 创建助手消息
func NewAssistantMessage(content string, reasonContent string, toolCalls []*ToolCall) *Message {
	return &Message{
		Role:          RoleAssistant,
		Content:       content,
		ReasonContent: reasonContent,
		ToolCalls:     toolCalls,
	}
}

// NewToolMessage 创建工具消息
func NewToolMessage(content string, toolCallID string, name string, base64Image string) *Message {
	return &Message{
		Role:        RoleTool,
		Content:     content,
		ToolCallID:  toolCallID,
		Name:        name,
		Base64Image: base64Image,
	}
}
