/*
 * @Author: Will nanan_zhao@163.com
 * @Date: 2025-05-12 17:12:12
 * @LastEditors: Will nanan_zhao@163.com
 * @LastEditTime: 2025-05-18 13:43:12
 * @FilePath: /joyshop_srvs/goods_srv/initialize/consul.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package initialize

import (
	"fmt"
	"goods_srv/global"

	"go.uber.org/zap"

	"github.com/hashicorp/consul/api"
)

func InitConsul() *api.Client {
	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%d", global.ServerConfig.Consul.Host, global.ServerConfig.Consul.Port)
	zap.S().Infof("Consul 地址: %s", config.Address)

	client, err := api.NewClient(config)
	if err != nil {
		zap.S().Errorf("连接 Consul 失败: %v", err)
		panic(fmt.Sprintf("连接 Consul 失败: %v", err))
	}

	zap.S().Info("Consul 客户端初始化成功")
	return client
}
