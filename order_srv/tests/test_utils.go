/*
 * @Description:
 */
/*
 * @Author: Will nanan_zhao@163.com
 * @Date: 2025-05-18 14:06:45
 * @LastEditors: Please set LastEditors
 * @LastEditTime: 2025-06-08 17:25:36
 * @FilePath: /joyshop_srvs/order_srv/tests/test_utils.go
 * @Description: 库存服务测试工具方法
 */
package tests

import (
	"fmt"
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

// 生成随机用户ID
func generateRandomUserID() int32 {
	return int32(rnd.Intn(1000) + 1)
}

// 生成随机订单号
func generateRandomOrderSn() string {
	return fmt.Sprintf("ORDER%06d", rnd.Intn(1000000))
}

// 生成随机商品名称
func generateRandomGoodsName() string {
	return fmt.Sprintf("商品%d", rnd.Intn(10000))
}

// 生成随机图片链接
func generateRandomGoodsImage() string {
	return fmt.Sprintf("http://example.com/image/%d.jpg", rnd.Intn(10000))
}

// 生成随机价格
func generateRandomGoodsPrice() float32 {
	return float32(rnd.Intn(10000)) / 100.0
}

// 生成随机数量
func generateRandomNums() int32 {
	return int32(rnd.Intn(10) + 1)
}
