package initialize

import (
	"fmt"
	"goods_srv/global"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func InitDB() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		global.ServerConfig.MySQL.User,
		global.ServerConfig.MySQL.Password,
		global.ServerConfig.MySQL.Host,
		global.ServerConfig.MySQL.Port,
		global.ServerConfig.MySQL.DBName,
	)

	var logLevel logger.LogLevel
	if global.ServerConfig.MySQL.LogMode == "info" {
		logLevel = logger.Info
	} else {
		logLevel = logger.Error
	}

	config := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: logger.Default.LogMode(logLevel),
	}

	db, err := gorm.Open(mysql.Open(dsn), config)
	if err != nil {
		panic(fmt.Sprintf("连接数据库失败: %v", err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Sprintf("获取数据库实例失败: %v", err))
	}

	sqlDB.SetMaxIdleConns(global.ServerConfig.MySQL.MaxIdleConns)
	sqlDB.SetMaxOpenConns(global.ServerConfig.MySQL.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	global.DB = db
}
