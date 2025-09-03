/*
 * @Description: 库存服务测试数据常量
 * 基于SQL脚本中的真实测试数据定义
 */
package tests

// TestInventoryConstants 测试库存数据常量
type TestInventoryConstants struct {
	GoodsID   int32
	GoodsName string
	Stock     int32
	Version   int32
	Status    string // 库存状态：充足、正常、低库存、零库存
}

// 预定义的测试库存数据（与SQL脚本保持一致）
var TestInventories = []TestInventoryConstants{
	// 充足库存商品
	{GoodsID: 1, GoodsName: "iPhone 15 Pro Max 256GB", Stock: 156, Version: 0, Status: "充足"},
	{GoodsID: 6, GoodsName: "MacBook Pro 14英寸", Stock: 23, Version: 0, Status: "正常"},
	{GoodsID: 9, GoodsName: "小米手环8 标准版", Stock: 500, Version: 0, Status: "充足"},
	{GoodsID: 10, GoodsName: "UNIQLO 圆领T恤", Stock: 1000, Version: 0, Status: "充足"},
	
	// 正常库存商品
	{GoodsID: 2, GoodsName: "HUAWEI Mate 60 Pro", Stock: 42, Version: 0, Status: "正常"},
	{GoodsID: 3, GoodsName: "小米14 Ultra", Stock: 33, Version: 0, Status: "正常"},
	{GoodsID: 7, GoodsName: "Nike Air Max 270", Stock: 88, Version: 0, Status: "正常"},
	{GoodsID: 8, GoodsName: "Adidas Ultraboost 22", Stock: 67, Version: 0, Status: "正常"},
	
	// 低库存商品（库存 ≤ 10）
	{GoodsID: 4, GoodsName: "OPPO Find X7", Stock: 8, Version: 0, Status: "低库存"},
	{GoodsID: 5, GoodsName: "vivo X100 Pro", Stock: 5, Version: 0, Status: "低库存"},
	
	// 零库存商品（用于测试库存不足场景）
	{GoodsID: 11, GoodsName: "测试商品（零库存）", Stock: 0, Version: 0, Status: "零库存"},
}

// 库存级别常量
var StockLevels = struct {
	Zero      string
	Low       string
	Normal    string
	Abundant  string
	Threshold struct {
		Low    int32
		Normal int32
		High   int32
	}
}{
	Zero:     "零库存",
	Low:      "低库存",
	Normal:   "正常库存",  
	Abundant: "充足库存",
	Threshold: struct {
		Low    int32
		Normal int32
		High   int32
	}{
		Low:    10,  // ≤10 为低库存
		Normal: 50,  // 11-50 为正常库存
		High:   200, // >50 为充足库存
	},
}

// GetTestInventoryByGoodsID 根据商品ID获取测试库存
func GetTestInventoryByGoodsID(goodsID int32) *TestInventoryConstants {
	for i := range TestInventories {
		if TestInventories[i].GoodsID == goodsID {
			return &TestInventories[i]
		}
	}
	return nil
}

// GetLowStockInventories 获取低库存商品
func GetLowStockInventories() []TestInventoryConstants {
	var inventories []TestInventoryConstants
	for _, inv := range TestInventories {
		if inv.Stock <= StockLevels.Threshold.Low {
			inventories = append(inventories, inv)
		}
	}
	return inventories
}

// GetZeroStockInventories 获取零库存商品
func GetZeroStockInventories() []TestInventoryConstants {
	var inventories []TestInventoryConstants
	for _, inv := range TestInventories {
		if inv.Stock == 0 {
			inventories = append(inventories, inv)
		}
	}
	return inventories
}

// GetAbundantStockInventories 获取充足库存商品
func GetAbundantStockInventories() []TestInventoryConstants {
	var inventories []TestInventoryConstants
	for _, inv := range TestInventories {
		if inv.Stock > StockLevels.Threshold.High {
			inventories = append(inventories, inv)
		}
	}
	return inventories
}

// GetTestInventoriesForConcurrency 获取用于并发测试的库存商品
func GetTestInventoriesForConcurrency() []TestInventoryConstants {
	// 返回库存充足的商品，适合并发扣减测试
	var inventories []TestInventoryConstants
	for _, inv := range TestInventories {
		if inv.Stock >= 50 { // 库存足够支持并发测试
			inventories = append(inventories, inv)
		}
	}
	return inventories
}

