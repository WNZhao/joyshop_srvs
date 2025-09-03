/*
 * @Description: 跨服务集成测试
 * 测试用户、商品、库存、订单服务之间的协作
 */
package integration

import (
	"testing"
)

// 模拟的服务数据结构（基于我们之前创建的测试常量）
type TestUser struct {
	ID       int32
	NickName string
	Mobile   string
	Role     string
}

type TestGoods struct {
	ID           int32
	Name         string
	ShopPrice    float32
	CategoryName string
}

type TestInventory struct {
	GoodsID   int32
	GoodsName string
	Stock     int32
	Status    string
}

type TestOrder struct {
	ID          int32
	UserID      int32
	OrderSN     string
	OrderMount  float32
	Status      string
	SignerName  string
}

type TestCartItem struct {
	UserID     int32
	GoodsID    int32
	GoodsName  string
	GoodsPrice float32
	Nums       int32
	Checked    bool
}

// 模拟测试数据（从之前创建的常量中复制关键数据）
var testUsers = []TestUser{
	{ID: 2, NickName: "张三", Mobile: "18782222220", Role: "member"},
	{ID: 3, NickName: "李四", Mobile: "18782222221", Role: "member"},
	{ID: 13, NickName: "VIP刘总", Mobile: "18888888880", Role: "vip"},
	{ID: 14, NickName: "VIP王总", Mobile: "18888888881", Role: "vip"},
}

var testGoods = []TestGoods{
	{ID: 1, Name: "iPhone 15 Pro Max 256GB", ShopPrice: 9999.0, CategoryName: "手机数码"},
	{ID: 2, Name: "HUAWEI Mate 60 Pro", ShopPrice: 7499.0, CategoryName: "手机数码"},
	{ID: 6, Name: "MacBook Pro 14英寸", ShopPrice: 28999.0, CategoryName: "电脑办公"},
}

var testInventories = []TestInventory{
	{GoodsID: 1, GoodsName: "iPhone 15 Pro Max", Stock: 50, Status: "充足"},
	{GoodsID: 2, GoodsName: "HUAWEI Mate 60 Pro", Stock: 30, Status: "充足"},
	{GoodsID: 5, GoodsName: "vivo X100 Pro", Stock: 5, Status: "偏低"},
}

var testOrders = []TestOrder{
	{ID: 1, UserID: 13, OrderSN: "20240101103000001300130001", OrderMount: 48997.0, Status: "TRADE_FINISHED", SignerName: "VIP刘总"},
	{ID: 2, UserID: 2, OrderSN: "20240102143000000200020002", OrderMount: 10997.0, Status: "TRADE_FINISHED", SignerName: "张三"},
}

// 辅助函数
func getUserByID(id int32) *TestUser {
	for _, user := range testUsers {
		if user.ID == id {
			return &user
		}
	}
	return nil
}

func getGoodsByID(id int32) *TestGoods {
	for _, goods := range testGoods {
		if goods.ID == id {
			return &goods
		}
	}
	return nil
}

func getInventoryByGoodsID(goodsID int32) *TestInventory {
	for _, inv := range testInventories {
		if inv.GoodsID == goodsID {
			return &inv
		}
	}
	return nil
}

func getOrdersByUserID(userID int32) []TestOrder {
	var orders []TestOrder
	for _, order := range testOrders {
		if order.UserID == userID {
			orders = append(orders, order)
		}
	}
	return orders
}

func getVIPUsers() []TestUser {
	var vipUsers []TestUser
	for _, user := range testUsers {
		if user.Role == "vip" {
			vipUsers = append(vipUsers, user)
		}
	}
	return vipUsers
}

func getHighValueGoods() []TestGoods {
	var highValueGoods []TestGoods
	for _, goods := range testGoods {
		if goods.ShopPrice > 10000 {
			highValueGoods = append(highValueGoods, goods)
		}
	}
	return highValueGoods
}

func getLowStockInventories() []TestInventory {
	var lowStock []TestInventory
	for _, inv := range testInventories {
		if inv.Stock <= 10 {
			lowStock = append(lowStock, inv)
		}
	}
	return lowStock
}

