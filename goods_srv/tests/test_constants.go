/*
 * @Description: 商品服务测试数据常量
 * 基于SQL脚本中的真实测试数据定义
 */
package tests

// TestGoodsConstants 测试商品数据常量
type TestGoodsConstants struct {
	ID          int32
	Name        string
	GoodsSN     string
	MarketPrice float32
	ShopPrice   float32
	BrandID     int32
	BrandName   string
	CategoryID  int32
	IsHot       bool
	IsNew       bool
	OnSale      bool
}

// TestBrandConstants 测试品牌数据常量
type TestBrandConstants struct {
	ID   int32
	Name string
	Logo string
}

// TestCategoryConstants 测试分类数据常量  
type TestCategoryConstants struct {
	ID       int32
	Name     string
	Level    int32
	ParentID int32
	IsTab    bool
}

// 预定义的测试品牌数据（与SQL脚本保持一致）
var TestBrands = []TestBrandConstants{
	{ID: 1, Name: "Apple", Logo: "https://img.example.com/brand/apple.png"},
	{ID: 2, Name: "HUAWEI", Logo: "https://img.example.com/brand/huawei.png"},
	{ID: 3, Name: "小米", Logo: "https://img.example.com/brand/xiaomi.png"},
	{ID: 4, Name: "OPPO", Logo: "https://img.example.com/brand/oppo.png"},
	{ID: 5, Name: "vivo", Logo: "https://img.example.com/brand/vivo.png"},
	{ID: 6, Name: "Nike", Logo: "https://img.example.com/brand/nike.png"},
	{ID: 7, Name: "Adidas", Logo: "https://img.example.com/brand/adidas.png"},
	{ID: 8, Name: "UNIQLO", Logo: "https://img.example.com/brand/uniqlo.png"},
	{ID: 9, Name: "IKEA", Logo: "https://img.example.com/brand/ikea.png"},
	{ID: 10, Name: "美的", Logo: "https://img.example.com/brand/midea.png"},
}

// 预定义的测试分类数据（3级分类结构）
var TestCategories = []TestCategoryConstants{
	// 一级分类
	{ID: 1, Name: "电子数码", Level: 1, ParentID: 0, IsTab: true},
	{ID: 2, Name: "服装鞋包", Level: 1, ParentID: 0, IsTab: true},
	{ID: 3, Name: "家居生活", Level: 1, ParentID: 0, IsTab: true},
	{ID: 4, Name: "图书文教", Level: 1, ParentID: 0, IsTab: false},
	{ID: 5, Name: "美妆个护", Level: 1, ParentID: 0, IsTab: false},
	
	// 二级分类（电子数码）
	{ID: 6, Name: "手机通讯", Level: 2, ParentID: 1, IsTab: true},
	{ID: 7, Name: "电脑办公", Level: 2, ParentID: 1, IsTab: true},
	{ID: 8, Name: "数码配件", Level: 2, ParentID: 1, IsTab: false},
	
	// 三级分类（手机通讯）
	{ID: 9, Name: "智能手机", Level: 3, ParentID: 6, IsTab: false},
	{ID: 10, Name: "老人机", Level: 3, ParentID: 6, IsTab: false},
}

// 预定义的测试商品数据（与SQL脚本保持一致）
var TestGoods = []TestGoodsConstants{
	// 高价值商品
	{ID: 1, Name: "iPhone 15 Pro Max 256GB 钛原色", GoodsSN: "IP15PM256TI", MarketPrice: 10999.0, ShopPrice: 9999.0, BrandID: 1, BrandName: "Apple", CategoryID: 9, IsHot: true, IsNew: true, OnSale: true},
	{ID: 2, Name: "HUAWEI Mate 60 Pro 12+512GB 雅川青", GoodsSN: "HWM60P512YQ", MarketPrice: 7999.0, ShopPrice: 7499.0, BrandID: 2, BrandName: "HUAWEI", CategoryID: 9, IsHot: true, IsNew: false, OnSale: true},
	{ID: 3, Name: "小米14 Ultra 16+512GB 钛金属黑", GoodsSN: "MI14U512TJ", MarketPrice: 6999.0, ShopPrice: 6499.0, BrandID: 3, BrandName: "小米", CategoryID: 9, IsHot: true, IsNew: true, OnSale: true},
	
	// 中等价位商品
	{ID: 4, Name: "OPPO Find X7 12+256GB 烟云紫", GoodsSN: "OPFX7256YZ", MarketPrice: 4999.0, ShopPrice: 4599.0, BrandID: 4, BrandName: "OPPO", CategoryID: 9, IsHot: false, IsNew: true, OnSale: true},
	{ID: 5, Name: "vivo X100 Pro 16+512GB 星迹蓝", GoodsSN: "VVX100P512XJ", MarketPrice: 5499.0, ShopPrice: 4999.0, BrandID: 5, BrandName: "vivo", CategoryID: 9, IsHot: false, IsNew: false, OnSale: true},
	
	// 电脑类商品
	{ID: 6, Name: "MacBook Pro 14英寸 M3 Pro芯片", GoodsSN: "MBP14M3PRO", MarketPrice: 31067.0, ShopPrice: 28999.0, BrandID: 1, BrandName: "Apple", CategoryID: 7, IsHot: true, IsNew: true, OnSale: true},
	
	// 服装类商品
	{ID: 7, Name: "Nike Air Max 270 男款跑鞋", GoodsSN: "NKAM270M", MarketPrice: 1299.0, ShopPrice: 999.0, BrandID: 6, BrandName: "Nike", CategoryID: 2, IsHot: true, IsNew: false, OnSale: true},
	{ID: 8, Name: "Adidas Ultraboost 22 跑鞋", GoodsSN: "ADUB22", MarketPrice: 1599.0, ShopPrice: 1199.0, BrandID: 7, BrandName: "Adidas", CategoryID: 2, IsHot: false, IsNew: false, OnSale: true},
	
	// 低价商品
	{ID: 9, Name: "小米手环8 标准版", GoodsSN: "MIHB8STD", MarketPrice: 299.0, ShopPrice: 249.0, BrandID: 3, BrandName: "小米", CategoryID: 8, IsHot: true, IsNew: true, OnSale: true},
	{ID: 10, Name: "UNIQLO 精梳棉圆领T恤", GoodsSN: "UQ_COTTON_T", MarketPrice: 79.0, ShopPrice: 69.0, BrandID: 8, BrandName: "UNIQLO", CategoryID: 2, IsHot: false, IsNew: false, OnSale: true},
}

