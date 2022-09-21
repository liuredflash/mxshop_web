package global

import (
	"mxshop_web/config"
	"mxshop_web/proto"
)

var (
	//指针类型，因为需要被动态改变
	ServerConfig  *config.ServerConfig
	UserSrvClient proto.UserClient
)
