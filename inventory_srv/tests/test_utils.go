/*
 * @Description:
 */
/*
 * @Author: Will nanan_zhao@163.com
 * @Date: 2025-05-18 14:06:45
 * @LastEditors: Please set LastEditors
 * @LastEditTime: 2025-06-05 17:06:46
 * @FilePath: /joyshop_srvs/inventory_srv/tests/test_utils.go
 * @Description: 库存服务测试工具方法
 */
package tests

import (
	"math/rand"
	"time"
)

// 创建一个本地随机数生成器
var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

// 生成随机库存数量
func generateRandomStock() int32 {
	// 生成 0-1000 之间的随机数
	return int32(rnd.Intn(1001))
}

// 生成随机商品ID
func generateRandomGoodsID() int32 {
	// 生成 1-1000 之间的随机数
	return int32(rnd.Intn(1000) + 1)
}

// 生成随机库存信息
func generateRandomInventory() (int32, int32) {
	return generateRandomGoodsID(), generateRandomStock()
}

// 生成测试用的库存扣减信息
func generateRandomSellInfo(count int) []map[string]int32 {
	result := make([]map[string]int32, count)
	for i := 0; i < count; i++ {
		result[i] = map[string]int32{
			"goods_id": generateRandomGoodsID(),
			"num":      int32(rnd.Intn(10) + 1), // 生成 1-10 之间的随机数
		}
	}
	return result
}
