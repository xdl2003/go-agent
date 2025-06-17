// Package tool
/**
@author: xdl2003
@desc: 流程规划工具
@date: 2025/6/5
**/
package tool

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"github.com/sashabaranov/go-openai/jsonschema"
	"go-manus/go-manus/model"
	"go-manus/go-manus/util"
	"sync"
	"time"
)

type PlanTool struct {
	*model.Tool
	Plans            map[string]*model.PlanInfo
	mu               sync.RWMutex
	CurrentPlanID    string
	CurrentStepIndex int
}

func NewPlanTool() *PlanTool {
	p, _ := jsonschema.GenerateSchemaForType(model.PlanCommand{})
	return &PlanTool{
		Tool: &model.Tool{
			Type: model.ToolTypeFunction,
			Function: model.FunctionDefinition{
				Name: "plan",
				Description: "A planning tool that allows the agent to create and manage plans for solving complex tasks.\n" +
					"The tool provides functionality for creating plans, updating plan steps, and tracking progress.",
				Parameters: p,
			},
		},
		Plans:         make(map[string]*model.PlanInfo),
		mu:            sync.RWMutex{},
		CurrentPlanID: "",
	}
}

func (pt *PlanTool) GetTool() *model.Tool {
	return pt.Tool
}

func (pt *PlanTool) Execute(input string) (string, error) {
	// fmt.Println("调用工具PlanTool, input=", input)
	var command model.PlanCommand
	err := jsoniter.UnmarshalFromString(input, &command)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Invalid input: %s", err.Error()))
	}
	if command.PlanID == "" {
		command.PlanID = fmt.Sprintf("plan_%s", time.Now().Format(time.DateTime))
	} else if command.Command == "create" {
		command.PlanID += time.Now().Format(time.DateTime)
	}
	switch command.Command {
	case "create":
		return pt.CreatePlan(command)
	case "update":
		return pt.UpdatePlan(command)
	case "list":
		return pt.List()
	case "get":
		return pt.Get(command.PlanID)
	case "set_active":
		return pt.SetActive(command.PlanID)
	case "mark_step":
		return pt.MarkStep(command.PlanID, command.StepIndex, command.StepStatus, command.StepNotes)
	case "delete":
		return pt.DeletePlan(command.PlanID)
	default:
		return "", errors.New(fmt.Sprintf("Invalid command: %s", command.Command))
	}
}

func (pt *PlanTool) GetPlan(planID string) *model.PlanInfo {
	pt.mu.RLock()
	defer pt.mu.RUnlock()
	return pt.Plans[planID]
}

func (pt *PlanTool) DelPlan(planID string) {
	pt.mu.Lock()
	defer pt.mu.Lock()
	delete(pt.Plans, planID)
}

func (pt *PlanTool) SetPlan(planInfo *model.PlanInfo) {
	pt.mu.Lock()
	defer pt.mu.Unlock()
	pt.Plans[planInfo.PlanID] = planInfo
}

func (pt *PlanTool) GetAllPlans() map[string]*model.PlanInfo {
	m := make(map[string]*model.PlanInfo)
	for key, value := range pt.Plans {
		m[key] = value
	}
	return m
}

func (pt *PlanTool) CreatePlan(pc model.PlanCommand) (string, error) {
	planID := pc.PlanID
	title := pc.Title
	steps := pc.Steps
	planningInfo := &model.PlanInfo{
		PlanID: planID,
		Title:  title,
		Steps:  make([]*model.Step, 0),
	}
	for _, step := range steps {
		planningInfo.Steps = append(planningInfo.Steps, &model.Step{
			Data:   step,
			Status: model.NotStarted,
			Notes:  "",
		})
	}
	pt.Plans[planID] = planningInfo
	return fmt.Sprintf("Plan %s created", planID), nil
}

