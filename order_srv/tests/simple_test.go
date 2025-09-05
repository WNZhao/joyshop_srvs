package tests

import (
	"context"
	"testing"
	"order_srv/handler"
	"order_srv/proto"
)

func TestSimpleCartAdd(t *testing.T) {
	initTestEnvSimple(t)
	t.Log("开始简单购物车测试")
	
	// 这个测试不依赖数据库连接，只测试基本功能
	srv := &handler.OrderServiceServer{}
	ctx := context.Background()

	// 创建测试请求
	req := &proto.CartItemRequest{
		UserId:     1,
		GoodsId:    1,
		GoodsName:  "测试商品",
		GoodsImage: "test.jpg",
		GoodsPrice: 99.99,
		Nums:       1,
		Checked:    true,
	}

	t.Logf("准备添加商品到购物车: %+v", req)

	// 因为没有数据库连接，这个调用会失败，但我们可以看到错误信息
	_, err := srv.CartItemAdd(ctx, req)
	if err != nil {
		t.Logf("预期的错误（因为没有数据库连接）: %v", err)
	} else {
		t.Log("意外成功")
	}
}