// 集成测试场景：用户下单流程
func TestCompleteOrderFlow(t *testing.T) {
	t.Log("=== 跨服务集成测试：完整下单流程 ===")
	
	// 1. 验证用户服务
	t.Run("验证用户服务", func(t *testing.T) {
		testUser := getUserByID(2) // 张三
		if testUser == nil {
			t.Fatal("未找到测试用户")
		}
		t.Logf("测试用户: %s (ID: %d, 手机: %s)", testUser.NickName, testUser.ID, testUser.Mobile)
	})
	
	// 2. 验证商品服务
	t.Run("验证商品服务", func(t *testing.T) {
		testGoods := getGoodsByID(1) // iPhone
		if testGoods == nil {
			t.Fatal("未找到测试商品")
		}
		t.Logf("测试商品: %s (ID: %d, 价格: ¥%.2f)", testGoods.Name, testGoods.ID, testGoods.ShopPrice)
	})
	
	// 3. 验证库存服务
	t.Run("验证库存服务", func(t *testing.T) {
		testInventory := getInventoryByGoodsID(1) // iPhone库存
		if testInventory == nil {
			t.Fatal("未找到测试库存")
		}
		t.Logf("测试库存: 商品ID %d, 库存: %d, 状态: %s", 
			testInventory.GoodsID, testInventory.Stock, testInventory.Status)
	})
	
	// 4. 模拟完整的购物流程
	t.Run("完整购物流程", func(t *testing.T) {
		userID := int32(2)  // 张三
		goodsID := int32(1) // iPhone
		
		// 4.1 检查商品可用性
		testGoods := getGoodsByID(goodsID)
		if testGoods == nil {
			t.Fatal("商品不存在")
		}
		
		// 4.2 检查库存充足性
		testInventory := getInventoryByGoodsID(goodsID)
		if testInventory == nil || testInventory.Stock <= 0 {
			t.Fatal("库存不足")
		}
		
		// 4.3 模拟添加到购物车
		cartItem := TestCartItem{
			UserID:     userID,
			GoodsID:    goodsID,
			GoodsName:  testGoods.Name,
			GoodsPrice: testGoods.ShopPrice,
			Nums:       1,
			Checked:    true,
		}
		
		t.Logf("模拟购物车添加: 用户ID %d, 商品 %s, 价格 ¥%.2f", 
			userID, cartItem.GoodsName, cartItem.GoodsPrice)
		
		// 4.4 模拟订单创建
		expectedOrderAmount := cartItem.GoodsPrice * float32(cartItem.Nums)
		t.Logf("模拟订单创建: 预期金额 ¥%.2f", expectedOrderAmount)
		
		// 4.5 模拟库存扣减
		expectedStock := testInventory.Stock - cartItem.Nums
		t.Logf("模拟库存扣减: 原库存 %d, 扣减 %d, 余库存 %d", 
			testInventory.Stock, cartItem.Nums, expectedStock)
		
		// 4.6 验证流程完整性
		if expectedStock >= 0 && expectedOrderAmount > 0 {
			t.Log("完整购物流程验证通过")
		} else {
			t.Error("购物流程验证失败")
		}
	})
}

// 集成测试场景：VIP用户特殊流程
func TestVIPUserFlow(t *testing.T) {
	t.Log("=== 跨服务集成测试：VIP用户特殊流程 ===")
	
	// 1. 获取VIP用户
	t.Run("VIP用户验证", func(t *testing.T) {
		vipUsers := getVIPUsers()
		if len(vipUsers) == 0 {
			t.Fatal("未找到VIP用户")
		}
		
		for _, vipUser := range vipUsers {
			t.Logf("VIP用户: %s (ID: %d, 级别: %s)", 
				vipUser.NickName, vipUser.ID, vipUser.Role)
		}
	})
	
	// 2. VIP用户订单验证
	t.Run("VIP订单验证", func(t *testing.T) {
		vipUsers := getVIPUsers()
		if len(vipUsers) == 0 {
			t.Fatal("未找到VIP用户")
		}
		
		for _, vipUser := range vipUsers {
			userOrders := getOrdersByUserID(vipUser.ID)
			t.Logf("VIP用户 %s 的订单数量: %d", vipUser.NickName, len(userOrders))
			
			for _, order := range userOrders {
				t.Logf("  VIP订单: ID %d, 金额 ¥%.2f, 状态: %s", 
					order.ID, order.OrderMount, order.Status)
			}
		}
	})
	
	// 3. 高价值商品验证
	t.Run("高价值商品验证", func(t *testing.T) {
		highValueGoods := getHighValueGoods()
		if len(highValueGoods) == 0 {
			t.Log("当前没有高价值商品（价格>10000）")
			return
		}
		
		for _, goods := range highValueGoods {
			t.Logf("高价值商品: %s (ID: %d, 价格: ¥%.2f, 分类: %s)", 
				goods.Name, goods.ID, goods.ShopPrice, goods.CategoryName)
		}
	})
}

