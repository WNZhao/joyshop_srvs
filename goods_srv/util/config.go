package util

import (
	"bytes"

	"github.com/fsnotify/fsnotify"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// LoadLocalConfig 从本地文件加载配置
func LoadLocalConfig(config interface{}) {
	v := viper.New()
	v.SetConfigName("config-develop")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")

	if err := v.ReadInConfig(); err != nil {
		zap.S().Fatalf("读取本地配置文件失败: %v", err)
	}

	if err := v.Unmarshal(config); err != nil {
		zap.S().Fatalf("解析本地配置失败: %v", err)
	}

	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		zap.S().Infof("本地配置文件发生变化: %s", e.Name)
		if err := v.Unmarshal(config); err != nil {
			zap.S().Fatalf("重新解析本地配置失败: %v", err)
		}
	})
}

// LoadRemoteConfig 从 Nacos 加载配置
func LoadRemoteConfig(nacosConfig interface{}, targetConfig interface{}) error {
	// 创建 Nacos 客户端
	sc := []constant.ServerConfig{
		{
			IpAddr: nacosConfig.(struct{ Host string }).Host,
			Port:   nacosConfig.(struct{ Port uint64 }).Port,
		},
	}

	cc := constant.ClientConfig{
		NamespaceId:         nacosConfig.(struct{ Namespace string }).Namespace,
		TimeoutMs:           nacosConfig.(struct{ Timeout uint64 }).Timeout,
		NotLoadCacheAtStart: true,
		LogDir:              nacosConfig.(struct{ LogDir string }).LogDir,
		CacheDir:            nacosConfig.(struct{ CacheDir string }).CacheDir,
		LogLevel:            nacosConfig.(struct{ LogLevel string }).LogLevel,
	}

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
		DataId: nacosConfig.(struct{ DataId string }).DataId,
		Group:  nacosConfig.(struct{ Group string }).Group,
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

	// 解析配置到目标结构体
	if err := v.Unmarshal(targetConfig); err != nil {
		return err
	}

	// 监听配置变化
	err = client.ListenConfig(vo.ConfigParam{
		DataId: nacosConfig.(struct{ DataId string }).DataId,
		Group:  nacosConfig.(struct{ Group string }).Group,
		OnChange: func(namespace, group, dataId, data string) {
			zap.S().Infof("Nacos 配置发生变化: namespace=%s, group=%s, dataId=%s", namespace, group, dataId)

			v := viper.New()
			v.SetConfigType("yaml")
			if err := v.ReadConfig(bytes.NewReader([]byte(data))); err != nil {
				zap.S().Errorf("读取新的配置内容失败: %v", err)
				return
			}

			if err := v.Unmarshal(targetConfig); err != nil {
				zap.S().Errorf("解析新的配置内容失败: %v", err)
				return
			}
			zap.S().Infof("成功更新配置: %+v", targetConfig)
		},
	})

	return err
}
