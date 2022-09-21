package initialize

import (
	"fmt"
	"mxshop_web/global"
	"mxshop_web/proto"

	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func InitSrvClient() {
	cfg := api.DefaultConfig()
	consulInfo := global.ServerConfig.ConsulInfo
	cfg.Address = fmt.Sprintf("%s:%d", consulInfo.Host, consulInfo.Port)
	//需要连接的srv的地址
	var userSrvHost string
	var userSrvPort int
	//连接consul
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	//使用server name请求它的地址
	data, err := client.Agent().ServicesWithFilter(fmt.Sprintf("Service == \"%s\"", global.ServerConfig.UserSrvInfo.Name))
	if err != nil {
		panic(err)
	}
	for _, value := range data {
		userSrvHost = value.Address
		userSrvPort = value.Port
		break
	}
	if userSrvHost == "" {
		zap.S().Fatal("[InitServConn] 连接 【用户服务失败】")
	}
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", userSrvHost,
		userSrvPort), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList] 连接 【用户服务失败】",
			"msg", err.Error(),
		)
	}
	// 获取全局的client
	// 后续如何跟着配置进行动态变化
	//初始化建立好了连接，后续不再需要tcp的三次握手
	//一个连接多个goroutine公用，会存在性能问题，使用连接池解决这个问题，再github上自己学习了解源码（grpc-go-pool）
	global.UserSrvClient = proto.NewUserClient(userConn)
}
