/*
 * @Author: Will nanan_zhao@163.com
 * @Date: 2025-05-08 17:09:45
 * @LastEditors: Will nanan_zhao@163.com
 * @LastEditTime: 2025-05-10 14:44:29
 * @FilePath: /joyshop_srvs/user_srv/initialize/db.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package initialize

import (
	"fmt"
	"log"
	"os"
	"time"

	"joyshop_srvs/user_srv/global"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// InitDB 初始化数据库
func InitDB() error {
	// 构建DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		global.GlobalConfig.DBConfig.User,
		global.GlobalConfig.DBConfig.Password,
		global.GlobalConfig.DBConfig.Host,
		global.GlobalConfig.DBConfig.Port,
		global.GlobalConfig.DBConfig.DBName,
	)

	zap.S().Debugf("数据库连接信息: %s", dsn)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // 使用标准输出
		logger.Config{
			SlowThreshold: time.Second, // 慢查询阈值
			LogLevel:      logger.Info, // 日志级别
			Colorful:      true,        // 彩色
		},
	)

	// 配置GORM
	gormConfig := &gorm.Config{
		// 使用单数表名
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		// 配置日志
		Logger: newLogger, //logger.Default.LogMode(logger.Info),
	}

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		return fmt.Errorf("连接数据库失败: %v", err)
	}

	// 设置全局变量
	global.DB = db

	zap.S().Info("数据库连接成功")
	return nil
}
