package main

import (
	"encoding/json"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"shop-api/user-web/nacos-test/config"
)

func main() {
	sc := []constant.ServerConfig{
		{
			IpAddr: "127.0.0.1",
			Port:   8848,
		},
	}

	cc := constant.ClientConfig{
		NamespaceId:         "e2595576-208b-4f75-b4e6-a76fbd1b45fa",
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "./tmp/nacos/log",
		CacheDir:            "./tmp/nacos/cache",
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
		DataId: "user-web.json",
		Group:  "dev"})
	if err != nil {
		panic(err)
	}
	serverConfig := config.ServerConfig{}
	json.Unmarshal([]byte(content), &serverConfig)
	fmt.Println(serverConfig)
	//err = configClient.ListenConfig(vo.ConfigParam{
	//	DataId: "user-web.json",
	//	Group:  "dev",
	//	OnChange: func(namespace, group, dataId, data string) {
	//		fmt.Println("配置文件产生变化")
	//		fmt.Println("group:" + group + ", dataId:" + dataId + ", data:" + data)
	//	},
	//})
	//
	//fmt.Println(content)
	//
	//time.Sleep(time.Second * 3000)

}
