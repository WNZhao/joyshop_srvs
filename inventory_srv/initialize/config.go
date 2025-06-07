/*
 * @Author: Will nanan_zhao@163.com
 * @Date: 2025-05-12 16:57:52
 * @LastEditors: Will nanan_zhao@163.com
 * @LastEditTime: 2025-05-17 16:16:19
 * @FilePath: /joyshop_srvs/inventory_srv/initialize/config.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package initialize

import (
	"inventory_srv/config"
	"inventory_srv/global"
	"inventory_srv/util"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// InitConfig 初始化配置
func InitConfig() {
	// 尝试从 Nacos 读取配置
	if err := loadNacosConfig(); err != nil {
		zap.S().Warnf("从 Nacos 读取配置失败，将使用本地配置: %v", err)
		// loadLocalConfig()
		util.LoadLocalConfig(&global.ServerConfig)
	}
}

// loadNacosConfig 从 Nacos 加载配置
func loadNacosConfig() error {
	// 读取 Nacos 配置信息
	v := viper.New()
	v.SetConfigFile("config/nacos-dev.yaml")
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		return err
	}

	var nacosConfig config.NacosConfig
	if err := v.UnmarshalKey("nacos", &nacosConfig); err != nil {
		zap.S().Errorf("解析 Nacos 配置失败: %v", err)
		return err
	}

	// 从 Nacos 加载配置
	return util.LoadRemoteConfig(&nacosConfig, &global.ServerConfig)
}
