package main

import (
	"context"
	"fmt"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
	"go-manus/go-manus/agent"
	"go-manus/go-manus/config"
	"go-manus/go-manus/log"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()
	log.InitLogger()
	defer log.Logger.Sync()
	log.Logger.Info("/go-manus start/\n")
	config.InitConfig()
	manus, err := agent.NewManus()
	if err != nil {
		log.Logger.Error("fail to new manus, err=%v", zap.Error(err))
		return
	}
	// 构建聊天完成请求，设置请求的模型和消息内容
	req := model.ChatCompletionRequest{
		// 将推理接入点 <Model>替换为 Model ID
		Model: "doubao-1.5-pro-32k-250115",
		Messages: []*model.ChatCompletionMessage{
			{
				// 消息的角色为用户
				Role: model.ChatMessageRoleUser,
				Content: &model.ChatCompletionMessageContent{
					StringValue: volcengine.String("你好"),
				},
			},
		},
	}

	// 发送聊天完成请求，并将结果存储在 resp 中，将可能出现的错误存储在 err 中
	resp, err := manus.LLM.DoubaoClient.CreateChatCompletion(ctx, req)
	if err != nil {
		// 若出现错误，打印错误信息并终止程序
		fmt.Printf("standard chat error: %v\n", err)
		return
	}
	// 打印聊天完成请求的响应结果
	fmt.Println(*resp.Choices[0].Message.Content.StringValue)
}
