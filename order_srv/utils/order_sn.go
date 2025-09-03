package utils

import (
	"fmt"
	"math/rand"
	"time"
)

// GenerateOrderSn 生成唯一订单号
// 
// 订单号生成规则：
// - 总长度：26位
// - 组成部分：时间戳(14位) + 用户ID(8位) + 随机数(4位)
// - 示例：20240101123045000001230456
//   - 20240101123045: 2024年01月01日12时30分45秒
//   - 00000123: 用户ID 123，补零到8位
//   - 0456: 随机数 456，补零到4位
//
// 参数：
//   userID: 用户ID，用于标识订单所属用户
//
// 返回值：
//   string: 26位唯一订单号
func GenerateOrderSn(userID int32) string {
	// 获取当前时间戳，格式：YYYYMMDDHHMMSS (14位)
	// 使用Go标准时间格式：2006-01-02 15:04:05
	timestamp := time.Now().Format("20060102150405")
	
	// 将用户ID格式化为8位字符串，不足8位前补零
	// 例如：用户ID 123 -> "00000123"
	userIDStr := fmt.Sprintf("%08d", userID)
	
	// 生成0-9999的随机数，确保4位长度
	randomNum := rand.Intn(10000)
	randomStr := fmt.Sprintf("%04d", randomNum)
	
	// 拼接生成最终的26位订单号
	// 格式：时间戳(14) + 用户ID(8) + 随机数(4) = 26位
	return fmt.Sprintf("%s%s%s", timestamp, userIDStr, randomStr)
}