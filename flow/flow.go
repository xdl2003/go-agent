// Package flow
/**
@author: xdl2003
@desc:
@date: 2025/6/6
**/
package flow

import (
	"fmt"
	"go-manus/go-manus/agent"
	"go-manus/go-manus/config"
	"go-manus/go-manus/llm"
	"go-manus/go-manus/log"
	"go-manus/go-manus/model"
	tool2 "go-manus/go-manus/tool"
	"go-manus/go-manus/util"
	"go.uber.org/zap"
)

// Flow
// @Description: 整体工作流
type Flow struct {
	//  agents
	//  @Description: 全部Agent
	agents map[string]agent.ToolCallAgentIF

	//  tools
	//  @Description: 全部工具
	tools map[string]model.BaseTool

	//  primaryAgentKey
	//  @Description: 主Agent的Key
	primaryAgentKey *string

	//  llm
	//  @Description: 大模型，用于生成和结束plan
	llm *llm.Client

	//  PlanTool
	//  @Description: 流程规划工具
	PlanTool *tool2.PlanTool
}

func NewFlow() *Flow {
	agents := make(map[string]agent.ToolCallAgentIF)
	tools := make(map[string]model.BaseTool)
	primaryAgentKey := "manus"
	llm, _ := llm.NewClient(config.AllConfig.PrimaryConfig)
	agents[primaryAgentKey], _ = agent.NewManus()
	planTool := tool2.NewPlanTool()
	return &Flow{
		agents:          agents,
		tools:           tools,
		primaryAgentKey: &primaryAgentKey,
		llm:             llm,
		PlanTool:        planTool,
	}
}

func (f *Flow) Execute(prompt *string) (string, error) {
	err := f.CreateInitialPlan(prompt)
	if err != nil {
		return "", err
	}

	result := ""

	for {
		currentStepIndex, stepInfo, err := f.getCurrentStepInfo()
		if err != nil {
			return "", err
		}

		if currentStepIndex < 0 {
			finalResult, err := f.finalizePlan()
			if err != nil {
				return "", err
			}
			result += finalResult
			fmt.Println("task completed: ", finalResult)
			break
		}

		stepResult, err := f.executeStep(stepInfo, int64(currentStepIndex))
		if err != nil {
			return "", err
		}
		result += stepResult + "\n"
		if f.agents[*f.primaryAgentKey].GetStatus() == model.AgentStateFINISHED {
			break
		}
	}

	return result, nil
}

// executeStep 执行当前步骤
func (f *Flow) executeStep(stepInfo *model.Step, currentStepIndex int64) (string, error) {
	planStatus, err := f.getPlanText()
	if err != nil {
		return "", err
	}

	stepText := stepInfo.Data
	if stepText == "" {
		stepText = fmt.Sprintf("Step %d", f.PlanTool.CurrentStepIndex)
	}
	stepPrompt := model.GetStepPrompt(planStatus, currentStepIndex, stepText)
	stepResult, err := f.agents["manus"].Run(stepPrompt)
	if err != nil {
		return "", err
	}
	if err := f.markStepCompleted(currentStepIndex); err != nil {
		return "", err
	}

	return stepResult, nil
}

// markStepCompleted 标记当前步骤为已完成
func (f *Flow) markStepCompleted(currentStepIndex int64) error {
	if currentStepIndex < 0 {
		return nil
	}
	command := model.PlanCommand{
		Command:    string(model.PlanCommandTypeMarkStep),
		PlanID:     f.PlanTool.CurrentPlanID,
		StepIndex:  int(currentStepIndex),
		StepStatus: string(model.Completed),
		StepNotes:  "",
	}
	result, err := f.PlanTool.Execute(util.MustJson(command))
	if err != nil {
		return err
	}
	fmt.Printf("Marked step %s as completed in plan %d, result=%v\n", f.PlanTool.CurrentPlanID, int(currentStepIndex), result)
	return nil
}

// getCurrentStepInfo 获取当前步骤信息
func (f *Flow) getCurrentStepInfo() (int, *model.Step, error) {

	planData := f.PlanTool.GetPlan(f.PlanTool.CurrentPlanID)
	for stepIndex, stepInfo := range planData.Steps {
		if !util.ContainArrStr(model.GetActiveStatuses(), string(stepInfo.Status)) {
			continue
		}
		result, err := f.PlanTool.MarkStep(f.PlanTool.CurrentPlanID, stepIndex, string(model.InProgress), "")
		if err != nil {
			log.Logger.Error("fail marking step as in_progress, err=%v", zap.Any("err", err))
			return -1, nil, err
		}
		log.Logger.Info("succ marking step as in_progress, result=%v", zap.Any("result", result))
		return stepIndex, stepInfo, nil
	}
	return -1, nil, nil
}

// finalizePlan 完成计划并提供总结
func (f *Flow) finalizePlan() (string, error) {
	planText, err := f.getPlanText()
	if err != nil {
		return "", err
	}

	systemMessage := model.NewSystemMessage(model.GetFinalizeSystemPrompt())
	userMessage := model.NewUserMessage(model.GetFinalizeUserPrompt(&planText), "")

	chatReq := &llm.AskRequest{
		Messages: []*model.Message{systemMessage, userMessage},
	}
	chatResp := f.llm.Ask(chatReq)
	if chatResp.Error != nil {
		return "", fmt.Errorf("plan completed. Error generating summary")
	}

	return fmt.Sprintf("Plan completed:\n\n%s", chatResp.Message.Content), nil
}

// getPlanText 获取当前计划的格式化文本
func (f *Flow) getPlanText() (string, error) {
	command := model.PlanCommand{
		Command: string(model.PlanCommandTypeGet),
		PlanID:  f.PlanTool.CurrentPlanID,
	}
	result, err := f.PlanTool.Execute(util.MustJson(command))
	if err != nil {
		return "", err
	}
	return result, nil
}

// CreateInitialPlan
//
//	@Description: 创建初始计划
//	@receiver f
func (f *Flow) CreateInitialPlan(prompt *string) error {
	log.Logger.Info("create initial plan")
	request := &llm.AskRequest{}
	request.Messages = []*model.Message{
		model.NewSystemMessage(model.GetPlanSystemPrompt()),
		model.NewUserMessage(model.GetPlanUserPrompt(prompt), ""),
	}
	request.Tools = []*model.Tool{
		f.PlanTool.GetTool(),
	}
	chatResp := f.llm.AskTool(request)
	log.Logger.Info("create initial plan, resp=%v", zap.Any("resp", chatResp))
	for _, toolCall := range chatResp.Message.ToolCalls {
		if toolCall.Function.Name == f.PlanTool.GetTool().Function.Name {
			_, err := f.PlanTool.Execute(toolCall.Function.Arguments)
			//fmt.Println("create initial plan, rep=", resp, "err=", err)
			//fmt.Println("plan=", util.MustJson(f.PlanTool.GetAllPlans()))
			if err != nil {
				return err
			}
			for plan, _ := range f.PlanTool.GetAllPlans() {
				f.PlanTool.SetActive(plan)
				break
			}
			return nil
		}
	}
	return nil
}
