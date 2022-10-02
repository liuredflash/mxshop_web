package utils

import (
	"net"

	"github.com/spf13/viper"
)

func GetFreePort() (int, error) {
	// 获取可用的tcp端口
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}

	defer l.Close()

	return l.Addr().(*net.TCPAddr).Port, nil

}

func GetEnvInfo(env string) bool {
	// 判定当前服务器环境
	viper.AutomaticEnv()
	return viper.GetBool(env)
}
