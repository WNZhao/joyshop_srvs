package initialize

import (
	"fmt"
	"joyshop_srvs/user_srv/global"

	"github.com/hashicorp/consul/api"
	"github.com/satori/uuid"
	"go.uber.org/zap"
)

// InitConsulRegister 初始化 Consul 服务注册
func InitConsulRegister() (client *api.Client, ID string, err error) {
	// 创建 Consul 客户端
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.GlobalConfig.ConsulInfo.Host, global.GlobalConfig.ConsulInfo.Port)
	zap.S().Infof(cfg.Address)
	client, er := api.NewClient(cfg)
	if er != nil {
		return nil, "", fmt.Errorf("创建Consul客户端失败: %v", er)
	}

	// 配置健康检查
	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", global.GlobalConfig.ConsulInfo.Host, global.GlobalConfig.ServerInfo.Port), // 这块不能用0.0.0.0 因为是同一台机器，我们暂用
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "15s",
	}
	serverId := fmt.Sprintf("%s", uuid.NewV4())
	// 配置服务注册信息
	registration := &api.AgentServiceRegistration{
		Name:    global.GlobalConfig.ServerInfo.Name,
		ID:      serverId, // , // global.GlobalConfig.ServerInfo.Name,
		Tags:    []string{"joyshop", "user_srv"},
		Port:    global.GlobalConfig.ServerInfo.Port,
		Address: global.GlobalConfig.ConsulInfo.Host,
		Check:   check,
	}

	// 注册服务
	if err := client.Agent().ServiceRegister(registration); err != nil {
		return nil, "", fmt.Errorf("服务注册失败: %v", err)
	}

	zap.S().Info("服务注册成功")
	return client, registration.ID, nil
}
