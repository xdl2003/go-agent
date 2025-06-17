// Package model
/**
@author: xdl2003
@desc:
@date: 2025/6/16
**/
package model

import "sync"

type Memory struct {
	Messages    []*Message
	MaxMessages int
	mu          sync.RWMutex
}

func NewMemory() *Memory {
	return &Memory{
		Messages:    make([]*Message, 0),
		MaxMessages: 100,
		mu:          sync.RWMutex{},
	}
}

func (m *Memory) AddMessage(msg *Message) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Messages = append(m.Messages, msg)
	if len(m.Messages) > m.MaxMessages {
		m.Messages = m.Messages[1:]
	}
}

func (m *Memory) AddMessages(msgs []*Message) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Messages = append(m.Messages, msgs...)
	if len(m.Messages) > m.MaxMessages {
		m.Messages = m.Messages[-m.MaxMessages:]
	}
}

func (m *Memory) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Messages = make([]*Message, 0)
}

func (m *Memory) GetRecentMessages(n int) []*Message {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.Messages[-n:]
}

// GetMessages 获取所有消息
func (m *Memory) GetMessages() []*Message {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return append([]*Message{}, m.Messages...)
}
