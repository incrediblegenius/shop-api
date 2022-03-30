package initialize

import (
	"encoding/json"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"shop-api/user-web/global"
)

func GetEnvInfo(env string) string {
	viper.AutomaticEnv()
	return viper.GetString(env)

}

func InitConfig() {
	env := GetEnvInfo("GOPATH")
	var configFileName string = "user-web/config-debug.yaml"

	if env != "" {
		configFileName = "user-web/config-debug.yaml"
	} else {
		configFileName = "user-web/config-pro.yaml"
	}
	//serverConfig := config.ServerConfig{}
	v := viper.New()
	//v.SetConfigFile("user-web/viper-test/config-debug.yaml")
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := v.Unmarshal(global.NacosConfig); err != nil {
		panic(err)
	}
	//fmt.Println(global.ServerConfig)
	zap.S().Infof("配置信息： %v", global.NacosConfig)
	//fmt.Println(v.Get("name"))

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
		LogDir:              "user-web/tmp/nacos/log",
		CacheDir:            "user-web/tmp/nacos/cache",
		LogLevel:            "debug",
	}

	configClient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		panic(err)
	}

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: global.NacosConfig.DataId,
		Group:  global.NacosConfig.Group})
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal([]byte(content), global.ServerConfig)
	if err != nil {
		zap.S().Fatalf("读取nacos失败：%s", err.Error())
	}
	fmt.Println(global.ServerConfig)
	//v.WatchConfig()
	//v.OnConfigChange(func(e fsnotify.Event) {
	//	//fmt.Println("config file changed :", e.Name)
	//	zap.S().Infof("配置文件产生变化：%s", e.Name)
	//	_ = v.ReadInConfig()
	//	_ = v.Unmarshal(global.ServerConfig)
	//	//fmt.Println(serverConfig)
	//})

	//time.Sleep(time.Second * 30)
}
