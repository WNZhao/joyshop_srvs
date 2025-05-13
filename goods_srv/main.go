package main

import (
	"fmt"
	"goods_srv/global"
	"goods_srv/handler"
	"goods_srv/initialize"
	"goods_srv/proto"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	// 初始化配置
	initialize.InitConfig()

	// 初始化日志
	initialize.InitLogger()

	// 初始化数据库
	initialize.InitDB()

	// 初始化 Consul
	consulClient := initialize.InitConsul()
	defer func() {
		if err := consulClient.Agent().ServiceDeregister(global.ServerConfig.Name); err != nil {
			panic(fmt.Sprintf("注销服务失败: %v", err))
		}
	}()

	// 创建 gRPC 服务器
	server := grpc.NewServer()

	// 注册商品服务
	proto.RegisterGoodsServer(server, &handler.GoodsServer{})

	// 注册健康检查服务
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(server, healthServer)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	// 启动服务
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", global.ServerConfig.Host, global.ServerConfig.Port))
	if err != nil {
		panic(fmt.Sprintf("启动服务失败: %v", err))
	}

	if err := server.Serve(lis); err != nil {
		panic(fmt.Sprintf("服务运行失败: %v", err))
	}
}
