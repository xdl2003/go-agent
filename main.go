package main

import (
	"go-manus/go-manus/config"
	"go-manus/go-manus/flow"
	"go-manus/go-manus/log"
)

func main() {
	log.InitLogger()
	defer func() {
		log.CloseLogger()
	}()
	log.Logger.Info("/go-manus start/\n")
	config.InitConfig()
	prompt := "请生成一篇对于manus的调研报告，包括历史，使用体验等"
	flow := flow.NewFlow()
	_, err := flow.Execute(&prompt)
	if err != nil {
		return
	}
	// fmt.Println(util.MustJson(resp))
}
