package model

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"inventory_srv/config"
	"inventory_srv/global"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// 数据库模型测试 商品信息
func setupTestDB(t *testing.T) {
	// 初始化日志
	zap.S().Info("开始初始化测试数据库...")

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

	// 在测试开始前检查表是否存在
	if !db.Migrator().HasTable(&Inventory{}) {
		if err := db.AutoMigrate(
			&Inventory{},
		); err != nil {
			t.Fatalf("自动迁移表结构失败: %v", err)
		}
	}
}

// TestInventoryCRUD 测试库存的增删改查
func TestInventoryCRUD(t *testing.T) {
	// 设置测试数据库
	setupTestDB(t)

	// 创建测试数据
	inventory := &Inventory{
		GoodsID: 1,
		Stock:   100,
		Version: 0,
	}

	// 测试创建
	t.Run("Create", func(t *testing.T) {
		result := global.DB.Create(inventory)
		if result.Error != nil {
			t.Fatalf("创建库存记录失败: %v", result.Error)
		}
		if inventory.ID == 0 {
			t.Fatal("创建库存记录后ID为0")
		}
	})

	// 测试查询
	t.Run("Read", func(t *testing.T) {
		var found Inventory
		result := global.DB.First(&found, inventory.ID)
		if result.Error != nil {
			t.Fatalf("查询库存记录失败: %v", result.Error)
		}
		if found.GoodsID != inventory.GoodsID {
			t.Errorf("查询到的商品ID不匹配，期望: %d, 实际: %d", inventory.GoodsID, found.GoodsID)
		}
		if found.Stock != inventory.Stock {
			t.Errorf("查询到的库存数量不匹配，期望: %d, 实际: %d", inventory.Stock, found.Stock)
		}
	})

	// 测试更新
	t.Run("Update", func(t *testing.T) {
		newStock := int32(200)
		result := global.DB.Model(inventory).Update("stock", newStock)
		if result.Error != nil {
			t.Fatalf("更新库存记录失败: %v", result.Error)
		}

		// 验证更新
		var updated Inventory
		global.DB.First(&updated, inventory.ID)
		if updated.Stock != newStock {
			t.Errorf("更新后的库存数量不匹配，期望: %d, 实际: %d", newStock, updated.Stock)
		}
	})

	// 测试删除
	t.Run("Delete", func(t *testing.T) {
		result := global.DB.Delete(inventory)
		if result.Error != nil {
			t.Fatalf("删除库存记录失败: %v", result.Error)
		}

		// 验证删除
		var deleted Inventory
		result = global.DB.First(&deleted, inventory.ID)
		if result.Error == nil {
			t.Error("删除后仍能查询到记录")
		}
	})
}
