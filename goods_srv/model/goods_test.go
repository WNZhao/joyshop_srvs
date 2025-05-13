package model

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"goods_srv/config"
	"goods_srv/global"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// 测试数据
var (
	testCategories = []struct {
		Name     string
		ParentId uint
		Level    int
		Sort     int
		IsTab    bool
	}{
		{"手机数码", 0, 1, 1, true},
		{"电脑办公", 0, 1, 2, true},
		{"家用电器", 0, 1, 3, true},
		{"智能手机", 1, 2, 1, false},
		{"平板电脑", 1, 2, 2, false},
		{"笔记本", 2, 2, 1, false},
		{"台式机", 2, 2, 2, false},
	}

	testBrands = []struct {
		Name string
		Logo string
		Desc string
	}{
		{
			"Apple",
			"https://example.com/brands/apple.png",
			"Apple Inc. 是一家美国跨国科技公司",
		},
		{
			"Samsung",
			"https://example.com/brands/samsung.png",
			"三星集团是韩国最大的跨国企业集团",
		},
		{
			"Huawei",
			"https://example.com/brands/huawei.png",
			"华为技术有限公司是一家中国科技公司",
		},
	}

	testGoods = []struct {
		Name            string
		GoodsSn         string
		MarketPrice     float64
		ShopPrice       float64
		GoodsBrief      string
		Images          []string
		DescImages      []string
		GoodsFrontImage string
	}{
		{
			"iPhone 15 Pro Max",
			"IP15PM256",
			9999.00,
			8999.00,
			"Apple iPhone 15 Pro Max 256GB 钛金属",
			[]string{
				"https://example.com/goods/iphone15/1.jpg",
				"https://example.com/goods/iphone15/2.jpg",
				"https://example.com/goods/iphone15/3.jpg",
			},
			[]string{
				"https://example.com/goods/iphone15/desc1.jpg",
				"https://example.com/goods/iphone15/desc2.jpg",
				"https://example.com/goods/iphone15/desc3.jpg",
			},
			"https://example.com/goods/iphone15/front.jpg",
		},
		{
			"Samsung Galaxy S24 Ultra",
			"SG24U512",
			8999.00,
			7999.00,
			"Samsung Galaxy S24 Ultra 512GB 钛金属",
			[]string{
				"https://example.com/goods/s24/1.jpg",
				"https://example.com/goods/s24/2.jpg",
				"https://example.com/goods/s24/3.jpg",
			},
			[]string{
				"https://example.com/goods/s24/desc1.jpg",
				"https://example.com/goods/s24/desc2.jpg",
				"https://example.com/goods/s24/desc3.jpg",
			},
			"https://example.com/goods/s24/front.jpg",
		},
		{
			"Huawei Mate 60 Pro",
			"HM60P512",
			7999.00,
			6999.00,
			"华为 Mate 60 Pro 512GB 曜金黑",
			[]string{
				"https://example.com/goods/mate60/1.jpg",
				"https://example.com/goods/mate60/2.jpg",
				"https://example.com/goods/mate60/3.jpg",
			},
			[]string{
				"https://example.com/goods/mate60/desc1.jpg",
				"https://example.com/goods/mate60/desc2.jpg",
				"https://example.com/goods/mate60/desc3.jpg",
			},
			"https://example.com/goods/mate60/front.jpg",
		},
	}
)

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
		cfg.DBConfig.User,
		cfg.DBConfig.Password,
		cfg.DBConfig.Host,
		cfg.DBConfig.Port,
		cfg.DBConfig.DBName,
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
	if !db.Migrator().HasTable(&Category{}) {
		if err := db.AutoMigrate(
			&Category{},
			&Brand{},
			&Goods{},
			&CategoryBrand{},
		); err != nil {
			t.Fatalf("自动迁移表结构失败: %v", err)
		}
	}
}

