/*
 * @Author: Will nanan_zhao@163.com
 * @Date: 2025-05-12 16:57:44
 * @LastEditors: Will nanan_zhao@163.com
 * @LastEditTime: 2025-05-17 14:28:12
 * @FilePath: /joyshop_srvs/goods_srv/global/global.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package global

import (
	"goods_srv/config"

	"github.com/hashicorp/consul/api"
	"gorm.io/gorm"
)

var (
	DB           *gorm.DB
	ServerConfig *config.ServerConfig
	ConsulClient *api.Client
)
