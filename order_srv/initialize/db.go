/*
 * @Author: Will nanan_zhao@163.com
 * @Date: 2025-05-12 16:58:21
 * @LastEditors: Will nanan_zhao@163.com
 * @LastEditTime: 2025-05-17 14:30:18
 * @FilePath: /joyshop_srvs/order_srv/initialize/db.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package initialize

import (
	"fmt"
	"log"
	"order_srv/global"
	"os"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func InitDB() {
	zap.S().Info("开始初始化数据库连接...")
	zap.S().Infof("数据库配置信息: Host=%s, Port=%d, User=%s, DBName=%s",
		global.ServerConfig.MySQL.Host,
		global.ServerConfig.MySQL.Port,
		global.ServerConfig.MySQL.User,
		global.ServerConfig.MySQL.DBName,
	)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		global.ServerConfig.MySQL.User,
		global.ServerConfig.MySQL.Password,
		global.ServerConfig.MySQL.Host,
		global.ServerConfig.MySQL.Port,
		global.ServerConfig.MySQL.DBName,
	)

	var logLevel logger.LogLevel
	switch global.ServerConfig.MySQL.LogMode {
	case "info":
		logLevel = logger.Info
	case "warn":
		logLevel = logger.Warn
	case "error":
		logLevel = logger.Error
	case "silent":
		logLevel = logger.Silent
	default:
		logLevel = logger.Error
	}
	zap.S().Infof("数据库日志级别设置为: %s", global.ServerConfig.MySQL.LogMode)

	// 创建自定义日志写入器
	logWriter := log.New(os.Stdout, "\r\n", log.LstdFlags)

	newLogger := logger.New(
		logWriter,
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logLevel,    // 日志级别
			IgnoreRecordNotFoundError: true,        // 忽略记录未找到错误
			Colorful:                  true,        // 彩色输出
		},
	)

	config := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: newLogger,
	}

	zap.S().Info("正在连接数据库...")
	db, err := gorm.Open(mysql.Open(dsn), config)
	if err != nil {
		zap.S().Errorf("连接数据库失败: %v", err)
		panic(fmt.Sprintf("连接数据库失败: %v", err))
	}
	zap.S().Info("数据库连接成功")

	sqlDB, err := db.DB()
	if err != nil {
		zap.S().Errorf("获取数据库实例失败: %v", err)
		panic(fmt.Sprintf("获取数据库实例失败: %v", err))
	}

	zap.S().Infof("设置数据库连接池参数: MaxIdleConns=%d, MaxOpenConns=%d",
		global.ServerConfig.MySQL.MaxIdleConns,
		global.ServerConfig.MySQL.MaxOpenConns,
	)
	sqlDB.SetMaxIdleConns(global.ServerConfig.MySQL.MaxIdleConns)
	sqlDB.SetMaxOpenConns(global.ServerConfig.MySQL.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	global.DB = db
	zap.S().Info("数据库初始化完成")
}