func TestGoodsCRUD(t *testing.T) {
	// 设置测试数据库
	setupTestDB(t)

	// 创建测试分类
	var categoryId uint
	for _, cat := range testCategories {
		category := &Category{
			Name:     cat.Name,
			ParentId: cat.ParentId,
			Level:    cat.Level,
			Sort:     cat.Sort,
			IsTab:    cat.IsTab,
		}
		if err := global.DB.Create(category).Error; err != nil {
			t.Fatalf("创建分类失败: %v", err)
		}
		if cat.Name == "智能手机" {
			categoryId = category.ID
		}
	}

	// 创建测试品牌
	var brandId uint
	for _, b := range testBrands {
		brand := &Brand{
			Name: b.Name,
			Logo: b.Logo,
			Desc: b.Desc,
		}
		if err := global.DB.Create(brand).Error; err != nil {
			t.Fatalf("创建品牌失败: %v", err)
		}
		if b.Name == "Apple" {
			brandId = brand.ID
		}
	}

	// 创建测试商品
	goods := &Goods{
		BrandId:         brandId,
		OnSale:          true,
		ShipFree:        true,
		IsNew:           true,
		IsHot:           true,
		Name:            testGoods[0].Name,
		GoodsSn:         testGoods[0].GoodsSn,
		MarketPrice:     testGoods[0].MarketPrice,
		ShopPrice:       testGoods[0].ShopPrice,
		GoodsBrief:      testGoods[0].GoodsBrief,
		Images:          testGoods[0].Images,
		DescImages:      testGoods[0].DescImages,
		GoodsFrontImage: testGoods[0].GoodsFrontImage,
		Status:          1,
	}

	// 关联分类
	var category Category
	if err := global.DB.First(&category, categoryId).Error; err != nil {
		t.Fatalf("获取分类失败: %v", err)
	}
	goods.Categories = []Category{category}

	// 测试创建商品
	if err := CreateGoods(goods); err != nil {
		t.Fatalf("创建商品失败: %v", err)
	}

	// 测试获取商品
	retrievedGoods, err := GetGoodsById(goods.ID)
	if err != nil {
		t.Fatalf("获取商品失败: %v", err)
	}
	if retrievedGoods.Name != goods.Name {
		t.Errorf("获取的商品名称不匹配，期望: %s, 实际: %s", goods.Name, retrievedGoods.Name)
	}

	// 测试更新商品
	retrievedGoods.Name = "iPhone 15 Pro Max 1TB"
	if err := UpdateGoods(retrievedGoods); err != nil {
		t.Fatalf("更新商品失败: %v", err)
	}

	// 验证更新
	updatedGoods, err := GetGoodsById(goods.ID)
	if err != nil {
		t.Fatalf("获取更新后的商品失败: %v", err)
	}
	if updatedGoods.Name != "iPhone 15 Pro Max 1TB" {
		t.Errorf("商品名称更新失败，期望: %s, 实际: %s", "iPhone 15 Pro Max 1TB", updatedGoods.Name)
	}

	// 测试删除商品
	if err := DeleteGoods(goods.ID); err != nil {
		t.Fatalf("删除商品失败: %v", err)
	}

	// 验证删除
	_, err = GetGoodsById(goods.ID)
	if err == nil {
		t.Error("商品未被成功删除")
	}
}

func TestGoodsList(t *testing.T) {
	// 设置测试数据库
	setupTestDB(t)

	// 先创建品牌
	var brandId uint
	for _, b := range testBrands {
		brand := &Brand{
			Name: b.Name,
			Logo: b.Logo,
			Desc: b.Desc,
		}
		if err := global.DB.Create(brand).Error; err != nil {
			t.Fatalf("创建品牌失败: %v", err)
		}
		if b.Name == "Apple" {
			brandId = brand.ID
		}
	}

	// 创建测试商品
	for _, g := range testGoods {
		goods := &Goods{
			BrandId:         brandId, // 设置品牌ID
			Name:            g.Name,
			GoodsSn:         g.GoodsSn,
			MarketPrice:     g.MarketPrice,
			ShopPrice:       g.ShopPrice,
			GoodsBrief:      g.GoodsBrief,
			Images:          g.Images,
			DescImages:      g.DescImages,
			GoodsFrontImage: g.GoodsFrontImage,
			Status:          1,
		}
		if err := CreateGoods(goods); err != nil {
			t.Fatalf("创建测试商品失败: %v", err)
		}
	}

	// 测试分页获取商品列表
	goods, total, err := GetGoodsList(1, 2)
	if err != nil {
		t.Fatalf("获取商品列表失败: %v", err)
	}
	if len(goods) != 2 {
		t.Errorf("获取的商品数量不正确，期望: 2, 实际: %d", len(goods))
	}
	if total < 3 {
		t.Errorf("商品总数不正确，期望: >=3, 实际: %d", total)
	}

	// 测试第二页
	goods, total, err = GetGoodsList(2, 2)
	if err != nil {
		t.Fatalf("获取第二页商品列表失败: %v", err)
	}
	if len(goods) != 1 {
		t.Errorf("第二页商品数量不正确，期望: 1, 实际: %d", len(goods))
	}
}
