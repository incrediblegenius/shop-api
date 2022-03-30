package main

import (
	"context"
	"log"
	"shop-api/user-web/proto"

	_ "github.com/mbobakov/grpc-consul-resolver" // It's important

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial(
		"consul://127.0.0.1:8500/user-srv?wait=14s&tag=user",
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	userSrvClient := proto.NewUserClient(conn)
	rsp, err := userSrvClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    1,
		PSize: 2,
	})

	if err != nil {
		panic(err)
	}
	for index, data := range rsp.Data {
		log.Printf("%d: %v", index, data)
	}
}
