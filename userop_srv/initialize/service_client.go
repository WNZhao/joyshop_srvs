package initialize

import (
	"go.uber.org/zap"
)

// InitServiceClients 初始化服务客户端连接
func InitServiceClients() {
	// userop_srv 目前不需要连接其他服务
	// 如果后续需要调用其他服务，可以在这里添加
	zap.S().Info("用户操作服务客户端初始化完成")
}

// CloseServiceClients 关闭所有服务客户端连接
func CloseServiceClients() {
	// 如果后续添加了服务客户端连接，在这里关闭
	zap.S().Info("所有服务客户端连接已关闭")
}