package util

import (
	"fmt"
	"user_srv/global"

	"github.com/hashicorp/consul/api"
	"github.com/satori/uuid"
	"go.uber.org/zap"
)

var (
	// ServiceID 服务ID
	ServiceID string
)

// RegisterService 注册服务到 Consul
func RegisterService() error {
	// 生成唯一服务ID
	ServiceID = uuid.NewV4().String()
	registration := &api.AgentServiceRegistration{
		ID:      ServiceID,
		Name:    global.ServerConfig.ServerInfo.Name,
		Tags:    []string{"user-srv", "user", "srv"},
		Port:    global.ServerConfig.ServerInfo.Port,
		Address: global.ServerConfig.ServerInfo.Host,
		Check: &api.AgentServiceCheck{
			GRPC:                           fmt.Sprintf("%s:%d", global.ServerConfig.ServerInfo.Host, global.ServerConfig.ServerInfo.Port),
			Timeout:                        "5s",
			Interval:                       "5s",
			DeregisterCriticalServiceAfter: "10s",
		},
	}

	if err := global.ConsulClient.Agent().ServiceRegister(registration); err != nil {
		return fmt.Errorf("注册服务失败: %v", err)
	}

	zap.S().Infof("服务注册成功，ID: %s", ServiceID)
	return nil
}

// DeregisterService 从 Consul 注销服务
func DeregisterService() error {
	if ServiceID == "" {
		return nil
	}

	if err := global.ConsulClient.Agent().ServiceDeregister(ServiceID); err != nil {
		return fmt.Errorf("注销服务失败: %v", err)
	}

	zap.S().Infof("服务注销成功，ID: %s", ServiceID)
	return nil
}
