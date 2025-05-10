/*
 * @Author: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @Date: 2025-04-30 16:33:02
 * @LastEditors: Will nanan_zhao@163.com
 * @LastEditTime: 2025-05-10 15:00:07
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

	"joyshop_srvs/user_srv/global"
	"joyshop_srvs/user_srv/handler"
	"joyshop_srvs/user_srv/initialize"
	"joyshop_srvs/user_srv/proto"
	"joyshop_srvs/user_srv/util"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	// 初始化日志
	initialize.InitLogger()

	// 初始化配置
	if err := initialize.InitConfig(); err != nil {
		zap.S().Fatalf("初始化配置失败: %v", err)
	}

	// 初始化数据库
	if err := initialize.InitDB(); err != nil {
		zap.S().Fatalf("初始化数据库失败: %v", err)
	}

	// 获取当前工作目录
	dir, err := os.Getwd()
	if err != nil {
		zap.S().Errorf("获取工作目录失败: %v", err)
	} else {
		zap.S().Infof("当前工作目录: %s", dir)
	}

	// 根据环境获取服务端口
	port, err := util.GetServerPort(global.GlobalConfig.ServerInfo.Port)
	if err != nil {
		zap.S().Fatalf("获取服务端口失败: %v", err)
	}
	// 更新全局配置中的端口
	global.GlobalConfig.ServerInfo.Port = port

	// 构建服务地址
	addr := fmt.Sprintf("%s:%d", global.GlobalConfig.ServerInfo.Host, global.GlobalConfig.ServerInfo.Port)
	zap.S().Infof("服务启动地址: %s", addr)

	// 监听端口
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		zap.S().Fatalf("监听端口失败: %v", err)
	}

	// 创建grpc服务
	s := grpc.NewServer()

	// 注册健康检查服务
	grpc_health_v1.RegisterHealthServer(s, health.NewServer())

	// 注册用户服务
	proto.RegisterUserServer(s, &handler.UserServer{})

	// 注册服务到 Consul
	consulClient, ID, err := initialize.InitConsulRegister()
	if err != nil {
		zap.S().Fatalf("注册服务到Consul失败: %v", err)
	}

	// 注册优雅退出信号监听
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		// 注销 Consul 服务
		if err := consulClient.Agent().ServiceDeregister(ID); err != nil {
			zap.S().Info("注销失败")
		} else {
			zap.S().Info("注销成功")
		}
		os.Exit(0)
	}()

	// 启动服务
	zap.S().Infof("服务启动成功，监听地址：%s", addr)
	if err := s.Serve(lis); err != nil {
		zap.S().Fatalf("服务启动失败: %v", err)
	}
}