func (pt *PlanTool) UpdatePlan(pc model.PlanCommand) (string, error) {
	if len(pc.PlanID) == 0 {
		return "", errors.New("parameter `plan_id` is required for command: update")
	}
	getPlan := pt.GetPlan(pc.PlanID)
	if getPlan == nil {
		return "", errors.New(fmt.Sprintf("No getPlan found with ID: %s", pc.PlanID))
	}
	if len(pc.Title) > 0 {
		getPlan.Title = pc.Title
	}
	if len(pc.Steps) > 0 {
		newStepInfos := make([]*model.Step, 0)
		for i, step := range pc.Steps {
			if i < len(getPlan.Steps) && step == getPlan.Steps[i].Data {
				newStepInfos = append(newStepInfos, getPlan.Steps[i])
			} else {
				newStepInfo := &model.Step{Data: step, Status: model.NotStarted, Notes: ""}
				newStepInfos = append(newStepInfos, newStepInfo)
			}
		}
	}
	pt.SetPlan(getPlan)
	return fmt.Sprintf("Plan %s updated", pc.PlanID), nil
}

func (pt *PlanTool) List() (string, error) {
	if len(pt.Plans) == 0 {
		return "No plans available. Create a plan with the 'create' command.", nil
	}
	output := "Available plans:\n"
	for planID, plan := range pt.GetAllPlans() {
		currentMarker := ""
		if planID == pt.CurrentPlanID {
			currentMarker = " (active)"
		}
		stats := plan.GetStats()
		progress := fmt.Sprintf("%d/%d steps completed", stats.Completed, stats.Total)
		output += fmt.Sprintf("• %s%s: %s - %s\n", planID, currentMarker, plan.Title, progress)
	}
	return output, nil
}

func (pt *PlanTool) Get(planID string) (string, error) {
	if len(planID) == 0 {
		return "", errors.New("parameter `plan_id` is required for command: get")
	}
	plan := pt.GetPlan(planID)
	if plan == nil {
		return fmt.Sprintf("Plan %s not found", planID), nil
	} else {
		return plan.String(), nil
	}
}

func (pt *PlanTool) SetActive(planID string) (string, error) {
	if len(planID) == 0 {
		return "", errors.New("parameter `plan_id` is required for command: set_active")
	}
	plan := pt.GetPlan(planID)
	if plan == nil {
		return "", fmt.Errorf("no plan found with ID: %s", planID)
	}
	pt.CurrentPlanID = planID
	return fmt.Sprintf("Plan '%s' is now the active plan.\n\n%s", planID, plan.String()), nil
}

func (pt *PlanTool) MarkStep(planID string, stepIndex int, stepStatus string, stepNotes string) (string, error) {
	if len(planID) == 0 {
		if len(pt.CurrentPlanID) == 0 {
			return "", errors.New("no active getPlan. Please specify a plan_id or set an active getPlan")
		}
		planID = pt.CurrentPlanID
	}
	getPlan := pt.GetPlan(planID)
	if getPlan == nil {
		return "", fmt.Errorf("no getPlan found with ID: %s", planID)
	}

	if stepIndex < 0 || stepIndex > len(getPlan.Steps) {
		return "", fmt.Errorf("invalid step_index: %d. Valid indices range from 0 to %d", stepIndex, len(getPlan.Steps)-1)
	}

	if len(stepStatus) > 0 && !util.ContainArrStr(model.GetAllStatuses(), stepStatus) {
		return "", fmt.Errorf("invalid step_status: %s. Valid statuses are: not_started, in_progress, completed, blocked", stepStatus)
	}

	if len(stepStatus) > 0 {
		getPlan.Steps[stepIndex].Status = model.StepStatus(stepStatus)
	}

	if len(stepNotes) > 0 {
		getPlan.Steps[stepIndex].Notes = stepNotes
	}

	return fmt.Sprintf("Step %d updated in getPlan '%s'.\n\n%s", stepIndex, planID, getPlan.String()), nil
}

func (pt *PlanTool) DeletePlan(planID string) (string, error) {
	if len(planID) == 0 {
		return "", errors.New("parameter `plan_id` is required for command: delete")
	}
	plan := pt.GetPlan(planID)
	if plan == nil {
		return "", errors.New(fmt.Sprintf("No plan found with ID: %s", planID))
	}

	pt.DelPlan(planID)

	if pt.CurrentPlanID == planID {
		pt.CurrentPlanID = ""
	}

	return fmt.Sprintf("Plan '%s' has been deleted.", planID), nil
}
