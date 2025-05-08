/*
 * @Author: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @Date: 2025-04-30 16:33:02
 * @LastEditors: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @LastEditTime: 2025-05-08 17:31:44
 * @FilePath: /joyshop_srvs/user_srv/main.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"joyshop_srvs/user_srv/handler"
	"joyshop_srvs/user_srv/initialize"
	"joyshop_srvs/user_srv/proto"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	// 初始化日志
	initialize.InitLogger()

	// 初始化数据库
	initialize.InitDB()

	// 获取并打印当前工作目录
	if dir, err := os.Getwd(); err != nil {
		zap.S().Errorf("获取工作目录失败: %v", err)
	} else {
		zap.S().Debugf("当前工作目录: %s", dir)
	}

	IP := flag.String("ip", "0.0.0.0", "ip地址")
	PORT := flag.Int("port", 50051, "端口号")

	flag.Parse()
	addr := fmt.Sprintf("%s:%d", *IP, *PORT)
	// 监听端口
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		zap.S().Fatalf("监听端口失败: %v", err)
	}

	// 创建grpc服务
	s := grpc.NewServer()
	// 注册服务
	proto.RegisterUserServer(s, &handler.UserServer{})
	zap.S().Infof("服务启动成功，监听地址：%s", addr)
	// 启动服务
	if err := s.Serve(lis); err != nil {
		zap.S().Fatalf("服务启动失败: %v", err)
	}
}
