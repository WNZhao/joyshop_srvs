package tests

import (
	"context"
	"fmt"
	"order_srv/handler"
	"order_srv/proto"
	"testing"
)

func setupOrderService() *handler.OrderServiceServer {
	return &handler.OrderServiceServer{}
}

func TestCartItemAddAndList(t *testing.T) {
	t.Log(TestScenarios.CartAdd)

	initTestEnv(t)
	srv := setupOrderService()
	ctx := context.Background()

	// 使用预定义的测试数据添加购物车
	testItem := TestCartItemConstants{
		UserID:     GetRandomTestUser(),
		GoodsID:    GetRandomTestGoods(),
		GoodsName:  "测试商品",
		GoodsImage: "https://img.example.com/test.jpg",
		GoodsPrice: 99.99,
		Nums:       2,
		Checked:    true,
	}
	addReq := &proto.CartItemRequest{
		UserId:     testItem.UserID,
		GoodsId:    testItem.GoodsID,
		GoodsName:  testItem.GoodsName,
		GoodsImage: testItem.GoodsImage,
		GoodsPrice: testItem.GoodsPrice,
		Nums:       testItem.Nums,
		Checked:    testItem.Checked,
	}
	_, err := srv.CartItemAdd(ctx, addReq)
	if err != nil {
		t.Fatalf("添加购物车失败: %v", err)
	}
	t.Logf("添加购物车成功 - 用户ID: %d, 商品: %s, 数量: %d",
		testItem.UserID, testItem.GoodsName, testItem.Nums)

	// 查询购物车
	listReq := &proto.UserInfo{Id: addReq.UserId}
	listResp, err := srv.CartItemList(ctx, listReq)
	if err != nil {
		t.Fatalf("查询购物车失败: %v", err)
	}

	t.Logf("用户ID %d 的购物车商品数量: %d", addReq.UserId, len(listResp.CartItems))
	for _, item := range listResp.CartItems {
		t.Logf("购物车商品 - ID: %d, 用户ID: %d, 商品ID: %d, 数量: %d, 已选中: %v",
			item.Id, item.UserId, item.GoodsId, item.Nums, item.Checked)
	}
}

func TestCartItemUpdateAndDelete(t *testing.T) {
	t.Log(TestScenarios.CartUpdate)

	initTestEnv(t)
	srv := setupOrderService()
	ctx := context.Background()

	// 使用预定义的测试数据
	testCartItem := TestCartItems[0] // 用户张三的iPhone购物车
	updateReq := &proto.CartItemRequest{
		Id:      testCartItem.ID,
		UserId:  testCartItem.UserID,
		GoodsId: testCartItem.GoodsID,
		Nums:    testCartItem.Nums + 1, // 数量+1
		Checked: !testCartItem.Checked, // 切换选中状态
	}
	_, err := srv.CartItemUpdate(ctx, updateReq)
	if err != nil {
		t.Fatalf("更新购物车失败: %v", err)
	}
	t.Logf("更新购物车成功 - 用户ID: %d, 商品ID: %d, 新数量: %d, 选中状态: %v",
		updateReq.UserId, updateReq.GoodsId, updateReq.Nums, updateReq.Checked)

	_, err = srv.CartItemDelete(ctx, updateReq)
	if err != nil {
		t.Fatalf("删除购物车失败: %v", err)
	}
	t.Logf("删除购物车成功 - 用户ID: %d, 商品ID: %d",
		updateReq.UserId, updateReq.GoodsId)
}

