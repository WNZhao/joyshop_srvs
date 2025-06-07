package util

import (
	"bytes"
	"inventory_srv/config"

	"github.com/fsnotify/fsnotify"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// LoadLocalConfig 从本地文件加载配置
func LoadLocalConfig(config interface{}) {
	zap.S().Info("开始从本地加载配置...")
	v := viper.New()
	v.SetConfigName("config-develop")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")

	if err := v.ReadInConfig(); err != nil {
		zap.S().Errorf("读取本地配置文件失败: %v", err)
		zap.S().Fatal("无法继续运行，程序退出")
	}

	if err := v.Unmarshal(config); err != nil {
		zap.S().Errorf("解析本地配置失败: %v", err)
		zap.S().Fatal("无法继续运行，程序退出")
	}

	zap.S().Info("成功加载本地配置")

	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		zap.S().Infof("检测到本地配置文件变化: %s", e.Name)
		if err := v.Unmarshal(config); err != nil {
			zap.S().Errorf("重新解析本地配置失败: %v", err)
			return
		}
		zap.S().Info("成功重新加载本地配置")
	})
}

// LoadRemoteConfig 从 Nacos 加载配置
func LoadRemoteConfig(nacosConfig *config.NacosConfig, targetConfig interface{}) error {
	zap.S().Info("开始从 Nacos 加载配置...")

	// 创建 Nacos 客户端
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

	zap.S().Infof("正在连接 Nacos 服务器: %s:%d", sc[0].IpAddr, sc[0].Port)
	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		zap.S().Errorf("创建 Nacos 客户端失败: %v", err)
		return err
	}

	// 从 Nacos 获取配置
	zap.S().Infof("正在从 Nacos 获取配置: DataId=%s, Group=%s", nacosConfig.DataId, nacosConfig.Group)

	content, err := client.GetConfig(vo.ConfigParam{
		DataId: nacosConfig.DataId,
		Group:  nacosConfig.Group,
	})
	if err != nil || content == "" {
		zap.S().Errorf("从 Nacos 获取配置失败: %v", err)
		return err
	}

	// 使用 viper 解析配置内容
	v := viper.New()
	v.SetConfigType("yaml")
	if err := v.ReadConfig(bytes.NewReader([]byte(content))); err != nil {
		zap.S().Errorf("解析 Nacos 配置内容失败: %v", err)
		return err
	}

	// 解析配置到目标结构体
	if err := v.Unmarshal(targetConfig); err != nil {
		zap.S().Errorf("将 Nacos 配置解析到结构体失败: %v", err)
		return err
	}

	zap.S().Info("成功从 Nacos 加载配置")

	// 监听配置变化
	err = client.ListenConfig(vo.ConfigParam{
		DataId: nacosConfig.DataId,
		Group:  nacosConfig.Group,
		OnChange: func(namespace, group, dataId, data string) {
			zap.S().Infof("检测到 Nacos 配置变化: namespace=%s, group=%s, dataId=%s", namespace, group, dataId)

			v := viper.New()
			v.SetConfigType("yaml")
			if err := v.ReadConfig(bytes.NewReader([]byte(data))); err != nil {
				zap.S().Errorf("读取新的 Nacos 配置内容失败: %v", err)
				return
			}

			if err := v.Unmarshal(targetConfig); err != nil {
				zap.S().Errorf("解析新的 Nacos 配置内容失败: %v", err)
				return
			}
			zap.S().Info("成功更新 Nacos 配置")
		},
	})

	if err != nil {
		zap.S().Errorf("设置 Nacos 配置监听失败: %v", err)
		return err
	}

	zap.S().Info("成功设置 Nacos 配置监听")
	return nil
}
