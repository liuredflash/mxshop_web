package initialize

import "go.uber.org/zap"

func InitLogger() {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger) //替换全局的zap
}