func TestOrderCreateAndList(t *testing.T) {
	t.Log(TestScenarios.OrderCreate)

	initTestEnv(t)
	srv := setupOrderService()
	ctx := context.Background()

	// 使用预定义的测试用户创建订单
	testUser := GetRandomTestUser() // 获取测试用户
	orderReq := &proto.OrderRequest{
		UserId:  testUser,
		Address: "测试地址 - 上海市浦东新区",
		Name:    "测试用户",
		Mobile:  "13800138000",
		Post:    "测试订单，请尽快发货",
	}
	orderResp, err := srv.OrderCreate(ctx, orderReq)
	if err != nil {
		t.Fatalf("创建订单失败: %v", err)
	}
	t.Logf("创建订单成功 - 用户ID: %d, 订单ID: %d, 订单号: %s",
		orderReq.UserId, orderResp.Id, orderResp.OrderSn)

	listReq := &proto.OrderFilterRequest{UserId: orderReq.UserId, Page: 1, PageSize: 10}
	listResp, err := srv.OrderList(ctx, listReq)
	if err != nil {
		t.Fatalf("查询订单列表失败: %v", err)
	}

	t.Logf("用户ID %d 的订单数量: %d", orderReq.UserId, listResp.Total)
	for _, order := range listResp.Data {
		t.Logf("订单信息 - ID: %d, 订单号: %s, 状态: %s, 金额: %.2f",
			order.Id, order.OrderSn, order.Status, order.Total)
	}
}

func TestOrderDetailUpdateDelete(t *testing.T) {
	t.Log(TestScenarios.OrderDetail)

	initTestEnv(t)
	srv := setupOrderService()
	ctx := context.Background()

	// 使用预定义的测试订单数据
	testOrder := GetTestOrderByID(1) // 使用第一个测试订单
	if testOrder == nil {
		t.Skip("未找到测试订单数据")
		return
	}

	orderReq := &proto.OrderRequest{Id: testOrder.ID, UserId: testOrder.UserID}
	detailResp, err := srv.OrderDetail(ctx, orderReq)
	if err != nil {
		t.Fatalf("查询订单详情失败: %v", err)
	}
	t.Logf("订单详情查询成功 - 订单ID: %d, 用户ID: %d",
		detailResp.OrderInfo.Id, detailResp.OrderInfo.UserId)

	updateReq := &proto.OrderStatus{Id: testOrder.ID, OrderSn: testOrder.OrderSN, Status: OrderStatus.TradeSuccess}
	_, err = srv.OrderUpdate(ctx, updateReq)
	if err != nil {
		t.Fatalf("更新订单状态失败: %v", err)
	}
	t.Logf("订单状态更新成功 - 订单ID: %d, 新状态: %s",
		testOrder.ID, OrderStatus.TradeSuccess)

	delReq := &proto.OrderDelRequest{Id: testOrder.ID, UserId: testOrder.UserID}
	_, err = srv.OrderDelete(ctx, delReq)
	if err != nil {
		t.Fatalf("删除订单失败: %v", err)
	}
	t.Logf("订单删除成功 - 订单ID: %d", testOrder.ID)
}

// 测试预定义购物车数据
func TestPredefinedCartItems(t *testing.T) {
	t.Log("预定义购物车数据测试")

	initTestEnv(t)
	srv := setupOrderService()
	ctx := context.Background()

	// 测试不同用户的购物车
	testUsers := []int32{2, 3, 4, 13} // 张三、李四、王五、VIP刘总

	for _, userID := range testUsers {
		t.Run(fmt.Sprintf("用户%d购物车", userID), func(t *testing.T) {
			// 获取用户的预定义购物车数据
			userCartItems := GetTestCartItemsByUserID(userID)

			if len(userCartItems) == 0 {
				t.Skipf("用户ID %d 没有预定义购物车数据", userID)
				return
			}

			t.Logf("用户ID %d 的预定义购物车商品数量: %d", userID, len(userCartItems))

			for _, item := range userCartItems {
				// 添加到购物车
				addReq := &proto.CartItemRequest{
					UserId:     item.UserID,
					GoodsId:    item.GoodsID,
					GoodsName:  item.GoodsName,
					GoodsImage: item.GoodsImage,
					GoodsPrice: item.GoodsPrice,
					Nums:       item.Nums,
					Checked:    item.Checked,
				}

				_, err := srv.CartItemAdd(ctx, addReq)
				if err != nil {
					t.Errorf("添加购物车商品失败: %v", err)
					continue
				}

				t.Logf("添加购物车成功 - 商品: %s, 价格: ¥%.2f, 数量: %d, 已选中: %v",
					item.GoodsName, item.GoodsPrice, item.Nums, item.Checked)
			}

			// 查询购物车验证
			listReq := &proto.UserInfo{Id: userID}
			listResp, err := srv.CartItemList(ctx, listReq)
			if err != nil {
				t.Errorf("查询购物车失败: %v", err)
				return
			}

			// 计算购物车总价
			totalAmount := CalculateCartTotal(userCartItems)
			t.Logf("用户ID %d 购物车总金额: ¥%.2f, 商品种类: %d",
				userID, totalAmount, len(listResp.CartItems))
		})
	}
}

