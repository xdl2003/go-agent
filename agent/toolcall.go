// Package agent
/**
@author: xdl2003
@desc:
@date: 2025/6/6
**/
package agent

import (
	"errors"
	"fmt"
	"go-manus/go-manus/llm"
	"go-manus/go-manus/model"
	"go-manus/go-manus/tool"
	"strings"
	"time"
)

type ToolCallAgentIF interface {
	GetName() string
	Run(request string) (string, error)
	GetStatus() model.AgentState
}

type ToolCallAgent struct {
	*ReActAgent
	// AvailableTools     map[string]tool.ToolIF
	// ToolChoices        tool.ToolChoiceType
	ToolCalls        []*model.ToolCall
	AvailableTools   map[string]model.BaseTool
	SpecialToolNames []string
	State            model.AgentState
}

func (tc *ToolCallAgent) GetStatus() model.AgentState {
	return tc.State
}

func NewToolCallAgent() (*ToolCallAgent, error) {
	agent, err := NewReActAgent()
	if err != nil {
		return nil, err
	}
	toolCallAgent := &ToolCallAgent{
		ReActAgent:     agent,
		State:          model.AgentStateIDLE,
		AvailableTools: tool.GetAvailableTools(),
	}
	return toolCallAgent, nil
}

func (tc *ToolCallAgent) GetName() string {
	return "manus"
}

func (tc *ToolCallAgent) Run(request string) (string, error) {
	fmt.Println("调用Manus, request=", request)
	tc.mu.RLock()
	if tc.State != model.AgentStateIDLE {
		tc.mu.RUnlock()
		return "", fmt.Errorf("cannot run agent from state: %s", tc.State.String())
	}
	tc.mu.RUnlock()

	//if request != "" {
	//	if err := a.UpdateMemory(schema.RoleUser, request, "", "", nil, "", ""); err != nil {
	//		return "", fmt.Errorf("failed to update memory: %w", err)
	//	}
	//}
	tc.State = model.AgentStateRUNNING
	defer func() {
		tc.State = model.AgentStateIDLE
	}()
	results := make([]string, 0)

	//stateCtx, err := tc.NewStateContext(model.AgentStateRUNNING)
	//if err != nil {
	//	return "", fmt.Errorf("failed to create state context: %w", err)
	//}
	//defer stateCtx.Done()

	for tc.CurrentStep < tc.MaxSteps && tc.State == model.AgentStateRUNNING {

		tc.mu.Lock()
		tc.CurrentStep++
		currentStep := tc.CurrentStep
		maxSteps := tc.MaxSteps
		tc.mu.Unlock()
		fmt.Println("Excuting step %d/%d", currentStep, maxSteps)

		// 执行步骤（由子类实现）
		stepResult, err := tc.Step()
		if err != nil {
			return "", fmt.Errorf("step %d failed: %w", currentStep, err)
		}

		//// 检查是否陷入停滞状态
		//if tc.IsStuck() {
		//	tc.HandleStuckState()
		//}

		results = append(results, fmt.Sprintf("Step %d: %s", currentStep, stepResult))

		// 短暂休眠，避免CPU占用过高
		time.Sleep(100 * time.Millisecond)
	}

	if tc.CurrentStep >= tc.MaxSteps {
		tc.mu.Lock()
		tc.CurrentStep = 0
		tc.State = model.AgentStateIDLE
		tc.mu.Unlock()
		results = append(results, fmt.Sprintf("Terminated: Reached max steps (%d)", tc.MaxSteps))
	}

	// 清理资源
	//if err := sandbox.Cleanup(ctx); err != nil {
	//	logs.Warnf("Failed to cleanup sandbox: %v", err)
	//}

	if len(results) == 0 {
		return "No steps executed", nil
	}

	return fmt.Sprintf("%s", results), nil
}

