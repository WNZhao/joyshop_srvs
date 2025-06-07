/*
 * @Author: Will nanan_zhao@163.com
 * @Date: 2025-05-18 14:06:45
 * @LastEditors: Please set LastEditors
 * @LastEditTime: 2025-06-07 21:49:34
 * @FilePath: /joyshop_srvs/inventory_srv/tests/inventory_test.go
 * @Description: 库存服务测试用例
 */
package tests

import (
	"context"
	"testing"

	"inventory_srv/global"
	"inventory_srv/handler"
	"inventory_srv/model"
	"inventory_srv/proto"

	"github.com/stretchr/testify/assert"
)

// TestSetInventory 测试设置库存
func TestSetInventory(t *testing.T) {
	// 初始化测试环境
	initTestEnv(t)

	// 创建服务实例
	server := &handler.InventoryServer{}

	// 生成测试数据
	goodsID, stock := generateRandomInventory()

	// 测试设置库存
	t.Run("SetInventory", func(t *testing.T) {
		req := &proto.GoodsInvInfo{
			GoodsId: goodsID,
			Num:     stock,
		}

		// 调用服务
		_, err := server.SetInventory(context.Background(), req)
		assert.NoError(t, err)

		// 验证数据库中的记录
		var inv model.Inventory
		result := global.DB.Where("goods_id = ?", goodsID).First(&inv)
		assert.NoError(t, result.Error)
		assert.Equal(t, goodsID, inv.GoodsID)
		assert.Equal(t, stock, inv.Stock)
	})
}

// TestGetInventory 测试获取库存
func TestGetInventory(t *testing.T) {
	// 初始化测试环境
	initTestEnv(t)

	// 创建服务实例
	server := &handler.InventoryServer{}

	// 生成测试数据
	goodsID, stock := generateRandomInventory()

	// 先设置库存
	inv := &model.Inventory{
		GoodsID: goodsID,
		Stock:   stock,
		Version: 0,
	}
	global.DB.Create(inv)

	// 测试获取库存
	t.Run("GetInventory", func(t *testing.T) {
		req := &proto.GoodsInvInfo{
			GoodsId: goodsID,
		}

		// 调用服务
		resp, err := server.GetInventory(context.Background(), req)
		assert.NoError(t, err)
		assert.Equal(t, goodsID, resp.GoodsId)
		assert.Equal(t, stock, resp.Num)
	})
}

// TestSell 测试库存扣减
func TestSell(t *testing.T) {
	// 初始化测试环境
	initTestEnv(t)

	// 创建服务实例
	server := &handler.InventoryServer{}

	// 生成测试数据
	goodsID, stock := generateRandomInventory()
	sellNum := int32(5)

	// 先设置库存
	inv := &model.Inventory{
		GoodsID: goodsID,
		Stock:   stock,
		Version: 0,
	}
	global.DB.Create(inv)

	// 测试库存扣减
	t.Run("Sell", func(t *testing.T) {
		req := &proto.SellInfo{
			GoodsInvInfo: []*proto.GoodsInvInfo{
				{
					GoodsId: goodsID,
					Num:     sellNum,
				},
			},
		}

		// 调用服务
		_, err := server.Sell(context.Background(), req)
		assert.NoError(t, err)

		// 验证数据库中的记录
		var updatedInv model.Inventory
		result := global.DB.Where("goods_id = ?", goodsID).First(&updatedInv)
		assert.NoError(t, result.Error)
		assert.Equal(t, stock-sellNum, updatedInv.Stock)
		assert.Equal(t, int32(1), updatedInv.Version)
	})
}

// TestReback 测试库存归还
func TestReback(t *testing.T) {
	// 初始化测试环境
	initTestEnv(t)

	// 创建服务实例
	server := &handler.InventoryServer{}

	// 生成测试数据
	goodsID, stock := generateRandomInventory()
	rebackNum := int32(5)

	// 先设置库存
	inv := &model.Inventory{
		GoodsID: goodsID,
		Stock:   stock,
		Version: 0,
	}
	global.DB.Create(inv)

	// 测试库存归还
	t.Run("Reback", func(t *testing.T) {
		req := &proto.SellInfo{
			GoodsInvInfo: []*proto.GoodsInvInfo{
				{
					GoodsId: goodsID,
					Num:     rebackNum,
				},
			},
		}

		// 调用服务
		_, err := server.Reback(context.Background(), req)
		assert.NoError(t, err)

		// 验证数据库中的记录
		var updatedInv model.Inventory
		result := global.DB.Where("goods_id = ?", goodsID).First(&updatedInv)
		assert.NoError(t, result.Error)
		assert.Equal(t, stock+rebackNum, updatedInv.Stock)
		assert.Equal(t, int32(1), updatedInv.Version)
	})
}

// TestConcurrentSellAndReback 测试并发扣减和归还
func TestConcurrentSellAndReback(t *testing.T) {
	// 初始化测试环境
	initTestEnv(t)

	// 创建服务实例
	server := &handler.InventoryServer{}

	// 生成测试数据
	goodsID, stock := generateRandomInventory()
	sellNum := int32(5)
	rebackNum := int32(3)

	// 先设置库存
	inv := &model.Inventory{
		GoodsID: goodsID,
		Stock:   stock,
		Version: 0,
	}
	global.DB.Create(inv)

	// 测试并发扣减和归还
	t.Run("ConcurrentSellAndReback", func(t *testing.T) {
		// 创建扣减请求
		sellReq := &proto.SellInfo{
			GoodsInvInfo: []*proto.GoodsInvInfo{
				{
					GoodsId: goodsID,
					Num:     sellNum,
				},
			},
		}

		// 创建归还请求
		rebackReq := &proto.SellInfo{
			GoodsInvInfo: []*proto.GoodsInvInfo{
				{
					GoodsId: goodsID,
					Num:     rebackNum,
				},
			},
		}

		// 模拟并发请求
		done := make(chan error, 2)
		go func() {
			_, err := server.Sell(context.Background(), sellReq)
			done <- err
		}()
		go func() {
			_, err := server.Reback(context.Background(), rebackReq)
			done <- err
		}()

		// 等待所有请求完成
		for i := 0; i < 2; i++ {
			err := <-done
			if err != nil {
				t.Logf("并发请求 %d 失败: %v", i+1, err)
			}
		}

		// 验证数据库中的记录
		var updatedInv model.Inventory
		result := global.DB.Where("goods_id = ?", goodsID).First(&updatedInv)
		assert.NoError(t, result.Error)
		assert.Equal(t, stock-sellNum+rebackNum, updatedInv.Stock)
		assert.Equal(t, int32(2), updatedInv.Version)
	})
}
