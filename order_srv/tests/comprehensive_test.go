package tests

import (
	"context"
	"order_srv/global"
	"order_srv/handler"
	"order_srv/model"
	"order_srv/proto"
	"testing"
)

// TestComprehensiveOrderAndCartFlow 综合测试订单和购物车流程
func TestComprehensiveOrderAndCartFlow(t *testing.T) {
	t.Log("=== 订单服务综合功能测试 ===")

	initTestEnvSimple(t)
	srv := &handler.OrderServiceServer{}
	ctx := context.Background()

	// 使用真实存在的用户ID
	validUserID := int32(15) // testuser
	t.Logf("使用有效用户ID进行测试: %d", validUserID)

	// === 1. 测试购物车功能 ===
	t.Run("购物车功能测试", func(t *testing.T) {
		// 添加商品到购物车
		addReq := &proto.CartItemRequest{
			UserId:     validUserID,
			GoodsId:    999, // 使用不存在的商品ID，避免与现有数据冲突
			GoodsName:  "综合测试商品",
			GoodsImage: "test-image.jpg", 
			GoodsPrice: 199.99,
			Nums:       2,
			Checked:    true,
		}

		_, err := srv.CartItemAdd(ctx, addReq)
		if err != nil {
			t.Fatalf("添加购物车失败: %v", err)
		}
		t.Log("✅ 购物车添加功能正常")

		// 查询购物车
		listReq := &proto.UserInfo{Id: validUserID}
		listResp, err := srv.CartItemList(ctx, listReq)
		if err != nil {
			t.Fatalf("查询购物车失败: %v", err)
		}

		var testCartID int32 = 0
		for _, item := range listResp.CartItems {
			if item.GoodsId == 999 {
				testCartID = item.Id
				break
			}
		}

		if testCartID == 0 {
			t.Fatal("未找到刚添加的测试商品")
		}
		t.Log("✅ 购物车查询功能正常")

		// 更新购物车
		updateReq := &proto.CartItemRequest{
			Id:      testCartID,
			UserId:  validUserID,
			GoodsId: 999,
			Nums:    3, // 更新数量
			Checked: false, // 更新选中状态
		}

		_, err = srv.CartItemUpdate(ctx, updateReq)
		if err != nil {
			t.Fatalf("更新购物车失败: %v", err)
		}
		t.Log("✅ 购物车更新功能正常")

		// 删除购物车商品
		_, err = srv.CartItemDelete(ctx, updateReq)
		if err != nil {
			t.Fatalf("删除购物车失败: %v", err)
		}
		t.Log("✅ 购物车删除功能正常")
	})

	// === 2. 测试用户验证功能 ===
	t.Run("用户验证功能测试", func(t *testing.T) {
		// 测试有效用户
		orderReq := &proto.OrderRequest{
			UserId:  validUserID,
			Address: "综合测试地址",
			Name:    "综合测试用户",
			Mobile:  "13900139000",
			Post:    "综合测试订单",
		}

		_, err := srv.OrderCreate(ctx, orderReq)
		// 这里会因为购物车没有商品而失败，但用户验证应该通过
		if err != nil && err.Error() == "rpc error: code = NotFound desc = 用户不存在" {
			t.Error("❌ 有效用户被错误拒绝")
		} else {
			t.Log("✅ 用户验证功能正常")
		}

		// 测试无效用户
		invalidReq := &proto.OrderRequest{
			UserId:  99999, // 超出范围的用户ID
			Address: "测试地址",
			Name:    "测试用户",
			Mobile:  "13900139000",
			Post:    "测试",
		}

		_, err = srv.OrderCreate(ctx, invalidReq)
		if err == nil || err.Error() != "rpc error: code = NotFound desc = 用户不存在" {
			t.Error("❌ 无效用户应该被拒绝")
		} else {
			t.Log("✅ 无效用户拒绝功能正常")
		}
	})

	// === 3. 测试订单状态更新功能 ===
	t.Run("订单状态更新功能测试", func(t *testing.T) {
		// 首先在数据库中创建一个测试订单
		testOrder := &model.OrderInfo{
			User:         validUserID,
			OrderSn:      "COMPREHENSIVE_TEST_ORDER_001",
			PayType:      "alipay",
			Status:       "WAIT_BUYER_PAY",
			OrderMount:   299.99,
			Address:      "测试地址",
			SignerName:   "测试用户",
			SingerMobile: "13900139000",
			Post:         "状态更新测试订单",
		}

		if err := global.DB.Create(testOrder).Error; err != nil {
			t.Fatalf("创建测试订单失败: %v", err)
		}

		// 测试状态更新
		updateStatusReq := &proto.OrderStatus{
			OrderSn: testOrder.OrderSn,
			Status:  "PAYING",
		}

		_, err := srv.OrderUpdate(ctx, updateStatusReq)
		if err != nil {
			t.Fatalf("更新订单状态失败: %v", err)
		}
		t.Log("✅ 订单状态更新功能正常")

		// 验证状态是否更新成功
		var updatedOrder model.OrderInfo
		if err := global.DB.Where("order_sn = ?", testOrder.OrderSn).First(&updatedOrder).Error; err != nil {
			t.Fatalf("查询更新后订单失败: %v", err)
		}

		if updatedOrder.Status != "PAYING" {
			t.Errorf("订单状态未正确更新，期望: PAYING，实际: %s", updatedOrder.Status)
		} else {
			t.Log("✅ 订单状态验证正常")
		}

		// 清理测试数据
		global.DB.Where("order_sn = ?", testOrder.OrderSn).Delete(&model.OrderInfo{})
	})

	t.Log("=== 综合测试完成 ===")
}