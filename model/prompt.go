// Package model
/**
@author: xdl2003
@desc:
@date: 2025/6/9
**/
package model

import "strconv"

/*****************************
 *  计划相关
 *****************************/

func GetPlanSystemPrompt() string {
	return "You are a planning assistant. Create a concise, actionable plan with clear steps. " +
		"Focus on key milestones rather than detailed sub-steps. " +
		"Optimize for clarity and efficiency."
}

func GetPlanUserPrompt(req *string) string {
	return "Create a reasonable plan with clear steps to accomplish the task: " + *req
}

/*****************************
 *  执行相关
 *****************************/

func GetStepPrompt(planStatus string, currentStepIndex int64, stepText string) string {
	return "CURRENT PLAN STATUS:\n        " +
		planStatus +
		"\n\n        YOUR CURRENT TASK:\n        You are now working on step " +
		strconv.FormatInt(currentStepIndex, 10) +
		": " +
		stepText +
		"\"\n\n        Please execute this step using the appropriate tools. When you're done, provide a summary of what you accomplished."
}

func GetNextStepPrompt() string {
	return "Based on user needs, proactively select the most appropriate tool or combination of tools." +
		" For complex tasks, you can break down the problem and use different tools step by step to solve it." +
		" After using each tool, clearly explain the execution results and suggest the next steps.\n\n\t\t" +
		"If you want to stop the interaction at any point, use the `terminate` tool/function call."
}

/*****************************
 *  manus相关
 *****************************/

func GetSystemPrompt() string {
	return "You are Manus, an all-capable AI assistant, aimed at solving any task presented by the user. " +
		"You have various tools at your disposal that you can call upon to efficiently complete complex requests. " +
		"Whether it's programming, information retrieval, file processing, web browsing, or human interaction (only for extreme cases), you can handle it all.\"" +
		"\n    \"The initial directory is: {directory}\""
}

/*****************************
 *  finalize相关
 *****************************/

func GetFinalizeSystemPrompt() string {
	return "You are a planning assistant. Your task is to summarize the completed plan."
}

func GetFinalizeUserPrompt(planText *string) string {
	return "The plan has been completed. Here is the final plan status:\\n\\n" +
		*planText +
		"\\n\\nPlease provide a summary of what was accomplished and any final thoughts."
}
