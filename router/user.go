package router

import (
	"mxshop_web/api"
	"mxshop_web/middlewares"

	"github.com/gin-gonic/gin"
)

func InitUserRouter(Router *gin.RouterGroup) {
	Router.GET("/user/list", middlewares.JWTAuth(), middlewares.IsAdminAuth(), api.GetUserList)
	Router.POST("/user/passwd_login", api.PasswordLogin)
}