// 测试订单状态流转
func TestOrderStatusFlow(t *testing.T) {
	t.Log(TestScenarios.PaymentFlow)

	initTestEnv(t)
	srv := setupOrderService()
	ctx := context.Background()

	// 测试预定义订单的状态流转
	testOrders := []struct {
		OrderID    int32
		UserID     int32
		InitStatus string
		NextStatus string
	}{
		{OrderID: 7, UserID: 6, InitStatus: OrderStatus.WaitBuyerPay, NextStatus: OrderStatus.Paying},
		{OrderID: 8, UserID: 7, InitStatus: OrderStatus.WaitBuyerPay, NextStatus: OrderStatus.TradeSuccess},
		{OrderID: 9, UserID: 8, InitStatus: OrderStatus.Paying, NextStatus: OrderStatus.TradeSuccess},
	}

	for _, testOrder := range testOrders {
		t.Run(fmt.Sprintf("订单%d状态流转", testOrder.OrderID), func(t *testing.T) {
			// 更新订单状态
			updateReq := &proto.OrderStatus{
				Id:      testOrder.OrderID,
				OrderSn: GenerateOrderSN(testOrder.UserID),
				Status:  testOrder.NextStatus,
			}

			_, err := srv.OrderUpdate(ctx, updateReq)
			if err != nil {
				t.Errorf("更新订单状态失败: %v", err)
				return
			}

			t.Logf("订单状态更新成功 - 订单ID: %d, 从 %s 更新为 %s",
				testOrder.OrderID, testOrder.InitStatus, testOrder.NextStatus)

			// 验证订单详情
			detailReq := &proto.OrderRequest{Id: testOrder.OrderID, UserId: testOrder.UserID}
			detailResp, err := srv.OrderDetail(ctx, detailReq)
			if err != nil {
				t.Errorf("查询订单详情失败: %v", err)
				return
			}

			if detailResp.OrderInfo.Status != testOrder.NextStatus {
				t.Errorf("订单状态不匹配: 期望 %s, 实际 %s",
					testOrder.NextStatus, detailResp.OrderInfo.Status)
			}
		})
	}
}

// 测试VIP用户订单
func TestVIPUserOrders(t *testing.T) {
	t.Log(TestScenarios.VIPService)

	initTestEnv(t)
	srv := setupOrderService()
	ctx := context.Background()

	// 获取VIP用户订单
	vipOrders := GetVIPUserOrders()
	t.Logf("VIP用户订单数量: %d", len(vipOrders))

	for _, vipOrder := range vipOrders {
		t.Run(fmt.Sprintf("VIP订单_%d", vipOrder.ID), func(t *testing.T) {
			// 查询VIP订单详情
			detailReq := &proto.OrderRequest{Id: vipOrder.ID, UserId: vipOrder.UserID}
			detailResp, err := srv.OrderDetail(ctx, detailReq)
			if err != nil {
				t.Errorf("查询VIP订单详情失败: %v", err)
				return
			}

			// 验证VIP订单特征（高金额）
			if detailResp.OrderInfo.Total < 10000 { // VIP订单通常金额较高
				t.Logf("注意：VIP订单金额较低: ¥%.2f", detailResp.OrderInfo.Total)
			}

			t.Logf("VIP订单详情 - 订单号: %s, 用户ID: %d, 金额: ¥%.2f, 状态: %s",
				detailResp.OrderInfo.OrderSn, detailResp.OrderInfo.UserId, detailResp.OrderInfo.Total, detailResp.OrderInfo.Status)

			// VIP订单应该有特殊备注
			if len(detailResp.OrderInfo.Post) > 0 && (len(detailResp.OrderInfo.Post) < 3 || detailResp.OrderInfo.Post[:3] != "VIP") {
				t.Logf("VIP订单备注: %s", detailResp.OrderInfo.Post)
			}
		})
	}
}

