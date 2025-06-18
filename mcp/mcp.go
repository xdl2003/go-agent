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
	"github.com/mark3labs/mcp-go/mcp"
	"go-manus/go-manus/util"
)

var (
	McpClient *client.Client
)

func InitMcp() {
	// sse, err := transport.NewSSE("http://appbuilder.baidu.com/v2/ai_search/mcp/sse?api_key=Bearer+bce-v3/ALTAK-vWqpUgUOM3oeYcxo2HM22/b23269e5f22b0f4f2f9b929d48f7068327748bc3")
	//if err != nil {
	//	panic(err)
	//}
	//McpClient = client.NewClient(sse)
	McpClient, _ = client.NewStdioMCPClient("npx", []string{"EXA_API_KEY=ea6ab437-7e2d-4111-b2dc-710e49dba8dd"},
		"-y",
		"exa-mcp-server",
		"--tools=web_search_exa,research_paper_search,company_research,crawling,competitor_finder,linkedin_search,wikipedia_search_exa,github_search")
	err := McpClient.Start(context.Background())
	if err != nil {
		panic(err)
	}
	McpClient.OnNotification(func(notification mcp.JSONRPCNotification) {
		fmt.Printf("Received notification: %s\n", notification.Method)
	})
	result, err := McpClient.Initialize(context.Background(), mcp.InitializeRequest{})
	tools, err := McpClient.ListTools(context.Background(), mcp.ListToolsRequest{})
	if err != nil {
		return
	}
	fmt.Printf("mcp client init success, %s\n", util.MustJson(result))
	fmt.Printf("mcp client init success, %s\n", util.MustJson(tools))
	//McpClient.CallTool(context.Background(), mcp.CallToolRequest{
	//
	//})
}

func Execute(input string, method string) (string, error) {
	callToolRequest := mcp.CallToolRequest{}
	callToolRequest.Method = method
	callToolRequest.Params.Name = method
	var data interface{}
	jsoniter.Unmarshal([]byte(input), &data)
	callToolRequest.Params.Arguments = data
	resp, err := McpClient.CallTool(context.Background(), callToolRequest)
	if err != nil {
		return "", err
	}
	return util.MustJson(resp), nil
}