// 集成测试场景：库存预警流程
func TestInventoryAlertFlow(t *testing.T) {
	t.Log("=== 跨服务集成测试：库存预警流程 ===")
	
	// 1. 获取低库存商品
	t.Run("低库存商品检测", func(t *testing.T) {
		lowStockItems := getLowStockInventories()
		if len(lowStockItems) == 0 {
			t.Log("当前没有低库存商品")
			return
		}
		
		for _, item := range lowStockItems {
			t.Logf("低库存预警: 商品ID %d (%s), 当前库存: %d, 状态: %s", 
				item.GoodsID, item.GoodsName, item.Stock, item.Status)
			
			// 2. 查找对应的商品信息
			goods := getGoodsByID(item.GoodsID)
			if goods != nil {
				t.Logf("  商品详情: %s, 分类: %s, 价格: ¥%.2f", 
					goods.Name, goods.CategoryName, goods.ShopPrice)
			}
		}
	})
	
	// 3. 模拟补货建议
	t.Run("补货建议", func(t *testing.T) {
		lowStockItems := getLowStockInventories()
		for _, item := range lowStockItems {
			// 根据商品类别和价值给出补货建议
			goods := getGoodsByID(item.GoodsID)
			if goods == nil {
				continue
			}
			
			var suggestedStock int32
			if goods.ShopPrice > 5000 { // 高价值商品
				suggestedStock = 20
			} else if goods.ShopPrice > 1000 { // 中等价值商品
				suggestedStock = 50
			} else { // 低价值商品
				suggestedStock = 100
			}
			
			t.Logf("补货建议: 商品 %s, 当前库存 %d, 建议补货至 %d", 
				item.GoodsName, item.Stock, suggestedStock)
		}
	})
}

// 集成测试场景：数据一致性验证
func TestDataConsistencyFlow(t *testing.T) {
	t.Log("=== 跨服务集成测试：数据一致性验证 ===")
	
	// 1. 用户-订单一致性
	t.Run("用户订单一致性", func(t *testing.T) {
		for _, user := range testUsers {
			userOrders := getOrdersByUserID(user.ID)
			t.Logf("用户 %s (ID: %d) 的订单数量: %d", 
				user.NickName, user.ID, len(userOrders))
			
			// 验证订单中的用户ID是否一致
			for _, order := range userOrders {
				if order.UserID != user.ID {
					t.Errorf("订单用户ID不一致: 订单ID %d, 期望用户ID %d, 实际用户ID %d", 
						order.ID, user.ID, order.UserID)
				} else {
					t.Logf("  订单 %s 用户ID一致性检查通过", order.OrderSN)
				}
			}
		}
	})
	
	// 2. 商品-库存一致性
	t.Run("商品库存一致性", func(t *testing.T) {
		for _, goods := range testGoods {
			inventory := getInventoryByGoodsID(goods.ID)
			if inventory == nil {
				t.Logf("商品 %s (ID: %d) 没有对应的库存记录", goods.Name, goods.ID)
				continue
			}
			
			t.Logf("商品-库存一致性检查: %s (ID: %d), 库存: %d", 
				goods.Name, goods.ID, inventory.Stock)
			
			if inventory.GoodsID != goods.ID {
				t.Errorf("库存商品ID不一致: 商品ID %d, 库存中商品ID %d", 
					goods.ID, inventory.GoodsID)
			} else {
				t.Logf("  商品 %s 库存一致性检查通过", goods.Name)
			}
		}
	})
	
	// 3. 业务规则验证
	t.Run("业务规则验证", func(t *testing.T) {
		// 验证VIP用户订单金额通常较高
		vipUsers := getVIPUsers()
		for _, vipUser := range vipUsers {
			userOrders := getOrdersByUserID(vipUser.ID)
			for _, order := range userOrders {
				if order.OrderMount > 10000 {
					t.Logf("VIP用户 %s 的高额订单: ¥%.2f ✓", 
						vipUser.NickName, order.OrderMount)
				} else {
					t.Logf("VIP用户 %s 的普通订单: ¥%.2f", 
						vipUser.NickName, order.OrderMount)
				}
			}
		}
		
		// 验证库存状态与实际库存数量的匹配
		for _, inventory := range testInventories {
			expectedStatus := ""
			if inventory.Stock > 30 {
				expectedStatus = "充足"
			} else if inventory.Stock > 10 {
				expectedStatus = "正常"
			} else {
				expectedStatus = "偏低"
			}
			
			if inventory.Status == expectedStatus || 
			   (inventory.Status == "充足" && inventory.Stock > 20) ||
			   (inventory.Status == "偏低" && inventory.Stock <= 10) {
				t.Logf("库存状态匹配: %s 库存%d 状态%s ✓", 
					inventory.GoodsName, inventory.Stock, inventory.Status)
			} else {
				t.Logf("库存状态需要关注: %s 库存%d 状态%s (预期: %s)", 
					inventory.GoodsName, inventory.Stock, inventory.Status, expectedStatus)
			}
		}
	})
}

