// Package model
/**
@author: xdl2003
@desc:
@date: 2025/6/6
**/
package model

import (
	"fmt"
	"strings"
)

type PlanInfo struct {
	PlanID string
	Title  string
	Steps  []*Step
}

type PlanCommandType string

const (
	PlanCommandTypeCreate    = PlanCommandType("create")
	PlanCommandTypeUpdate    = PlanCommandType("update")
	PlanCommandTypeList      = PlanCommandType("list")
	PlanCommandTypeGet       = PlanCommandType("get")
	PlanCommandTypeSetActive = PlanCommandType("set_active")
	PlanCommandTypeMarkStep  = PlanCommandType("mark_step")
	PlanCommandTypeDelete    = PlanCommandType("delete")
)

type PlanStats struct {
	Completed  int
	InProgress int
	Blocked    int
	NotStarted int
	Total      int
}

// GetAllStatuses 返回所有可能的步骤状态值
func GetAllStatuses() []string {
	return []string{string(NotStarted), string(InProgress), string(Completed), string(Blocked)}
}

// GetActiveStatuses 返回代表活动状态的值列表
func GetActiveStatuses() []string {
	return []string{string(NotStarted), string(InProgress)}
}

func (pi *PlanInfo) GetStats() *PlanStats {
	totalSteps := len(pi.Steps)
	completed := 0
	inProgress := 0
	blocked := 0
	notStarted := 0
	for _, step := range pi.Steps {
		switch step.Status {
		case NotStarted:
			notStarted += 1
		case Blocked:
			blocked += 1
		case Completed:
			completed += 1
		case InProgress:
			inProgress += 1
		}
	}
	return &PlanStats{
		Completed:  completed,
		InProgress: inProgress,
		Blocked:    blocked,
		NotStarted: notStarted,
		Total:      totalSteps,
	}
}

func (pi *PlanInfo) String() string {
	output := fmt.Sprintf("Plan: %s (ID: %s)\\n", pi.Title, pi.PlanID)
	output += strings.Repeat("=", len(output)) + "\n\n"

	stats := pi.GetStats()
	output += fmt.Sprintf("Progress: %d/%d steps completed ", stats.Completed, stats.Total)
	if stats.Total > 0 {
		percentage := float64(stats.Completed) / float64(stats.Total) * 100
		output += fmt.Sprintf("(%.1f%%)\n", percentage)
	} else {
		output += fmt.Sprint("(0%%)\n")
	}
	output += fmt.Sprintf("Status: %d completed, %d in progress, %d blocked, %d not started\n\n",
		stats.Completed, stats.InProgress, stats.Blocked, stats.NotStarted)
	output += fmt.Sprintf("Steps:\n")
	for index, step := range pi.Steps {
		output += fmt.Sprintf("%d. %s %s\n", index, GetStatusMarks(step.Status), step.Data)
		if len(step.Notes) > 0 {
			output += fmt.Sprintf("   Notes: %s\n", step.Notes)
		}
	}
	return output
}

type PlanCommand struct {
	Command    string   `json:"command" description:"The command to execute. Available commands: create, update, list, get, set_active, mark_step, delete."`
	PlanID     string   `json:"plan_id" description:"Unique identifier for the plan. Required for create, update, set_active, and delete commands. Optional for get and mark_step (uses active plan if not specified)." enum:"create,update,set_active,delete,get,mark_step"`
	Title      string   `json:"title" description:"Title for the plan. Required for create command, optional for update command."`
	Steps      []string `json:"steps" description:"List of plan steps. Required for create command, optional for update command."`
	StepIndex  int      `json:"step_index" description:"Index of the step to update (0-based). Required for mark_step command."`
	StepStatus string   `json:"step_status" description:"Status to set for a step. Used with mark_step command. Available status: not_started, in_progress, completed, blocked." enum:"not_started,in_progress,completed,blocked"`
	StepNotes  string   `json:"step_notes" description:"Additional notes for a step. Optional for mark_step command."`
}

type Step struct {
	Data   string
	Status StepStatus
	Notes  string
}

type StepStatus string

const (
	NotStarted StepStatus = "not_started"
	InProgress StepStatus = "in_progress"
	Completed  StepStatus = "completed"
	Blocked    StepStatus = "blocked"
)

var (
	StatusMarks = map[StepStatus]string{
		Completed:  "[✓]",
		InProgress: "[→]",
		Blocked:    "[!]",
		NotStarted: "[ ]",
	}
)

// GetStatusMarks 返回状态到标记符号的映射
func GetStatusMarks(status StepStatus) string {
	mark, ok := StatusMarks[status]
	if !ok {
		return "[ ]"
	}
	return mark
}
