/*
 * @Author: Will nanan_zhao@163.com
 * @Date: 2025-05-17 17:18:18
 * @LastEditors: Will nanan_zhao@163.com
 * @LastEditTime: 2025-05-17 17:19:52
 * @FilePath: /joyshop_srvs/goods_srv/util/consul.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package util

import (
	"fmt"
	"goods_srv/global"

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
		Name:    global.ServerConfig.Name,
		Tags:    global.ServerConfig.Tags,
		Port:    global.ServerConfig.Port,
		Address: global.ServerConfig.Host,
		Check: &api.AgentServiceCheck{
			GRPC:                           fmt.Sprintf("%s:%d", global.ServerConfig.Host, global.ServerConfig.Port),
			Timeout:                        "5s",
			Interval:                       "5s",
			DeregisterCriticalServiceAfter: "10s",
		},
	}

	if err := global.ConsulClient.Agent().ServiceRegister(registration); err != nil {
		return fmt.Errorf("注册服务失败: %v", err)
	}

	zap.S().Infof("商品服务注册成功，ID: %s", ServiceID)
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

	zap.S().Infof("商品服务注销成功，ID: %s", ServiceID)
	return nil
}
