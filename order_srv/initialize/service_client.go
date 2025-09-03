package initialize

import (
	"fmt"
	"order_srv/global"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// InitServiceClients 初始化服务客户端连接
func InitServiceClients() {
	initGoodsClient()
	initInventoryClient()
}

// initGoodsClient 初始化商品服务客户端
func initGoodsClient() {
	// 从Consul获取商品服务地址
	services, _, err := global.ConsulClient.Health().Service("goods_srv", "", true, nil)
	if err != nil {
		zap.S().Errorf("从Consul获取商品服务失败: %v", err)
		return
	}
	
	if len(services) == 0 {
		zap.S().Error("没有可用的商品服务实例")
		return
	}
	
	// 简单选择第一个健康的服务实例
	service := services[0]
	addr := fmt.Sprintf("%s:%d", service.Service.Address, service.Service.Port)
	
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.S().Errorf("连接商品服务失败: %v", err)
		return
	}
	
	global.GoodsClient = conn
	zap.S().Infof("商品服务客户端连接成功: %s", addr)
}

// initInventoryClient 初始化库存服务客户端
func initInventoryClient() {
	// 从Consul获取库存服务地址
	services, _, err := global.ConsulClient.Health().Service("inventory_srv", "", true, nil)
	if err != nil {
		zap.S().Errorf("从Consul获取库存服务失败: %v", err)
		return
	}
	
	if len(services) == 0 {
		zap.S().Error("没有可用的库存服务实例")
		return
	}
	
	// 简单选择第一个健康的服务实例
	service := services[0]
	addr := fmt.Sprintf("%s:%d", service.Service.Address, service.Service.Port)
	
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.S().Errorf("连接库存服务失败: %v", err)
		return
	}
	
	global.InventoryClient = conn
	zap.S().Infof("库存服务客户端连接成功: %s", addr)
}

// CloseServiceClients 关闭所有服务客户端连接
func CloseServiceClients() {
	if global.GoodsClient != nil {
		global.GoodsClient.Close()
	}
	if global.InventoryClient != nil {
		global.InventoryClient.Close()
	}
}