// 测试订单创建场景
func TestOrderCreationScenarios(t *testing.T) {
	t.Log("订单创建场景测试")

	initTestEnv(t)
	srv := setupOrderService()
	ctx := context.Background()

	// 测试预定义的订单创建场景
	for i, scenario := range TestOrderCreationScenarios {
		t.Run(fmt.Sprintf("创建场景_%d", i+1), func(t *testing.T) {
			// 先添加购物车商品
			for _, item := range scenario.GoodsItems {
				addReq := &proto.CartItemRequest{
					UserId:     scenario.UserID,
					GoodsId:    item.GoodsID,
					GoodsName:  item.GoodsName,
					GoodsImage: item.GoodsImage,
					GoodsPrice: item.GoodsPrice,
					Nums:       item.Nums,
					Checked:    true, // 设为选中状态
				}

				_, err := srv.CartItemAdd(ctx, addReq)
				if err != nil {
					t.Errorf("添加购物车商品失败: %v", err)
				}
			}

			// 创建订单
			orderReq := &proto.OrderRequest{
				UserId:  scenario.UserID,
				Address: scenario.Address,
				Name:    scenario.SignerName,
				Mobile:  scenario.SignerMobile,
				Post:    scenario.Post,
			}

			orderResp, err := srv.OrderCreate(ctx, orderReq)

			// 验证结果
			if scenario.ExpectedResult == "success" {
				if err != nil {
					t.Errorf("订单创建应该成功但失败了: %v", err)
				} else {
					t.Logf("订单创建成功 - 订单ID: %d, 订单号: %s",
						orderResp.Id, orderResp.OrderSn)
				}
			} else if scenario.ExpectedResult == "insufficient_stock" {
				if err == nil {
					t.Error("库存不足时订单创建应该失败但成功了")
				} else {
					t.Logf("库存不足订单创建正确失败: %v", err)
				}
			}
		})
	}
}

// 测试购物车批量操作
func TestBatchCartOperations(t *testing.T) {
	t.Log("购物车批量操作测试")

	initTestEnv(t)
	srv := setupOrderService()
	ctx := context.Background()

	testUserID := int32(2) // 使用张三的用户ID

	// 批量添加购物车商品
	batchItems := []TestCartItemConstants{
		{UserID: testUserID, GoodsID: 1, GoodsName: "iPhone 15", GoodsPrice: 9999.0, Nums: 1, Checked: true},
		{UserID: testUserID, GoodsID: 2, GoodsName: "HUAWEI Mate", GoodsPrice: 7499.0, Nums: 1, Checked: true},
		{UserID: testUserID, GoodsID: 3, GoodsName: "小米14", GoodsPrice: 6499.0, Nums: 2, Checked: false},
	}

	t.Run("批量添加商品", func(t *testing.T) {
		for _, item := range batchItems {
			addReq := &proto.CartItemRequest{
				UserId:     item.UserID,
				GoodsId:    item.GoodsID,
				GoodsName:  item.GoodsName,
				GoodsPrice: item.GoodsPrice,
				Nums:       item.Nums,
				Checked:    item.Checked,
			}

			_, err := srv.CartItemAdd(ctx, addReq)
			if err != nil {
				t.Errorf("批量添加商品失败: %v", err)
				continue
			}

			t.Logf("批量添加成功 - 商品: %s, 数量: %d, 选中: %v",
				item.GoodsName, item.Nums, item.Checked)
		}
	})

	t.Run("验证购物车内容", func(t *testing.T) {
		listReq := &proto.UserInfo{Id: testUserID}
		listResp, err := srv.CartItemList(ctx, listReq)
		if err != nil {
			t.Errorf("查询购物车失败: %v", err)
			return
		}

		totalItems := len(listResp.CartItems)
		checkedItems := 0
		totalAmount := float32(0)

		// 由于购物车返回的是ShopCartInfoResponse，不包含价格信息
		// 所以使用预定义的测试数据来计算总金额
		for _, item := range batchItems {
			if item.Checked {
				checkedItems++
				totalAmount += item.GoodsPrice * float32(item.Nums)
			}
		}

		t.Logf("购物车统计 - 总商品: %d, 已选中: %d, 已选金额: ¥%.2f",
			totalItems, checkedItems, totalAmount)
	})
}
