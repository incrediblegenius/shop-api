package initialize

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"shop-api/user-web/global"
	"shop-api/user-web/proto"
)

func InitSrvConn() {
	conn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port, global.ServerConfig.UserSrcConfig.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("Failed to connect to consul", zap.Error(err))
		return
	}
	global.UserSrvClient = proto.NewUserClient(conn)
}

func InitSrvConn2() {
	// 从注册中心获取用户服务的信息
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	userSrvHost := ""
	userSrvPort := 0
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	data, err := client.Agent().ServicesWithFilter(fmt.Sprintf(`Service == "%s"`, global.ServerConfig.UserSrcConfig.Name))
	if err != nil {
		panic(err)
	}
	for _, v := range data {
		userSrvHost = v.Address
		userSrvPort = v.Port
		break
		//fmt.Println(k)
	}
	if userSrvHost == "" || userSrvPort == 0 {
		zap.S().Fatal("[InitSrvConn] 服务主机域名解析错误")
		return
	}
	//fmt.Println(userSrvHost, userSrvPort)
	//ip := "127.0.0.1"
	//port := 8088
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", userSrvHost, userSrvPort), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList] 连接 【用户服务失败】",
			"msg", err.Error())
	}
	// 1.后续用户服务下线了， 2。改端口了， 3。 ip改了（负载均衡）
	// 事先创建好了连接，不用进行多次三次握手
	// 多个连接影响性能（连接池、负载均衡）
	global.UserSrvClient = proto.NewUserClient(userConn)
}
