package initialize

import (
	"fmt"
	"user_srv/global"

	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
)

// InitConsul 初始化 Consul 客户端
func InitConsul() {
	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)

	client, err := api.NewClient(config)
	if err != nil {
		zap.S().Fatalf("创建 Consul 客户端失败: %v", err)
	}

	global.ConsulClient = client
	zap.S().Info("Consul 客户端初始化成功")
}
