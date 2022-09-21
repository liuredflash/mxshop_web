package main

import (
	"fmt"
	"mxshop_web/global"
	"mxshop_web/initialize"
	myvalidator "mxshop_web/validator"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func main() {
	initialize.InitLogger()
	initialize.InitConfig()
	Router := initialize.Routers()
	//初始化srv的连接,需要再初始化config之后
	initialize.InitSrvClient()

	//注册验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("mobile", myvalidator.ValidateMobile) //给mobile这个tag(再required后面使用)限定规则
	}
	zap.S().Debugf("启动服务器，端口:%d", global.ServerConfig.Port)
	if err := Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
		zap.S().Panic("启动失败:", err.Error())
	}

}
