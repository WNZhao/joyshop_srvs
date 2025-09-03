/*
 * @Description: 订单服务测试数据常量
 * 基于SQL脚本中的真实测试数据定义
 */
package tests

// TestCartItemConstants 测试购物车数据常量
type TestCartItemConstants struct {
	ID         int32
	UserID     int32
	GoodsID    int32
	GoodsName  string
	GoodsImage string
	GoodsPrice float32
	Nums       int32
	Checked    bool
}

// TestOrderConstants 测试订单数据常量
type TestOrderConstants struct {
	ID          int32
	UserID      int32
	OrderSN     string
	OrderMount  float32
	Status      string
	PayType     string
	Address     string
	SignerName  string
	SignerMobile string
	Post        string
}

// TestOrderGoodsConstants 测试订单商品数据常量
type TestOrderGoodsConstants struct {
	ID         int32
	OrderID    int32
	GoodsID    int32
	GoodsName  string
	GoodsImage string
	GoodsPrice float32
	Nums       int32
}

// 预定义的测试购物车数据（与SQL脚本保持一致）
var TestCartItems = []TestCartItemConstants{
	// 用户2（张三）的购物车
	{ID: 1, UserID: 2, GoodsID: 1, GoodsName: "iPhone 15 Pro Max 256GB", GoodsImage: "https://img.example.com/iphone15promax.jpg", GoodsPrice: 9999.0, Nums: 1, Checked: true},
	{ID: 2, UserID: 2, GoodsID: 7, GoodsName: "Nike Air Max 270 男款跑鞋", GoodsImage: "https://img.example.com/nike_airmax270.jpg", GoodsPrice: 999.0, Nums: 2, Checked: true},
	{ID: 3, UserID: 2, GoodsID: 9, GoodsName: "小米手环8 标准版", GoodsImage: "https://img.example.com/mi_band8.jpg", GoodsPrice: 249.0, Nums: 1, Checked: false},

	// 用户3（李四）的购物车
	{ID: 4, UserID: 3, GoodsID: 2, GoodsName: "HUAWEI Mate 60 Pro", GoodsImage: "https://img.example.com/huawei_mate60pro.jpg", GoodsPrice: 7499.0, Nums: 1, Checked: true},
	{ID: 5, UserID: 3, GoodsID: 8, GoodsName: "Adidas Ultraboost 22 跑鞋", GoodsImage: "https://img.example.com/adidas_ultraboost22.jpg", GoodsPrice: 1199.0, Nums: 1, Checked: true},

	// 用户4（王五）的购物车
	{ID: 6, UserID: 4, GoodsID: 3, GoodsName: "小米14 Ultra", GoodsImage: "https://img.example.com/mi14ultra.jpg", GoodsPrice: 6499.0, Nums: 1, Checked: true},
	{ID: 7, UserID: 4, GoodsID: 10, GoodsName: "UNIQLO 圆领T恤", GoodsImage: "https://img.example.com/uniqlo_tshirt.jpg", GoodsPrice: 69.0, Nums: 3, Checked: true},

	// VIP用户（刘总）的购物车
	{ID: 8, UserID: 13, GoodsID: 6, GoodsName: "MacBook Pro 14英寸", GoodsImage: "https://img.example.com/macbook_pro14.jpg", GoodsPrice: 28999.0, Nums: 1, Checked: true},
	{ID: 9, UserID: 13, GoodsID: 1, GoodsName: "iPhone 15 Pro Max 256GB", GoodsImage: "https://img.example.com/iphone15promax.jpg", GoodsPrice: 9999.0, Nums: 2, Checked: true},
}

