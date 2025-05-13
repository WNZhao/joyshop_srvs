package util

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func LoadConfig(config interface{}) {
	viper.SetConfigName("config-develop")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		zap.S().Fatalf("读取配置文件失败: %v", err)
	}

	if err := viper.Unmarshal(config); err != nil {
		zap.S().Fatalf("解析配置文件失败: %v", err)
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		zap.S().Infof("配置文件发生变化: %s", e.Name)
		if err := viper.Unmarshal(config); err != nil {
			zap.S().Fatalf("重新解析配置文件失败: %v", err)
		}
	})
}
