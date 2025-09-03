/*
 * @Author: Will nanan_zhao@163.com
 * @Date: 2025-05-18 14:06:45
 * @LastEditors: Please set LastEditors
 * @LastEditTime: 2025-06-07 21:56:56
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
	t.Log(TestScenarios.SetInventory)
	
	// 初始化测试环境
	initTestEnv(t)

	// 创建服务实例
	server := &handler.InventoryServer{}

	// 使用预定义的测试数据
	testInv := TestInventories[0] // iPhone库存
	goodsID, stock := testInv.GoodsID, testInv.Stock

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
		t.Logf("设置库存成功 - 商品ID: %d, 库存: %d", goodsID, stock)
	})
}

// TestGetInventory 测试获取库存
func TestGetInventory(t *testing.T) {
	t.Log(TestScenarios.GetInventory)
	
	// 初始化测试环境
	initTestEnv(t)

	// 创建服务实例
	server := &handler.InventoryServer{}

	// 使用预定义的测试数据
	testInv := TestInventories[1] // HUAWEI库存
	goodsID, stock := testInv.GoodsID, testInv.Stock

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
		t.Logf("获取库存成功 - 商品ID: %d (%s), 库存: %d", goodsID, testInv.GoodsName, resp.Num)
	})
}

// TestSell 测试库存扣减
func TestSell(t *testing.T) {
	t.Log(TestScenarios.SellInventory)
	
	// 初始化测试环境
	initTestEnv(t)

	// 创建服务实例
	server := &handler.InventoryServer{}

	// 使用预定义的测试数据
	testInv := TestInventories[2] // 小米库存（库存充足）
	goodsID, stock := testInv.GoodsID, testInv.Stock
	sellNum := int32(10) // 扩减10件

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
		t.Logf("扣减库存成功 - 商品ID: %d (%s), 原库存: %d, 扣减: %d, 余库存: %d", 
			goodsID, testInv.GoodsName, stock, sellNum, updatedInv.Stock)
	})
}

// TestReback 测试库存归还
func TestReback(t *testing.T) {
	t.Log(TestScenarios.RebackInventory)
	
	// 初始化测试环境
	initTestEnv(t)

	// 创建服务实例
	server := &handler.InventoryServer{}

	// 使用预定义的测试数据
	testInv := TestInventories[3] // Nike库存
	goodsID, stock := testInv.GoodsID, testInv.Stock
	rebackNum := int32(5) // 归还5件

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
		t.Logf("归还库存成功 - 商品ID: %d (%s), 原库存: %d, 归还: %d, 新库存: %d", 
			goodsID, testInv.GoodsName, stock, rebackNum, updatedInv.Stock)
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

// TestLowStockAlert 测试低库存预警
func TestLowStockAlert(t *testing.T) {
	t.Log(TestScenarios.LowStockAlert)
	
	// 初始化测试环境
	initTestEnv(t)
	
	// 创建服务实例
	server := &handler.InventoryServer{}
	
	// 获取所有低库存商品
	lowStockItems := GetLowStockInventories()
	t.Logf("预期低库存商品数量: %d", len(lowStockItems))
	
	for _, item := range lowStockItems {
		t.Run(fmt.Sprintf("低库存商品_%d", item.GoodsID), func(t *testing.T) {
			// 先设置库存
			inv := &model.Inventory{
				GoodsID: item.GoodsID,
				Stock:   item.Stock,
				Version: 0,
			}
			global.DB.Create(inv)
			
			// 查询库存
			resp, err := server.GetInventory(context.Background(), &proto.GoodsInvInfo{
				GoodsId: item.GoodsID,
			})
			if err != nil {
				t.Errorf("查询低库存商品失败: %v", err)
				return
			}
			
			// 验证是否为低库存
			if resp.Num > StockLevels.Threshold.Low {
				t.Errorf("商品不是低库存: 商品ID %d, 库存 %d, 阈值 %d", 
					item.GoodsID, resp.Num, StockLevels.Threshold.Low)
			}
			
			stockStatus := GetStockStatus(resp.Num)
			t.Logf("低库存商品 - ID: %d (%s), 库存: %d, 状态: %s", 
				item.GoodsID, item.GoodsName, resp.Num, stockStatus)
		})
	}
}

// TestInventoryValidation 测试库存数据验证
func TestInventoryValidation(t *testing.T) {
	t.Log("库存数据验证测试")
	
	// 初始化测试环境
	initTestEnv(t)
	
	// 创建服务实例
	server := &handler.InventoryServer{}
	
	// 验证所有预定义库存数据
	for _, testInv := range TestInventories {
		t.Run(fmt.Sprintf("验证库存_%d", testInv.GoodsID), func(t *testing.T) {
			// 设置库存
			inv := &model.Inventory{
				GoodsID: testInv.GoodsID,
				Stock:   testInv.Stock,
				Version: testInv.Version,
			}
			global.DB.Create(inv)
			
			// 查询库存
			resp, err := server.GetInventory(context.Background(), &proto.GoodsInvInfo{
				GoodsId: testInv.GoodsID,
			})
			if err != nil {
				t.Errorf("查询库存失败: %v", err)
				return
			}
			
			// 验证库存数据
			if resp.Num != testInv.Stock {
				t.Errorf("库存数量不匹配: 期望 %d, 实际 %d", testInv.Stock, resp.Num)
			}
			
			// 验证库存状态
			expectedStatus := GetStockStatus(testInv.Stock)
			if expectedStatus != testInv.Status {
				t.Errorf("库存状态不匹配: 期望 %s, 实际 %s", testInv.Status, expectedStatus)
			}
			
			t.Logf("库存验证通过 - ID: %d (%s), 库存: %d, 状态: %s", 
				testInv.GoodsID, testInv.GoodsName, resp.Num, expectedStatus)
		})
	}
}

// TestSellScenarios 测试各种扣减场景
func TestSellScenarios(t *testing.T) {
	t.Log("库存扣减场景测试")
	
	// 初始化测试环境
	initTestEnv(t)
	
	// 创建服务实例
	server := &handler.InventoryServer{}
	
	// 测试预定义的扣减场景
	for i, scenario := range TestSellScenarios {
		t.Run(fmt.Sprintf("扣减场景_%d", i+1), func(t *testing.T) {
			// 获取商品库存信息
			testInv := GetTestInventoryByGoodsID(scenario.GoodsID)
			if testInv == nil {
				t.Skipf("未找到商品ID %d 的库存信息", scenario.GoodsID)
				return
			}
			
			// 设置初始库存
			inv := &model.Inventory{
				GoodsID: testInv.GoodsID,
				Stock:   testInv.Stock,
				Version: 0,
			}
			global.DB.Create(inv)
			
			// 尝试扣减库存
			_, err := server.Sell(context.Background(), &proto.SellInfo{
				GoodsInvInfo: []*proto.GoodsInvInfo{
					{
						GoodsId: scenario.GoodsID,
						Num:     scenario.SellNum,
					},
				},
			})
			
			// 验证结果
			if scenario.ExpectedResult == "success" {
				if err != nil {
					t.Errorf("扣减应该成功但失败了: %v", err)
				} else {
					t.Logf("扣减成功 - 商品ID: %d (%s), 扣减: %d", 
						scenario.GoodsID, testInv.GoodsName, scenario.SellNum)
				}
			} else if scenario.ExpectedResult == "insufficient" {
				if err == nil {
					t.Errorf("库存不足时扣减应该失败但成功了")
				} else {
					t.Logf("库存不足扣减正确失败 - 商品ID: %d (%s), 尝试扣减: %d, 可用库存: %d", 
						scenario.GoodsID, testInv.GoodsName, scenario.SellNum, testInv.Stock)
				}
			}
		})
	}
}
