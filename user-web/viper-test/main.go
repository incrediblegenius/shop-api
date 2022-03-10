package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"time"
)

//type ServerConfig struct {
//	Name string `mapstructure:"name"`
//	Port int    `mapstructure:"port"`
//}
type MySQLConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

// 线上线下配置隔离
//不用改代码，线上线下文件隔离开

type ServerConfig struct {
	Name string `mapstructure:"name"`
	//Port int    `mapstructure:"port"`
	MysqlInfo MySQLConfig `mapstructure:"mysql"`
}

func GetEnvInfo(env string) string {
	viper.AutomaticEnv()
	return viper.GetString(env)

}

func main() {
	env := GetEnvInfo("GOPATH")
	var configFileName string = "user-web/viper-test/config-debug.yaml"

	if env != "" {
		configFileName = "config-debug.yaml"
	} else {
		configFileName = "config-pro.yaml"
	}
	serverConfig := ServerConfig{}
	v := viper.New()
	//v.SetConfigFile("user-web/viper-test/config-debug.yaml")
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := v.Unmarshal(&serverConfig); err != nil {
		panic(err)
	}
	fmt.Println(serverConfig)
	fmt.Println(v.Get("name"))

	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed :", e.Name)
		_ = v.ReadInConfig()
		_ = v.Unmarshal(&serverConfig)
		fmt.Println(serverConfig)
	})

	time.Sleep(time.Second * 30)

}