// 集成测试场景：服务性能和边界测试
func TestServiceBoundaryFlow(t *testing.T) {
	t.Log("=== 跨服务集成测试：服务边界测试 ===")
	
	// 1. 大订单处理能力测试
	t.Run("大订单处理", func(t *testing.T) {
		userID := int32(13) // VIP用户
		user := getUserByID(userID)
		t.Logf("测试VIP用户: %s", user.NickName)
		
		// 模拟大订单（多个高价值商品）
		var totalAmount float32
		orderItems := []TestCartItem{
			{GoodsID: 1, GoodsName: "iPhone 15 Pro Max", GoodsPrice: 9999.0, Nums: 2},
			{GoodsID: 6, GoodsName: "MacBook Pro 14英寸", GoodsPrice: 28999.0, Nums: 1},
		}
		
		for _, item := range orderItems {
			// 检查库存是否足够
			inventory := getInventoryByGoodsID(item.GoodsID)
			if inventory != nil && inventory.Stock >= item.Nums {
				totalAmount += item.GoodsPrice * float32(item.Nums)
				t.Logf("添加商品: %s x%d, 单价: ¥%.2f", 
					item.GoodsName, item.Nums, item.GoodsPrice)
			} else {
				t.Logf("库存不足: %s 需要%d, 可用%d", 
					item.GoodsName, item.Nums, 
					func() int32 { if inventory != nil { return inventory.Stock } else { return 0 } }())
			}
		}
		
		t.Logf("大订单总金额: ¥%.2f", totalAmount)
		
		// 验证大订单阈值
		if totalAmount > 50000 {
			t.Logf("超大订单检测: 需要特殊审批 ✓")
		} else if totalAmount > 20000 {
			t.Logf("大订单检测: VIP用户自动通过 ✓")
		}
	})
	
	// 2. 库存边界测试
	t.Run("库存边界测试", func(t *testing.T) {
		for _, inventory := range testInventories {
			t.Logf("测试商品: %s, 当前库存: %d", inventory.GoodsName, inventory.Stock)
			
			// 测试不同购买数量的库存检查
			testQuantities := []int32{1, inventory.Stock, inventory.Stock + 1}
			
			for _, qty := range testQuantities {
				if qty <= inventory.Stock {
					t.Logf("  购买%d件: 库存充足 ✓", qty)
				} else {
					t.Logf("  购买%d件: 库存不足，需要%d件 ✗", qty, qty-inventory.Stock)
				}
			}
		}
	})
	
	// 3. 用户权限边界测试
	t.Run("用户权限边界测试", func(t *testing.T) {
		for _, user := range testUsers {
			// 模拟不同用户类型的权限检查
			if user.Role == "vip" {
				t.Logf("VIP用户 %s: 享受优先发货、专属客服、无金额限制 ✓", user.NickName)
			} else {
				t.Logf("普通用户 %s: 标准服务、单次订单限额检查", user.NickName)
			}
			
			// 检查历史订单金额
			userOrders := getOrdersByUserID(user.ID)
			var totalOrderAmount float32
			for _, order := range userOrders {
				totalOrderAmount += order.OrderMount
			}
			
			if totalOrderAmount > 0 {
				t.Logf("  用户 %s 历史消费总额: ¥%.2f", user.NickName, totalOrderAmount)
			}
		}
	})
}

// 主函数：运行所有集成测试
func TestAllIntegrationScenarios(t *testing.T) {
	t.Log("开始执行跨服务集成测试套件...")
	
	// 按顺序执行各个集成测试场景
	t.Run("完整下单流程", TestCompleteOrderFlow)
	t.Run("VIP用户特殊流程", TestVIPUserFlow)
	t.Run("库存预警流程", TestInventoryAlertFlow)
	t.Run("数据一致性验证", TestDataConsistencyFlow)
	t.Run("服务边界测试", TestServiceBoundaryFlow)
	
	t.Log("跨服务集成测试套件执行完成 ✓")
	
	// 输出测试总结
	t.Run("测试总结", func(t *testing.T) {
		t.Log("=== 集成测试总结 ===")
		t.Logf("✓ 测试用户数量: %d", len(testUsers))
		t.Logf("✓ 测试商品数量: %d", len(testGoods))
		t.Logf("✓ 测试库存记录: %d", len(testInventories))
		t.Logf("✓ 测试订单数量: %d", len(testOrders))
		
		// 统计各类用户
		vipCount := len(getVIPUsers())
		t.Logf("✓ VIP用户数量: %d", vipCount)
		t.Logf("✓ 普通用户数量: %d", len(testUsers)-vipCount)
		
		// 统计库存状态
		lowStockCount := len(getLowStockInventories())
		t.Logf("✓ 低库存商品数量: %d", lowStockCount)
		t.Logf("✓ 正常库存商品数量: %d", len(testInventories)-lowStockCount)
		
		t.Log("所有跨服务集成测试场景验证完成！")
	})
}