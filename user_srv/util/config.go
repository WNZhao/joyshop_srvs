/*
 * @Author: Will nanan_zhao@163.com
 * @Date: 2025-05-08 17:09:35
 * @LastEditors: Will nanan_zhao@163.com
 * @LastEditTime: 2025-05-11 11:21:45
 * @FilePath: /joyshop_srvs/user_srv/util/config.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
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
