/*
 * @Author: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @Date: 2025-05-03 14:17:12
 * @LastEditors: Will nanan_zhao@163.com
 * @LastEditTime: 2025-05-09 09:26:27
 * @FilePath: /joyshop_srvs/user_srv/global/global.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package global

import (
	"user_srv/config"

	"github.com/hashicorp/consul/api"
	"gorm.io/gorm"
)

var (
	// DB 全局数据库连接
	DB *gorm.DB
	// ServerConfig 全局配置变量
	ServerConfig config.ServerConfig
	// ConsulClient Consul 客户端
	ConsulClient *api.Client
)