// 价格区间常量
var PriceRanges = struct {
	Low    struct{ Min, Max float32 }
	Medium struct{ Min, Max float32 }
	High   struct{ Min, Max float32 }
}{
	Low:    struct{ Min, Max float32 }{Min: 0, Max: 500},
	Medium: struct{ Min, Max float32 }{Min: 500, Max: 5000},
	High:   struct{ Min, Max float32 }{Min: 5000, Max: 50000},
}

// GetTestGoodsByID 根据ID获取测试商品
func GetTestGoodsByID(id int32) *TestGoodsConstants {
	for i := range TestGoods {
		if TestGoods[i].ID == id {
			return &TestGoods[i]
		}
	}
	return nil
}

// GetTestGoodsByBrand 根据品牌ID获取测试商品
func GetTestGoodsByBrand(brandID int32) []TestGoodsConstants {
	var goods []TestGoodsConstants
	for _, g := range TestGoods {
		if g.BrandID == brandID {
			goods = append(goods, g)
		}
	}
	return goods
}

// GetHotTestGoods 获取热销测试商品
func GetHotTestGoods() []TestGoodsConstants {
	var goods []TestGoodsConstants
	for _, g := range TestGoods {
		if g.IsHot {
			goods = append(goods, g)
		}
	}
	return goods
}

// GetNewTestGoods 获取新品测试商品
func GetNewTestGoods() []TestGoodsConstants {
	var goods []TestGoodsConstants
	for _, g := range TestGoods {
		if g.IsNew {
			goods = append(goods, g)
		}
	}
	return goods
}

// GetTestBrandByID 根据ID获取测试品牌
func GetTestBrandByID(id int32) *TestBrandConstants {
	for i := range TestBrands {
		if TestBrands[i].ID == id {
			return &TestBrands[i]
		}
	}
	return nil
}

// GetTestCategoryByID 根据ID获取测试分类
func GetTestCategoryByID(id int32) *TestCategoryConstants {
	for i := range TestCategories {
		if TestCategories[i].ID == id {
			return &TestCategories[i]
		}
	}
	return nil
}

// GetTestCategoriesByLevel 根据级别获取测试分类
func GetTestCategoriesByLevel(level int32) []TestCategoryConstants {
	var categories []TestCategoryConstants
	for _, c := range TestCategories {
		if c.Level == level {
			categories = append(categories, c)
		}
	}
	return categories
}

// TestSearchKeywords 测试搜索关键词
var TestSearchKeywords = []string{
	"iPhone",
	"手机",
	"Apple",
	"Nike",
	"跑鞋",
	"小米",
}

// TestScenarios 测试场景常量
var TestScenarios = struct {
	GoodsList     string
	GoodsDetail   string
	GoodsSearch   string
	BrandList     string
	CategoryList  string
	HotGoods      string
	NewGoods      string
}{
	GoodsList:    "商品列表查询测试",
	GoodsDetail:  "商品详情查询测试",
	GoodsSearch:  "商品搜索测试",
	BrandList:    "品牌列表查询测试", 
	CategoryList: "分类列表查询测试",
	HotGoods:     "热销商品查询测试",
	NewGoods:     "新品商品查询测试",
}