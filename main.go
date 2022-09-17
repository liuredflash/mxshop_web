package main

import (
	"fmt"
	"mxshop_web/initialize"

	"go.uber.org/zap"
)

func main() {
	initialize.InitLogger()
	Router := initialize.Routers()
	port := 8080
	zap.S().Debugf("启动服务器，端口:%d", port)
	if err := Router.Run(fmt.Sprintf(":%d", port)); err != nil {
		zap.S().Panic("启动失败:", err.Error())
	}

}
