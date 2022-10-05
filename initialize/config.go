package initialize

import (
	"encoding/json"
	"fmt"
	"mxshop_web/global"
	"mxshop_web/utils"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func InitConfig() {
	viper.SetConfigType("yaml") //当文件不是yaml结尾时需要显示指定，是yaml结尾时可以不要
	debug := utils.GetEnvInfo("MXSHOP_DEBUG")
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
	if err := v.Unmarshal(&global.NacosConfig); err != nil {
		panic(err)
	}
	zap.S().Infof("配置信息: %v", global.NacosConfig)

	//从nacos中读取配置
	sc := []constant.ServerConfig{
		{
			IpAddr: global.NacosConfig.Host,
			Port:   global.NacosConfig.Port,
		},
	}

	cc := constant.ClientConfig{
		NamespaceId:         global.NacosConfig.Namespace,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "tmp/nacos/log",
		CacheDir:            "tmp/nacos/cache",
		LogLevel:            "debug",
	}

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})

	if err != nil {
		panic(err)
	}
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: global.NacosConfig.DataId,
		Group:  global.NacosConfig.Group,
	})
	if err != nil {
		panic(err)
	}

	json.Unmarshal([]byte(content), &global.ServerConfig)
	fmt.Println(global.ServerConfig)
	//动态监控变化
	// v.WatchConfig()
	// v.OnConfigChange(func(e fsnotify.Event) {
	// 	fmt.Println("config file changed: ", e.Name)
	// 	_ = v.ReadInConfig()
	// 	_ = v.Unmarshal(&global.ServerConfig)
	// 	fmt.Println(global.ServerConfig)

	// })

}
