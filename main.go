package main

import (
	"go.uber.org/zap"
)

func main() {
	logger := zap.NewExample()
	defer logger.Sync()
	logger.Info("/go-manus start/\n")
}
