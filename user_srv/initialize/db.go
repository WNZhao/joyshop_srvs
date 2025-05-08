package initialize

import (
	"fmt"
	"time"

	"joyshop_srvs/user_srv/config/db"
	"joyshop_srvs/user_srv/global"
	"joyshop_srvs/user_srv/util"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// InitDB 初始化数据库连接
func InitDB() {
	// 加载配置
	config := &db.Config{}
	configPath := util.GetConfigPath()
	zap.S().Debugf("正在加载配置文件: %s", configPath)

	err := util.LoadConfig(configPath, config)
	if err != nil {
		zap.S().Errorf("加载配置文件失败: %v", err)
		zap.S().Fatalf("配置文件路径: %s", configPath)
	}

	// 拼接DSN
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Database.User,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		config.Database.Name,
	)
	zap.S().Debugf("数据库连接信息: %s:%d/%s", config.Database.Host, config.Database.Port, config.Database.Name)

	// 设置打印日志
	newLogger := logger.New(
		zap.NewStdLog(zap.L()), // 使用zap作为日志输出
		logger.Config{
			SlowThreshold: time.Second * 10,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)
	// 连接数据库
	var errDB error

	global.DB, errDB = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
		Logger: newLogger,
	})
	if errDB != nil {
		zap.S().Errorf("连接数据库失败: %v", errDB)
		zap.S().Fatalf("数据库连接信息: %s", dsn)
	}

	zap.S().Info("数据库连接成功")
	// 定义表结构，将表结构直接生成对应的表 自动建表
	//db.AutoMigrate(&model.User{})
}
