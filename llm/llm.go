// Package llm
/**
@author: xudongliu.666@bytedance.com
@desc:
@date: 2025/6/5
**/
package llm

import (
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	"go-manus/go-manus/config"
	"go-manus/go-manus/log"
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
