package initialize

import (
	"goods_srv/config"
	"goods_srv/global"
	"goods_srv/util"
)

func InitConfig() {
	// 初始化配置
	config := config.ServerConfig{}
	util.LoadConfig(&config)
	global.ServerConfig = config
}
