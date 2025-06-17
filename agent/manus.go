// Package agent
/**
@author: xdl2003
@desc: manus agent
@date: 2025/6/5
**/
package agent

type Manus struct {
	*ToolCallAgent
}

func NewManus() (*Manus, error) {
	agent, err := NewToolCallAgent()
	if err != nil {
		return nil, err
	}
	manus := &Manus{
		ToolCallAgent: agent,
	}
	return manus, nil
}
