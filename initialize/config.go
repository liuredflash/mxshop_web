package initialize

import (
	"fmt"
	"mxshop_web/global"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}
func InitConfig() {
	viper.SetConfigType("yaml")
	debug := GetEnvInfo("MXSHOP_DEBUG")
	configFilePrefix := "config"
	var configFileName string
	if debug {
		configFileName = fmt.Sprintf("%s-debug.yaml", configFilePrefix)
		zap.S().Infof("当前是开发环境")
	} else {
		configFileName = fmt.Sprintf("%s-pro.yaml", configFilePrefix)
		zap.S().Infof("当前是生产环境")
	}
	v := viper.New()
	zap.S().Debugf("当前配置文件：%s", configFileName)
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	//这个对象如何在其它文件中使用，使用全局变量
	// serverConfig := config.ServerConfig{}
	if err := v.Unmarshal(&global.ServerConfig); err != nil {
		panic(err)
	}
	zap.S().Infof("配置信息: %v", global.ServerConfig)

	//动态监控变化
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed: ", e.Name)
		_ = v.ReadInConfig()
		_ = v.Unmarshal(&global.ServerConfig)
		fmt.Println(global.ServerConfig)

	})

}
