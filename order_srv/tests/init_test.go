/*
 * @Author: Will nanan_zhao@163.com
 * @Date: 2025-05-18 14:06:45
 * @LastEditors: Please set LastEditors
 * @LastEditTime: 2025-06-07 21:36:14
 * @FilePath: /joyshop_srvs/order_srv/tests/init_test.go
 * @Description: 测试初始化文件
 */
package tests

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"order_srv/config"
	"order_srv/global"
	"order_srv/initialize"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// initTestEnv 初始化测试环境
func initTestEnv(t *testing.T) {
	// 初始化日志
	zap.S().Info("开始初始化测试环境...")

	// 获取当前工作目录
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("获取工作目录失败: %v", err)
	}

	// 设置配置文件路径
	configFile := filepath.Join(dir, "..", "config", "config-develop.yaml")
	zap.S().Infof("正在加载配置文件: %s", configFile)

	// 检查配置文件是否存在
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		t.Fatalf("配置文件不存在: %s", configFile)
	}

	// 初始化viper
	v := viper.New()
	v.SetConfigFile(configFile)

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		t.Fatalf("读取配置文件失败: %v", err)
	}

	// 解析配置到结构体
	var cfg config.ServerConfig
	if err := v.Unmarshal(&cfg); err != nil {
		t.Fatalf("解析配置文件失败: %v", err)
	}

	// 设置全局配置
	global.ServerConfig = &cfg

	// 构建DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.MySQL.User,
		cfg.MySQL.Password,
		cfg.MySQL.Host,
		cfg.MySQL.Port,
		cfg.MySQL.DBName,
	)

	// 配置GORM
	gormConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
		Logger: logger.Default.LogMode(logger.Info),
	}

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		t.Fatalf("连接数据库失败: %v", err)
	}

	// 设置全局DB实例
	global.DB = db

	// 初始化Redis
	initialize.InitRedis()
}
