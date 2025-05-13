/*
 * @Author: Will nanan_zhao@163.com
 * @Date: 2025-05-12 17:12:12
 * @LastEditors: Will nanan_zhao@163.com
 * @LastEditTime: 2025-05-12 17:14:33
 * @FilePath: /joyshop_srvs/goods_srv/initialize/consul.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package initialize

import (
	"fmt"
	"goods_srv/global"

	"github.com/hashicorp/consul/api"
)

func InitConsul() *api.Client {
	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%d", global.ServerConfig.Consul.Host, global.ServerConfig.Consul.Port)

	client, err := api.NewClient(config)
	if err != nil {
		panic(fmt.Sprintf("连接 Consul 失败: %v", err))
	}

	registration := new(api.AgentServiceRegistration)
	registration.Name = global.ServerConfig.Name
	registration.ID = global.ServerConfig.Name
	registration.Port = global.ServerConfig.Port
	registration.Tags = global.ServerConfig.Tags
	registration.Address = global.ServerConfig.Host

	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", global.ServerConfig.Host, global.ServerConfig.Port),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "10s",
	}
	registration.Check = check

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic(fmt.Sprintf("注册服务失败: %v", err))
	}

	return client
}
