/*
 * @Description:
 */
package model

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"order_srv/config"
	"order_srv/global"

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

	// 自动迁移订单相关表结构
	if err := db.AutoMigrate(&OrderInfo{}, &OrderGoods{}, &ShoppingCart{}); err != nil {
		t.Fatalf("自动迁移表结构失败: %v", err)
	}
}

// TestOrderInfoCRUD 测试订单信息的增删改查
func TestOrderInfoCRUD(t *testing.T) {
	setupTestDB(t)

	order := &OrderInfo{
		User:         1,
		OrderSn:      "ORDER123456",
		PayType:      "alipay",
		Status:       "PAYING",
		TradeNo:      "TRADE123456",
		OrderMount:   99.99,
		PayTime:      nil, // 使用nil，让数据库设置为NULL
		Address:      "测试地址",
		SignerName:   "张三",
		SingerMobile: "13800000000",
		Post:         "请尽快发货",
	}

	// 创建
	t.Run("Create", func(t *testing.T) {
		result := global.DB.Create(order)
		if result.Error != nil {
			t.Fatalf("创建订单信息失败: %v", result.Error)
		}
		if order.ID == 0 {
			t.Fatal("创建订单信息后ID为0")
		}
	})

	// 查询
	t.Run("Read", func(t *testing.T) {
		var found OrderInfo
		result := global.DB.First(&found, order.ID)
		if result.Error != nil {
			t.Fatalf("查询订单信息失败: %v", result.Error)
		}
		if found.OrderSn != order.OrderSn {
			t.Errorf("订单号不匹配，期望: %s, 实际: %s", order.OrderSn, found.OrderSn)
		}
	})

	// 更新
	t.Run("Update", func(t *testing.T) {
		newStatus := "TRADE_SUCCESS"
		result := global.DB.Model(order).Update("status", newStatus)
		if result.Error != nil {
			t.Fatalf("更新订单状态失败: %v", result.Error)
		}
		var updated OrderInfo
		global.DB.First(&updated, order.ID)
		if updated.Status != newStatus {
			t.Errorf("更新后的订单状态不匹配，期望: %s, 实际: %s", newStatus, updated.Status)
		}
	})

	// 删除
	t.Run("Delete", func(t *testing.T) {
		result := global.DB.Delete(order)
		if result.Error != nil {
			t.Fatalf("删除订单信息失败: %v", result.Error)
		}
		var deleted OrderInfo
		result = global.DB.First(&deleted, order.ID)
		if result.Error == nil {
			t.Error("删除后仍能查询到订单信息")
		}
	})
}

// TestOrderGoodsCRUD 测试订单商品的增删改查
func TestOrderGoodsCRUD(t *testing.T) {
	setupTestDB(t)

	orderGoods := &OrderGoods{
		Order:      1,
		Goods:      2,
		GoodsName:  "测试商品",
		GoodsImage: "http://example.com/image.jpg",
		GoodsPrice: 88.88,
		Nums:       3,
	}

	// 创建
	t.Run("Create", func(t *testing.T) {
		result := global.DB.Create(orderGoods)
		if result.Error != nil {
			t.Fatalf("创建订单商品失败: %v", result.Error)
		}
		if orderGoods.ID == 0 {
			t.Fatal("创建订单商品后ID为0")
		}
	})

	// 查询
	t.Run("Read", func(t *testing.T) {
		var found OrderGoods
		result := global.DB.First(&found, orderGoods.ID)
		if result.Error != nil {
			t.Fatalf("查询订单商品失败: %v", result.Error)
		}
		if found.GoodsName != orderGoods.GoodsName {
			t.Errorf("商品名称不匹配，期望: %s, 实际: %s", orderGoods.GoodsName, found.GoodsName)
		}
	})

	// 更新
	t.Run("Update", func(t *testing.T) {
		newNums := int32(5)
		result := global.DB.Model(orderGoods).Update("nums", newNums)
		if result.Error != nil {
			t.Fatalf("更新商品数量失败: %v", result.Error)
		}
		var updated OrderGoods
		global.DB.First(&updated, orderGoods.ID)
		if updated.Nums != newNums {
			t.Errorf("更新后的商品数量不匹配，期望: %d, 实际: %d", newNums, updated.Nums)
		}
	})

	// 删除
	t.Run("Delete", func(t *testing.T) {
		result := global.DB.Delete(orderGoods)
		if result.Error != nil {
			t.Fatalf("删除订单商品失败: %v", result.Error)
		}
		var deleted OrderGoods
		result = global.DB.First(&deleted, orderGoods.ID)
		if result.Error == nil {
			t.Error("删除后仍能查询到订单商品")
		}
	})
}

// TestShoppingCartCRUD 测试购物车的增删改查
func TestShoppingCartCRUD(t *testing.T) {
	setupTestDB(t)

	cart := &ShoppingCart{
		User:    1,
		Goods:   2,
		Nums:    3,
		Checked: true,
	}

	// 创建
	t.Run("Create", func(t *testing.T) {
		result := global.DB.Create(cart)
		if result.Error != nil {
			t.Fatalf("创建购物车失败: %v", result.Error)
		}
		if cart.ID == 0 {
			t.Fatal("创建购物车后ID为0")
		}
	})

	// 查询
	t.Run("Read", func(t *testing.T) {
		var found ShoppingCart
		result := global.DB.First(&found, cart.ID)
		if result.Error != nil {
			t.Fatalf("查询购物车失败: %v", result.Error)
		}
		if found.Goods != cart.Goods {
			t.Errorf("商品ID不匹配，期望: %d, 实际: %d", cart.Goods, found.Goods)
		}
	})

	// 更新
	t.Run("Update", func(t *testing.T) {
		newNums := int32(5)
		result := global.DB.Model(cart).Update("nums", newNums)
		if result.Error != nil {
			t.Fatalf("更新购物车数量失败: %v", result.Error)
		}
		var updated ShoppingCart
		global.DB.First(&updated, cart.ID)
		if updated.Nums != newNums {
			t.Errorf("更新后的购物车数量不匹配，期望: %d, 实际: %d", newNums, updated.Nums)
		}
	})

	// 删除
	t.Run("Delete", func(t *testing.T) {
		result := global.DB.Delete(cart)
		if result.Error != nil {
			t.Fatalf("删除购物车失败: %v", result.Error)
		}
		var deleted ShoppingCart
		result = global.DB.First(&deleted, cart.ID)
		if result.Error == nil {
			t.Error("删除后仍能查询到购物车")
		}
	})
}
