package router

import (
	"mxshop_web/api"

	"github.com/gin-gonic/gin"
)

func InitUserRouter(Router *gin.RouterGroup) {
	Router.GET("/user/list", api.GetUserList)
}