// Think 处理当前状态并使用工具决定下一步行动
func (tc *ToolCallAgent) Think() (bool, error) {
	if tc.NextStepPrompt != "" {
		userMsg := model.NewUserMessage(tc.NextStepPrompt, "")
		tc.Memory.AddMessage(userMsg)
	}

	messages := tc.Memory.GetMessages()
	var systemMsgs []*model.Message

	if tc.SystemPrompt != "" {
		systemMsgs = []*model.Message{model.NewSystemMessage(tc.SystemPrompt)}
	}

	userMessage := model.NewUserMessage(tc.NextStepPrompt, "")

	request := &llm.AskRequest{}
	request.Messages = []*model.Message{}
	request.Messages = append(request.Messages, messages...)
	request.Messages = append(request.Messages, systemMsgs...)
	request.Messages = append(request.Messages, userMessage)

	request.Tools = []*model.Tool{}
	for _, tool := range tc.AvailableTools {
		request.Tools = append(request.Tools, tool.GetTool())
	}
	chatResp := tc.LLM.AskTool(request)

	if len(chatResp.Message.ToolCalls) == 0 && len(chatResp.Message.Content) == 0 &&
		len(chatResp.Message.Base64Image) == 0 && len(chatResp.Message.ReasonContent) == 0 {
		return false, errors.New("no response received from the LLM")
	}
	tc.ToolCalls = chatResp.Message.ToolCalls
	content := chatResp.Message.Content
	reasonContent := chatResp.Message.ReasonContent

	// 记录响应信息
	fmt.Printf("✨ %s's thoughts: %s, reason: %s\n", tc.Name, content, reasonContent)
	fmt.Printf("🛠️ %s selected %d tools to use\n", tc.Name, len(tc.ToolCalls))
	if len(tc.ToolCalls) > 0 {
		toolNames := make([]string, 0, len(tc.ToolCalls))
		for _, call := range tc.ToolCalls {
			toolNames = append(toolNames, call.Function.Name)
		}
		fmt.Printf("🧰 Tools being prepared: %v\n", toolNames)
		fmt.Printf("🔧 Tool arguments: %v\n", tc.ToolCalls[0].Function.Arguments)
	}

	tc.Memory.AddMessage(model.NewAssistantMessage(content, reasonContent, tc.ToolCalls))

	return len(tc.ToolCalls) > 0, nil
}

// Act 执行步骤
func (tc *ToolCallAgent) Act() (string, error) {
	if len(tc.ToolCalls) == 0 {
		return "No content or toolcalls to execute", nil
	}
	results := make([]string, 0, len(tc.ToolCalls))
	for _, command := range tc.ToolCalls {
		result, err := tc.executeTool(command)
		if err != nil {
			return "", err
		}
		if tc.MaxObserve > 0 && len(result) > tc.MaxObserve {
			result = result[:tc.MaxObserve]
		}
		fmt.Printf("🎯 Tool '%s' completed its mission! Result: %s\n", command.Function.Name, result)

		// 将工具响应添加到内存
		toolMsg := model.NewToolMessage(result, command.ID, command.Function.Name, tc.CurrentBase64Image)
		tc.Memory.AddMessage(toolMsg)
		results = append(results, result)
	}
	return strings.Join(results, "\n"), nil
}

// Step 执行Agent工作流程中的单个步骤
func (tc *ToolCallAgent) Step() (string, error) {
	// 思考阶段：决定要使用哪些工具
	shouldAct, err := tc.Think()
	if err != nil {
		return "", fmt.Errorf("thinking phase failed: %w", err)
	}

	if !shouldAct {
		return "No action required", nil
	}

	// 行动阶段：执行工具调用
	result, err := tc.Act()
	if err != nil {
		return "", fmt.Errorf("action phase failed: %w", err)
	}

	return result, nil
}

func (tc *ToolCallAgent) executeTool(command *model.ToolCall) (string, error) {
	// 格式化为显示结果
	var observation string
	if command.Function.Name == "" {
		return "Error: Invalid command format", nil
	}

	if _, ok := tc.AvailableTools[command.Function.Name]; !ok {
		return fmt.Sprintf("Error: Unknown tool '%s'", command.Function.Name), nil
	}
	// 执行工具
	result, err := tc.AvailableTools[command.Function.Name].Execute(command.Function.Arguments)
	if err != nil {
		errorMsg := fmt.Sprintf("⚠️ Tool '%s' encountered a problem: %v", command.Function.Name, err)
		return fmt.Sprintf("Error: %s", errorMsg), nil
	}

	// 处理特殊工具
	if err := tc.handleSpecialTool(command.Function.Name, result); err != nil {
		return "", err
	}

	if len(result) > 0 {
		observation = fmt.Sprintf("Observed output of cmd `%s` executed:\n%v", command.Function.Name, result)
	} else {
		observation = fmt.Sprintf("Cmd `%s` completed with no output", command.Function)
	}

	return observation, nil
}

// handleSpecialTool 处理特殊工具执行和状态变更
func (tc *ToolCallAgent) handleSpecialTool(name string, result interface{}) error {
	if name != "terminate" {
		return nil
	}
	// 设置代理状态为已完成
	fmt.Printf("🏁 Special tool '%s' has completed the task!\n", name)
	tc.State = model.AgentStateFINISHED
	return nil
}
