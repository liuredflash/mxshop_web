package global

import "mxshop_web/config"

var (
	//指针类型，因为需要被动态改变
	ServerConfig *config.ServerConfig
)
