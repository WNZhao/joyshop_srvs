package main

import (
	"crypto/sha512"
	"fmt"
	"os"
	"path/filepath"

	"joyshop_srvs/user_srv/config"
	"joyshop_srvs/user_srv/model"

	"github.com/anaskhan96/go-password-encoder"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func main() {
	// 初始化日志
	zap.S().Info("开始初始化测试数据...")

	// 获取当前工作目录
	dir, err := os.Getwd()
	if err != nil {
		zap.S().Fatalf("获取工作目录失败: %v", err)
	}

	// 设置配置文件路径
	configFile := filepath.Join(dir, "user_srv", "config", "config-develop.yaml")
	zap.S().Infof("正在加载配置文件: %s", configFile)

	// 检查配置文件是否存在
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		zap.S().Fatalf("配置文件不存在: %s", configFile)
	}

	// 初始化viper
	v := viper.New()
	v.SetConfigFile(configFile)

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		zap.S().Fatalf("读取配置文件失败: %v", err)
	}

	// 解析配置到结构体
	var cfg config.ServerConfig
	if err := v.Unmarshal(&cfg); err != nil {
		zap.S().Fatalf("解析配置文件失败: %v", err)
	}

	// 构建DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBConfig.User,
		cfg.DBConfig.Password,
		cfg.DBConfig.Host,
		cfg.DBConfig.Port,
		cfg.DBConfig.DBName,
	)

	zap.S().Debugf("数据库连接信息: %s", dsn)

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
		zap.S().Fatalf("连接数据库失败: %v", err)
	}

	zap.S().Info("数据库连接成功")

	// 自动迁移表结构
	if err := db.AutoMigrate(&model.User{}); err != nil {
		zap.S().Fatalf("自动迁移表结构失败: %v", err)
	}
	zap.S().Info("表结构迁移成功")

	// 创建测试用户
	options := &password.Options{10, 100, 32, sha512.New}
	salt, encodedPwd := password.Encode("admin123", options)
	newPassword := fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)

	zap.S().Info("开始创建测试用户...")
	for i := 0; i < 10; i++ {
		mobile := fmt.Sprintf("1380013800%d", i)

		// 检查用户是否已存在
		var count int64
		if err := db.Model(&model.User{}).Where("mobile = ?", mobile).Count(&count).Error; err != nil {
			zap.S().Errorf("检查用户是否存在失败: %v", err)
			continue
		}

		if count > 0 {
			zap.S().Infof("用户已存在，跳过创建: %s", mobile)
			continue
		}

		user := &model.User{
			Mobile:   mobile,
			NickName: fmt.Sprintf("小明%d", i),
			UserName: fmt.Sprintf("xiaoming%d", i),
			Password: newPassword,
			Email:    fmt.Sprintf("xiaoming%d@123.com", i),
		}
		if err := db.Save(user).Error; err != nil {
			zap.S().Errorf("创建用户失败: %v", err)
			continue
		}
		zap.S().Infof("创建用户成功: %s", user.UserName)
	}

	zap.S().Info("测试数据初始化完成")
}
