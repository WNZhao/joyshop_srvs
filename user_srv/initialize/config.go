/*
 * @Author: Will nanan_zhao@163.com
 * @Date: 2025-05-09 09:23:15
 * @LastEditors: Will nanan_zhao@163.com
 * @LastEditTime: 2025-05-09 10:31:15
 * @FilePath: /joyshop_srvs/user_srv/initialize/config.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package initialize

import (
	"fmt"
	"os"
	"path/filepath"

	"joyshop_srvs/user_srv/global"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// InitConfig 初始化配置
func InitConfig() error {
	// 获取当前工作目录
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("获取工作目录失败: %v", err)
	}

	// 设置配置文件路径
	configFile := filepath.Join(dir, "user_srv", "config", "config-develop.yaml")
	zap.S().Infof("正在加载配置文件: %s", configFile)

	// 检查配置文件是否存在
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return fmt.Errorf("配置文件不存在: %s", configFile)
	}

	// 初始化viper
	v := viper.New()
	v.SetConfigFile(configFile)

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 解析配置到结构体
	if err := v.Unmarshal(&global.GlobalConfig); err != nil {
		return fmt.Errorf("解析配置文件失败: %v", err)
	}

	zap.S().Info("配置加载成功")
	return nil
}
