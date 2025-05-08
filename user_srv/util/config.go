package util

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// GetConfigPath 获取配置文件路径
func GetConfigPath() string {
	// 获取当前工作目录
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return filepath.Join(dir, "user_srv", "config", "db", "db.yaml")
}

// LoadConfig 加载配置文件
func LoadConfig(path string, config interface{}) error {
	// 读取配置文件
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// 解析YAML
	return yaml.Unmarshal(data, config)
}