// 预定义的测试订单数据（与SQL脚本保持一致）
var TestOrders = []TestOrderConstants{
	// 已完成订单
	{ID: 1, UserID: 13, OrderSN: "20240101103000001300130001", OrderMount: 48997.0, Status: "TRADE_FINISHED", PayType: "alipay", Address: "北京市朝阳区国贸CBD中心A座2501", SignerName: "VIP刘总", SignerMobile: "18888888880", Post: "VIP用户，请优先发货"},
	{ID: 2, UserID: 2, OrderSN: "20240102143000000200020002", OrderMount: 10997.0, Status: "TRADE_FINISHED", PayType: "wechat", Address: "上海市浦东新区陆家嘴金茂大厦88楼", SignerName: "张三", SignerMobile: "18782222220", Post: "工作日上午送达"},

	// 支付成功订单
	{ID: 3, UserID: 3, OrderSN: "20240103160000000300030003", OrderMount: 8698.0, Status: "TRADE_SUCCESS", PayType: "alipay", Address: "广州市天河区珠江新城CBD广州塔旁", SignerName: "李四", SignerMobile: "18782222221", Post: "请联系收件人预约时间"},
	{ID: 4, UserID: 4, OrderSN: "20240104111500000400040004", OrderMount: 6706.0, Status: "TRADE_SUCCESS", PayType: "wechat", Address: "深圳市南山区科技园腾讯大厦", SignerName: "王五", SignerMobile: "18782222222", Post: "上班时间请送到前台"},
	{ID: 5, UserID: 5, OrderSN: "20240105094500000500050005", OrderMount: 1199.0, Status: "TRADE_SUCCESS", PayType: "alipay", Address: "成都市锦江区春熙路步行街", SignerName: "赵六", SignerMobile: "18782222223", Post: "快递柜收取"},
	{ID: 6, UserID: 14, OrderSN: "20240106142000001400140006", OrderMount: 28999.0, Status: "TRADE_SUCCESS", PayType: "credit_card", Address: "杭州市西湖区文三路阿里巴巴总部", SignerName: "VIP王总", SignerMobile: "18888888881", Post: "VIP服务，请确保包装完好"},

	// 待支付订单
	{ID: 7, UserID: 6, OrderSN: "20240107173000000600060007", OrderMount: 1248.0, Status: "WAIT_BUYER_PAY", PayType: "alipay", Address: "西安市雁塔区高新区软件园", SignerName: "孙七", SignerMobile: "18782222224", Post: "周末在家收取"},
	{ID: 8, UserID: 7, OrderSN: "20240108201500000700070008", OrderMount: 6499.0, Status: "WAIT_BUYER_PAY", PayType: "wechat", Address: "武汉市洪山区光谷软件园", SignerName: "周八", SignerMobile: "18782222225", Post: "请提前电话联系"},

	// 支付中订单
	{ID: 9, UserID: 8, OrderSN: "20240109112000000800080009", OrderMount: 999.0, Status: "PAYING", PayType: "alipay", Address: "南京市鼓楼区新街口商圈", SignerName: "吴九", SignerMobile: "18782222226", Post: "门卫代收"},

	// 已关闭订单
	{ID: 10, UserID: 9, OrderSN: "20240110165000000900090010", OrderMount: 249.0, Status: "TRADE_CLOSED", PayType: "wechat", Address: "长沙市岳麓区橘子洲头", SignerName: "郑十", SignerMobile: "18782222227", Post: "用户取消"},
}

// 订单状态常量
var OrderStatus = struct {
	WaitBuyerPay  string
	Paying        string
	TradeSuccess  string
	TradeFinished string
	TradeClosed   string
}{
	WaitBuyerPay:  "WAIT_BUYER_PAY",  // 待支付
	Paying:        "PAYING",          // 支付中
	TradeSuccess:  "TRADE_SUCCESS",   // 支付成功
	TradeFinished: "TRADE_FINISHED",  // 交易完成
	TradeClosed:   "TRADE_CLOSED",    // 交易关闭
}

// 支付方式常量
var PaymentType = struct {
	Alipay     string
	Wechat     string
	CreditCard string
	UnionPay   string
}{
	Alipay:     "alipay",
	Wechat:     "wechat",
	CreditCard: "credit_card",
	UnionPay:   "union_pay",
}

// GetTestCartItemsByUserID 根据用户ID获取购物车商品
func GetTestCartItemsByUserID(userID int32) []TestCartItemConstants {
	var items []TestCartItemConstants
	for _, item := range TestCartItems {
		if item.UserID == userID {
			items = append(items, item)
		}
	}
	return items
}

// GetCheckedCartItemsByUserID 根据用户ID获取已选中的购物车商品
func GetCheckedCartItemsByUserID(userID int32) []TestCartItemConstants {
	var items []TestCartItemConstants
	for _, item := range TestCartItems {
		if item.UserID == userID && item.Checked {
			items = append(items, item)
		}
	}
	return items
}

// GetTestOrdersByUserID 根据用户ID获取订单
func GetTestOrdersByUserID(userID int32) []TestOrderConstants {
	var orders []TestOrderConstants
	for _, order := range TestOrders {
		if order.UserID == userID {
			orders = append(orders, order)
		}
	}
	return orders
}