// TestSellOperations 测试扣减操作数据
type TestSellOperations struct {
	GoodsID     int32
	SellNum     int32
	ExpectedResult string // "success" or "insufficient"
}

// 预定义的扣减测试场景
var TestSellScenarios = []TestSellOperations{
	// 正常扣减场景
	{GoodsID: 1, SellNum: 10, ExpectedResult: "success"},   // iPhone库存156，扣减10应该成功
	{GoodsID: 2, SellNum: 5, ExpectedResult: "success"},    // HUAWEI库存42，扣减5应该成功
	{GoodsID: 9, SellNum: 50, ExpectedResult: "success"},   // 小米手环库存500，扣减50应该成功
	
	// 库存不足场景
	{GoodsID: 4, SellNum: 20, ExpectedResult: "insufficient"}, // OPPO库存8，扣减20应该失败
	{GoodsID: 5, SellNum: 10, ExpectedResult: "insufficient"}, // vivo库存5，扣减10应该失败
	{GoodsID: 11, SellNum: 1, ExpectedResult: "insufficient"}, // 零库存商品，扣减1应该失败
	
	// 边界测试场景
	{GoodsID: 4, SellNum: 8, ExpectedResult: "success"},     // OPPO库存8，扣减8应该成功（扣完）
	{GoodsID: 5, SellNum: 5, ExpectedResult: "success"},     // vivo库存5，扣减5应该成功（扣完）
}

// TestRebackOperations 测试回滚操作数据
type TestRebackOperations struct {
	GoodsID    int32
	RebackNum  int32
	OriginalStock int32
}

// 预定义的回滚测试场景
var TestRebackScenarios = []TestRebackOperations{
	{GoodsID: 1, RebackNum: 10, OriginalStock: 156}, // iPhone回滚10件
	{GoodsID: 2, RebackNum: 5, OriginalStock: 42},   // HUAWEI回滚5件
	{GoodsID: 4, RebackNum: 3, OriginalStock: 8},    // OPPO回滚3件
}

// TestBatchOperations 批量操作测试数据
type TestBatchOperations struct {
	Operations []TestSellOperations
	ExpectedResult string
}

// 预定义的批量操作测试场景
var TestBatchScenarios = []TestBatchOperations{
	// 全部成功的批量操作
	{
		Operations: []TestSellOperations{
			{GoodsID: 1, SellNum: 5, ExpectedResult: "success"},
			{GoodsID: 2, SellNum: 3, ExpectedResult: "success"},
			{GoodsID: 9, SellNum: 10, ExpectedResult: "success"},
		},
		ExpectedResult: "all_success",
	},
	// 包含失败的批量操作（应该全部回滚）
	{
		Operations: []TestSellOperations{
			{GoodsID: 1, SellNum: 5, ExpectedResult: "success"},
			{GoodsID: 4, SellNum: 20, ExpectedResult: "insufficient"}, // 这个会失败
			{GoodsID: 9, SellNum: 10, ExpectedResult: "success"},
		},
		ExpectedResult: "all_failed", // 因为一个失败，全部应该回滚
	},
}

// TestScenarios 测试场景常量
var TestScenarios = struct {
	SetInventory     string
	GetInventory     string
	SellInventory    string
	RebackInventory  string
	ConcurrentSell   string
	BatchOperation   string
	LowStockAlert    string
	ZeroStockCheck   string
}{
	SetInventory:    "设置库存测试",
	GetInventory:    "查询库存测试",
	SellInventory:   "扣减库存测试",
	RebackInventory: "回滚库存测试",
	ConcurrentSell:  "并发扣减测试",
	BatchOperation:  "批量操作测试",
	LowStockAlert:   "低库存预警测试",
	ZeroStockCheck:  "零库存检查测试",
}

// GetStockStatus 根据库存数量获取库存状态
func GetStockStatus(stock int32) string {
	if stock == 0 {
		return StockLevels.Zero
	} else if stock <= StockLevels.Threshold.Low {
		return StockLevels.Low
	} else if stock <= StockLevels.Threshold.Normal {
		return StockLevels.Normal
	} else {
		return StockLevels.Abundant
	}
}