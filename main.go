package main

import (
	"go-manus/go-manus/config"
	"go-manus/go-manus/flow"
	"go-manus/go-manus/log"
	"go-manus/go-manus/mcp"
)

func main() {
	log.InitLogger()
	defer func() {
		log.CloseLogger()
	}()
	log.Logger.Info("/go-manus start/\n")
	config.InitConfig()
	mcp.InitMcp()
	prompt := "你是一个AI产品经理，请生成一篇对于devin工具的调研报告"
	flow := flow.NewFlow()
	_, err := flow.Execute(&prompt)
	if err != nil {
		return
	}
	// fmt.Println(util.MustJson(resp))
}
