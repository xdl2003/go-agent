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
	prompt := "使用工具生成关于广告系统的知识图谱，你有create_entities,create_relations, add_observations等工具来构建知识图谱，你已经完成了一些调研在workspace中"
	flow := flow.NewFlow()
	_, err := flow.Execute(&prompt)
	if err != nil {
		return
	}
	// fmt.Println(util.MustJson(resp))
}
