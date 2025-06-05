// Package agent
/**
@author: xudongliu.666@bytedance.com
@desc:
@date: 2025/6/5
**/
package agent

type Manus struct {
	*ReActAgent
}

func NewManus() (*Manus, error) {
	agent, err := NewReActAgent()
	if err != nil {
		return nil, err
	}
	manus := &Manus{
		ReActAgent: agent,
	}
	return manus, nil
}
