package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetUserList(*gin.Context) {
	zap.S().Debugf("获取用户列表")
}
