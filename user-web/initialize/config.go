package initialize

import (
	"github.com/fsnotify/fsnotify"
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
	if err := v.Unmarshal(global.ServerConfig); err != nil {
		panic(err)
	}
	//fmt.Println(global.ServerConfig)
	zap.S().Infof("配置信息： %v", global.ServerConfig)
	//fmt.Println(v.Get("name"))

	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		//fmt.Println("config file changed :", e.Name)
		zap.S().Infof("配置文件产生变化：%s", e.Name)
		_ = v.ReadInConfig()
		_ = v.Unmarshal(global.ServerConfig)
		//fmt.Println(serverConfig)
	})

	//time.Sleep(time.Second * 30)
}