// GetTestOrdersByStatus 根据状态获取订单
func GetTestOrdersByStatus(status string) []TestOrderConstants {
	var orders []TestOrderConstants
	for _, order := range TestOrders {
		if order.Status == status {
			orders = append(orders, order)
		}
	}
	return orders
}

// GetTestOrderByID 根据ID获取订单
func GetTestOrderByID(id int32) *TestOrderConstants {
	for i := range TestOrders {
		if TestOrders[i].ID == id {
			return &TestOrders[i]
		}
	}
	return nil
}

// GetVIPUserOrders 获取VIP用户订单
func GetVIPUserOrders() []TestOrderConstants {
	var orders []TestOrderConstants
	vipUserIDs := []int32{13, 14} // VIP用户ID
	for _, order := range TestOrders {
		for _, vipID := range vipUserIDs {
			if order.UserID == vipID {
				orders = append(orders, order)
				break
			}
		}
	}
	return orders
}

// TestOrderCreation 订单创建测试数据
type TestOrderCreation struct {
	UserID      int32
	GoodsItems  []TestCartItemConstants
	Address     string
	SignerName  string
	SignerMobile string
	Post        string
	ExpectedResult string // "success" or "insufficient_stock" or "error"
}

// 预定义的订单创建测试场景
var TestOrderCreationScenarios = []TestOrderCreation{
	// 正常创建订单场景
	{
		UserID: 2,
		GoodsItems: []TestCartItemConstants{
			{GoodsID: 9, GoodsName: "小米手环8", GoodsPrice: 249.0, Nums: 1},
		},
		Address:     "测试地址1",
		SignerName:  "张三",
		SignerMobile: "18782222220",
		Post:        "测试订单",
		ExpectedResult: "success",
	},
	// 库存不足场景
	{
		UserID: 3,
		GoodsItems: []TestCartItemConstants{
			{GoodsID: 5, GoodsName: "vivo X100 Pro", GoodsPrice: 4999.0, Nums: 10}, // 库存只有5，购买10应该失败
		},
		Address:     "测试地址2",
		SignerName:  "李四",
		SignerMobile: "18782222221",
		Post:        "测试库存不足",
		ExpectedResult: "insufficient_stock",
	},
	// VIP用户大额订单
	{
		UserID: 13,
		GoodsItems: []TestCartItemConstants{
			{GoodsID: 1, GoodsName: "iPhone 15 Pro Max", GoodsPrice: 9999.0, Nums: 2},
			{GoodsID: 6, GoodsName: "MacBook Pro 14英寸", GoodsPrice: 28999.0, Nums: 1},
		},
		Address:     "VIP专用地址",
		SignerName:  "VIP刘总",
		SignerMobile: "18888888880",
		Post:        "VIP用户，请优先处理",
		ExpectedResult: "success",
	},
}

// TestScenarios 测试场景常量
var TestScenarios = struct {
	CartAdd       string
	CartList      string
	CartUpdate    string
	CartDelete    string
	OrderCreate   string
	OrderList     string
	OrderDetail   string
	OrderUpdate   string
	OrderDelete   string
	PaymentFlow   string
	VIPService    string
}{
	CartAdd:     "添加购物车测试",
	CartList:    "查询购物车测试",
	CartUpdate:  "更新购物车测试",
	CartDelete:  "删除购物车测试",
	OrderCreate: "创建订单测试",
	OrderList:   "查询订单列表测试",
	OrderDetail: "查询订单详情测试",
	OrderUpdate: "更新订单状态测试",
	OrderDelete: "删除订单测试",
	PaymentFlow: "支付流程测试",
	VIPService:  "VIP用户服务测试",
}

// GetRandomTestUser 获取随机测试用户ID
func GetRandomTestUser() int32 {
	// 从普通用户中选择（ID: 2-12）
	return 2 // 简单返回张三的ID
}

// GetRandomTestGoods 获取随机测试商品ID
func GetRandomTestGoods() int32 {
	// 返回库存充足的商品
	return 1 // iPhone 15 Pro Max，库存充足
}

// CalculateCartTotal 计算购物车总金额
func CalculateCartTotal(items []TestCartItemConstants) float32 {
	var total float32 = 0
	for _, item := range items {
		if item.Checked {
			total += item.GoodsPrice * float32(item.Nums)
		}
	}
	return total
}

// GenerateOrderSN 生成测试订单号
func GenerateOrderSN(userID int32) string {
	// 使用简单的格式：日期 + 用户ID + 随机数
	return "TEST2024010100000001"  // 测试专用订单号格式
}