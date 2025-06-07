/*
 * @Author: Will nanan_zhao@163.com
 * @Date: 2025-05-12 17:13:18
 * @LastEditors: Will nanan_zhao@163.com
 * @LastEditTime: 2025-05-18 11:00:52
 * @FilePath: /joyshop_srvs/inventory_srv/main.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"inventory_srv/global"
	"inventory_srv/handler"
	"inventory_srv/initialize"
	"inventory_srv/proto"
	"inventory_srv/util"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	// 初始化日志
	initialize.InitLogger()

	// 初始化配置
	initialize.InitConfig()

	// 初始化数据库
	initialize.InitDB()

	// 初始化 Redis
	initialize.InitRedis()

	// 初始化 Consul
	global.ConsulClient = initialize.InitConsul()

	// 创建 gRPC 服务器
	server := grpc.NewServer()

	// 注册健康检查服务
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(server, healthServer)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	// 注册库存服务
	proto.RegisterInventoryServiceServer(server, &handler.InventoryServer{})

	// 获取服务地址和端口
	addr := fmt.Sprintf("%s:%d", global.ServerConfig.Host, global.ServerConfig.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		zap.S().Fatalf("监听端口失败: %v", err)
	}

	// 启动服务
	go func() {
		zap.S().Infof("库存服务启动成功，监听地址: %s", addr)
		if err := server.Serve(lis); err != nil {
			zap.S().Fatalf("启动服务失败: %v", err)
		}
	}()

	// 注册服务到 Consul
	if err := util.RegisterService(); err != nil {
		zap.S().Fatalf("注册服务到 Consul 失败: %v", err)
	}

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// 注销服务
	if err := util.DeregisterService(); err != nil {
		zap.S().Errorf("从 Consul 注销服务失败: %v", err)
	}

	// 优雅关闭
	server.GracefulStop()
	zap.S().Info("库存服务已关闭")
}
