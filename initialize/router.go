package initialize

import (
	"mxshop_web/router"

	"github.com/gin-gonic/gin"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	ApiGroup := Router.Group("/v1")
	router.InitUserRouter(ApiGroup)
	return Router
}
