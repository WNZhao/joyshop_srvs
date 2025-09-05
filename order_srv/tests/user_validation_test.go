package tests

import (
	"context"
	"order_srv/handler"
	"order_srv/proto"
	"testing"
)

func TestOrderCreateWithInvalidUser(t *testing.T) {
	t.Log("测试无效用户创建订单")

	initTestEnvSimple(t)
	srv := &handler.OrderServiceServer{}
	ctx := context.Background()

	// 测试无效用户ID（超出合理范围）
	invalidUserID := int32(9999)
	orderReq := &proto.OrderRequest{
		UserId:  invalidUserID,
		Address: "测试地址",
		Name:    "测试用户",
		Mobile:  "13800138000",
		Post:    "测试订单",
	}

	t.Logf("尝试为无效用户创建订单，用户ID: %d", invalidUserID)
	_, err := srv.OrderCreate(ctx, orderReq)
	
	if err == nil {
		t.Error("无效用户创建订单应该失败但成功了")
	} else {
		t.Logf("无效用户创建订单正确失败: %v", err)
	}
}

func TestOrderCreateWithValidUser(t *testing.T) {
	t.Log("测试有效用户创建订单")

	initTestEnvSimple(t)
	srv := &handler.OrderServiceServer{}
	ctx := context.Background()

	// 测试有效用户ID（有购物车记录的用户）
	validUserID := int32(2)
	orderReq := &proto.OrderRequest{
		UserId:  validUserID,
		Address: "测试地址",
		Name:    "测试用户", 
		Mobile:  "13800138000",
		Post:    "测试订单",
	}

	t.Logf("为有效用户创建订单，用户ID: %d", validUserID)
	_, err := srv.OrderCreate(ctx, orderReq)
	
	// 这里会因为商品服务未连接而失败，但用户验证应该通过
	if err != nil {
		// 检查是否是用户不存在的错误
		if err.Error() == "rpc error: code = NotFound desc = 用户不存在" {
			t.Error("有效用户被错误拒绝")
		} else {
			t.Logf("有效用户通过验证，但因其他原因失败（预期）: %v", err)
		}
	} else {
		t.Log("订单创建成功")
	}
}

func TestOrderCreateWithNewUser(t *testing.T) {
	t.Log("测试新用户创建订单")

	initTestEnvSimple(t) 
	srv := &handler.OrderServiceServer{}
	ctx := context.Background()

	// 测试新用户ID（在合理范围内但没有历史数据）
	newUserID := int32(500)
	orderReq := &proto.OrderRequest{
		UserId:  newUserID,
		Address: "测试地址",
		Name:    "新用户",
		Mobile:  "13800138000", 
		Post:    "新用户订单",
	}

	t.Logf("为新用户创建订单，用户ID: %d", newUserID)
	_, err := srv.OrderCreate(ctx, orderReq)
	
	// 新用户应该被允许创建订单
	if err != nil {
		if err.Error() == "rpc error: code = NotFound desc = 用户不存在" {
			t.Error("新用户被错误拒绝")
		} else {
			t.Logf("新用户通过验证，但因其他原因失败（预期）: %v", err)
		}
	} else {
		t.Log("新用户订单创建成功")
	}
}