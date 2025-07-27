/*
 * @Author: Will nanan_zhao@163.com
 * @Date: 2025-05-09 09:23:15
 * @LastEditors: Will nanan_zhao@163.com
 * @LastEditTime: 2025-05-17 16:07:56
 * @FilePath: /joyshop_srvs/user_srv/initialize/config.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package initialize

import (
	"bytes"
	"os"
	"path/filepath"
	"user_srv/config"
	"user_srv/global"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// InitConfig 初始化配置
func InitConfig() {
	// 尝试从 Nacos 读取配置
	if err := initNacosConfig(); err != nil {
		zap.S().Warnf("从 Nacos 读取配置失败，将使用本地配置: %v", err)
		// 如果 Nacos 配置失败，使用本地配置
		loadLocalConfig()
	}
}

// initNacosConfig 初始化 Nacos 配置
func initNacosConfig() error {
	// 首先读取本地 Nacos 配置
	nacosConfig, err := loadLocalNacosConfig()
	if err != nil {
		return err
	}

	// 创建 Nacos 客户端
	zap.S().Infof("Nacos 配置信息: Host=%s, Port=%d, Namespace=%s, DataId=%s, Group=%s",
		nacosConfig.Host, nacosConfig.Port, nacosConfig.Namespace, nacosConfig.DataId, nacosConfig.Group)

	sc := []constant.ServerConfig{
		{
			IpAddr: nacosConfig.Host,
			Port:   nacosConfig.Port,
		},
	}

	cc := constant.ClientConfig{
		NamespaceId:         nacosConfig.Namespace,
		TimeoutMs:           nacosConfig.Timeout,
		NotLoadCacheAtStart: true,
		LogDir:              nacosConfig.LogDir,
		CacheDir:            nacosConfig.CacheDir,
		LogLevel:            nacosConfig.LogLevel,
	}

	// 如果配置了用户名和密码，则设置认证信息
	if nacosConfig.Username != "" && nacosConfig.Password != "" {
		cc.Username = nacosConfig.Username
		cc.Password = nacosConfig.Password
		zap.S().Infof("使用 Nacos 认证: Username=%s", nacosConfig.Username)
	} else {
		zap.S().Info("未配置 Nacos 认证信息，使用匿名访问")
	}

	zap.S().Infof("ServerConfig: %+v", sc)
	zap.S().Infof("ClientConfig: %+v", cc)

	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		return err
	}

	// 从 Nacos 获取配置
	content, err := client.GetConfig(vo.ConfigParam{
		DataId: nacosConfig.DataId,
		Group:  nacosConfig.Group,
	})
	if err != nil || content == "" {
		return err
	}

	// 使用 viper 解析配置内容
	v := viper.New()
	v.SetConfigType("yaml")
	if err := v.ReadConfig(bytes.NewReader([]byte(content))); err != nil {
		return err
	}

	// 解析配置到全局变量
	if err := v.Unmarshal(&global.ServerConfig); err != nil {
		return err
	}

	zap.S().Infof("成功从 Nacos 获取配置: %+v", global.ServerConfig)

	// 监听配置变化
	err = client.ListenConfig(vo.ConfigParam{
		DataId: nacosConfig.DataId,
		Group:  nacosConfig.Group,
		OnChange: func(namespace, group, dataId, data string) {
			zap.S().Infof("Nacos 配置发生变化: namespace=%s, group=%s, dataId=%s", namespace, group, dataId)

			// 使用 viper 解析新的配置内容
			v := viper.New()
			v.SetConfigType("yaml")
			if err := v.ReadConfig(bytes.NewReader([]byte(data))); err != nil {
				zap.S().Errorf("读取新的配置内容失败: %v", err)
				return
			}

			if err := v.Unmarshal(&global.ServerConfig); err != nil {
				zap.S().Errorf("解析新的配置内容失败: %v", err)
				return
			}
			zap.S().Infof("成功更新配置: %+v", global.ServerConfig)
		},
	})
	if err != nil {
		zap.S().Errorf("设置 Nacos 配置监听失败: %v", err)
		return err
	}

	return nil
}

// loadLocalNacosConfig 加载本地 Nacos 配置
func loadLocalNacosConfig() (*config.NacosConfig, error) {
	configName := "nacos-dev.yaml"
	if os.Getenv("APP_ENV") == "production" {
		configName = "nacos-prod.yaml"
	}

	configPath := filepath.Join("config", configName)
	zap.S().Infof("正在加载 Nacos 配置文件: %s", configPath)

	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		zap.S().Errorf("读取 Nacos 配置文件失败: %v", err)
		return nil, err
	}

	var nacosConfig config.NacosConfig
	if err := v.UnmarshalKey("nacos", &nacosConfig); err != nil {
		zap.S().Errorf("解析 Nacos 配置失败: %v", err)
		return nil, err
	}

	zap.S().Infof("成功解析 Nacos 配置: %+v", nacosConfig)

	// 根据环境变量设置 group
	if os.Getenv("APP_ENV") == "production" {
		nacosConfig.Group = "prod"
	} else {
		nacosConfig.Group = "dev"
	}

	return &nacosConfig, nil
}

// loadLocalConfig 加载本地配置
func loadLocalConfig() {
	configName := "config-develop.yaml"
	if os.Getenv("APP_ENV") == "production" {
		configName = "config-prod.yaml"
	}

	v := viper.New()
	v.SetConfigFile(filepath.Join("config", configName))
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		zap.S().Fatalf("读取本地配置文件失败: %v", err)
	}

	if err := v.Unmarshal(&global.ServerConfig); err != nil {
		zap.S().Fatalf("解析本地配置失败: %v", err)
	}

	zap.S().Infof("成功加载本地配置: %+v", global.ServerConfig)
}
