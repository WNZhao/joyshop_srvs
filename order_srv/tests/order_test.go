package tests

import (
	"context"
	"order_srv/handler"
	"order_srv/proto"
	"testing"
)

func setupOrderService() *handler.OrderServiceServer {
	return &handler.OrderServiceServer{}
}

func TestCartItemAddAndList(t *testing.T) {
	initTestEnv(t)
	srv := setupOrderService()
	ctx := context.Background()

	// 添加购物车
	addReq := &proto.CartItemRequest{
		UserId:     generateRandomUserID(),
		GoodsId:    generateRandomGoodsID(),
		GoodsName:  generateRandomGoodsName(),
		GoodsImage: generateRandomGoodsImage(),
		GoodsPrice: generateRandomGoodsPrice(),
		Nums:       generateRandomNums(),
		Checked:    true,
	}
	_, err := srv.CartItemAdd(ctx, addReq)
	if err != nil {
		t.Fatalf("添加购物车失败: %v", err)
	}

	// 查询购物车
	listReq := &proto.UserInfo{Id: addReq.UserId}
	_, err = srv.CartItemList(ctx, listReq)
	if err != nil {
		t.Fatalf("查询购物车失败: %v", err)
	}
}

func TestCartItemUpdateAndDelete(t *testing.T) {
	initTestEnv(t)
	srv := setupOrderService()
	ctx := context.Background()

	updateReq := &proto.CartItemRequest{
		Id:      1,
		UserId:  generateRandomUserID(),
		GoodsId: generateRandomGoodsID(),
		Nums:    generateRandomNums(),
		Checked: false,
	}
	_, err := srv.CartItemUpdate(ctx, updateReq)
	if err != nil {
		t.Fatalf("更新购物车失败: %v", err)
	}

	_, err = srv.CartItemDelete(ctx, updateReq)
	if err != nil {
		t.Fatalf("删除购物车失败: %v", err)
	}
}

func TestOrderCreateAndList(t *testing.T) {
	initTestEnv(t)
	srv := setupOrderService()
	ctx := context.Background()

	orderReq := &proto.OrderRequest{
		UserId:  generateRandomUserID(),
		Address: "测试地址",
		Name:    "张三",
		Mobile:  "13800000000",
		Post:    "请尽快发货",
	}
	_, err := srv.OrderCreate(ctx, orderReq)
	if err != nil {
		t.Fatalf("创建订单失败: %v", err)
	}

	listReq := &proto.OrderFilterRequest{UserId: orderReq.UserId, Page: 1, PageSize: 10}
	_, err = srv.OrderList(ctx, listReq)
	if err != nil {
		t.Fatalf("查询订单列表失败: %v", err)
	}
}

func TestOrderDetailUpdateDelete(t *testing.T) {
	initTestEnv(t)
	srv := setupOrderService()
	ctx := context.Background()

	orderReq := &proto.OrderRequest{Id: 1, UserId: generateRandomUserID()}
	_, err := srv.OrderDetail(ctx, orderReq)
	if err != nil {
		t.Fatalf("查询订单详情失败: %v", err)
	}

	updateReq := &proto.OrderStatus{Id: 1, OrderSn: generateRandomOrderSn(), Status: "TRADE_SUCCESS"}
	_, err = srv.OrderUpdate(ctx, updateReq)
	if err != nil {
		t.Fatalf("更新订单状态失败: %v", err)
	}

	delReq := &proto.OrderDelRequest{Id: 1, UserId: orderReq.UserId}
	_, err = srv.OrderDelete(ctx, delReq)
	if err != nil {
		t.Fatalf("删除订单失败: %v", err)
	}
}
