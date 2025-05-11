/*
 * @Author: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @Date: 2025-04-30 16:33:02
 * @LastEditors: Will nanan_zhao@163.com
 * @LastEditTime: 2025-05-11 11:45:46
 * @FilePath: /joyshop_srvs/user_srv/main.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"user_srv/global"
	"user_srv/handler"
	"user_srv/initialize"
	"user_srv/proto"
	"user_srv/util"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	// 初始化日志
	initialize.InitLogger()

	// 初始化配置
	initialize.InitConfig()

	// 初始化数据库
	initialize.InitDB()

	// 初始化 Consul
	initialize.InitConsul()

	// 创建 gRPC 服务器
	server := grpc.NewServer()
	proto.RegisterUserServer(server, &handler.UserServer{})

	// 获取服务地址和端口
	addr := fmt.Sprintf("%s:%d", global.ServerConfig.ServerInfo.Host, global.ServerConfig.ServerInfo.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		zap.S().Fatalf("监听端口失败: %v", err)
	}

	// 启动服务
	go func() {
		zap.S().Infof("服务启动成功，监听地址: %s", addr)
		if err := server.Serve(lis); err != nil {
			zap.S().Fatalf("启动服务失败: %v", err)
		}
	}()

	// 注册服务到 Consul
	if err := util.RegisterService(); err != nil {
		zap.S().Fatalf("注册服务失败: %v", err)
	}

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// 注销服务
	if err := util.DeregisterService(); err != nil {
		zap.S().Errorf("注销服务失败: %v", err)
	}

	// 优雅关闭
	server.GracefulStop()
	zap.S().Info("服务已关闭")
}
