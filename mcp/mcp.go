// Package mcp
/**
@author: xudongliu.666@bytedance.com
@desc:
@date: 2025/6/17
**/
package mcp

import (
	"context"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
	"go-manus/go-manus/config"
	"go-manus/go-manus/log"
	"go-manus/go-manus/util"
	"go.uber.org/zap"
)

var (
	AllMcpClient []*client.Client
	ToolMcpMap   map[string]*client.Client
	AllTools     []*mcp.Tool
)

func InitMcp() {
	ToolMcpMap = make(map[string]*client.Client)
	AllMcpClient = make([]*client.Client, 0)
	AllTools = make([]*mcp.Tool, 0)
	for _, mcpConfig := range config.AllConfig.AllMcpConfig {
		var trans transport.Interface
		switch mcpConfig.Type {
		case "stdio":
			trans = transport.NewStdio(mcpConfig.Command, mcpConfig.Env, mcpConfig.Args...)
		case "sse":
			trans, _ = transport.NewSSE(mcpConfig.BaseUrl)
		case "streamableHTTP":
			trans, _ = transport.NewStreamableHTTP(mcpConfig.BaseUrl)
		default:
			log.Logger.Error("unsupported mcp type", zap.String("type", mcpConfig.Type))
			continue
		}
		mcpClient := client.NewClient(trans)
		err := mcpClient.Start(context.Background())
		if err != nil {
			panic(err)
		}
		mcpClient.OnNotification(func(notification mcp.JSONRPCNotification) {})
		_, err = mcpClient.Initialize(context.Background(), mcp.InitializeRequest{})
		if err != nil {
			continue
		}
		tools, err := mcpClient.ListTools(context.Background(), mcp.ListToolsRequest{})
		AllMcpClient = append(AllMcpClient, mcpClient)
		for _, tool := range tools.Tools {
			ToolMcpMap[tool.Name] = mcpClient
			AllTools = append(AllTools, &tool)
		}
	}
	//McpClient, _ = client.NewStdioMCPClient("npx", []string{"EXA_API_KEY=ea6ab437-7e2d-4111-b2dc-710e49dba8dd"},
	//	"-y",
	//	"exa-mcp-server",
	//	"--tools=web_search_exa,research_paper_search,company_research,crawling,competitor_finder,linkedin_search,wikipedia_search_exa,github_search")
	//err := McpClient.Start(context.Background())
	//if err != nil {
	//	panic(err)
	//}
	//McpClient.OnNotification(func(notification mcp.JSONRPCNotification) {
	//	fmt.Printf("Received notification: %s\n", notification.Method)
	//})
	//result, err := McpClient.Initialize(context.Background(), mcp.InitializeRequest{})
	//tools, err := McpClient.ListTools(context.Background(), mcp.ListToolsRequest{})
	//if err != nil {
	//	return
	//}
	//fmt.Printf("mcp client init success, %s\n", util.MustJson(result))
	//fmt.Printf("mcp client init success, %s\n", util.MustJson(tools))

}

func Execute(input string, method string) (string, error) {
	callToolRequest := mcp.CallToolRequest{}
	callToolRequest.Method = method
	callToolRequest.Params.Name = method

	if ToolMcpMap[method] == nil {
		return "", fmt.Errorf("mcp client not found")
	}
	var data interface{}
	jsoniter.Unmarshal([]byte(input), &data)
	callToolRequest.Params.Arguments = data
	resp, err := ToolMcpMap[method].CallTool(context.Background(), callToolRequest)
	if err != nil {
		return "", err
	}
	return util.MustJson(resp), nil
}
