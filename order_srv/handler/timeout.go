package handler

import (
	"context"
	"fmt"
	"order_srv/global"
	"order_srv/model"
	"order_srv/utils"
	"time"

	inventorypb "order_srv/proto/inventory"
)

// CheckOrderTimeout 检查并关闭超时订单
func CheckOrderTimeout() {
	global.Logger.Info("开始检查超时订单...")
	
	// 查询所有超时未支付的订单
	var orders []model.OrderInfo
	now := time.Now()
	
	// 查询条件：
	// 1. 状态为 PAYING（待支付）
	// 2. 支付截止时间已过
	result := global.DB.Where("status = ? AND pay_deadline < ? AND pay_deadline IS NOT NULL", 
		"PAYING", now).Find(&orders)
	
	if result.Error != nil {
		global.Logger.Errorf("查询超时订单失败: %v", result.Error)
		return
	}
	
	if len(orders) == 0 {
		global.Logger.Debug("没有超时订单需要处理")
		return
	}
	
	global.Logger.Infof("发现 %d 个超时订单需要关闭", len(orders))
	
	// 批量关闭超时订单
	for _, order := range orders {
		if err := CloseTimeoutOrder(&order); err != nil {
			global.Logger.Errorf("关闭超时订单失败，订单ID: %d，错误: %v", order.ID, err)
			continue
		}
		global.Logger.Infof("成功关闭超时订单，订单号: %s", order.OrderSn)
	}
	
	global.Logger.Infof("超时订单检查完成，共处理 %d 个订单", len(orders))
}

// CloseTimeoutOrder 关闭单个超时订单
func CloseTimeoutOrder(order *model.OrderInfo) error {
	// 开启事务
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			global.Logger.Errorf("关闭超时订单异常: %v", r)
		}
	}()
	
	// 更新订单状态为超时关闭
	if err := tx.Model(order).Update("status", "TRADE_CLOSED").Error; err != nil {
		tx.Rollback()
		return err
	}
	
	// 恢复库存（如果已经扣减的话）
	// 这里需要调用库存服务恢复库存
	// 由于订单商品信息存在OrderGoods表中，需要查询并恢复
	var orderGoods []model.OrderGoods
	if err := tx.Where("order = ?", order.ID).Find(&orderGoods).Error; err != nil {
		tx.Rollback()
		return err
	}
	
	// 调用已有的库存回滚函数
	if err := RollbackOrderInventory(order.ID); err != nil {
		tx.Rollback()
		return fmt.Errorf("归还库存失败: %w", err)
	}
	
	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return err
	}
	
	global.Logger.Infof("订单 %s 已超时关闭", order.OrderSn)
	return nil
}

// StartTimeoutChecker 启动定时检查超时订单
func StartTimeoutChecker() {
	global.Logger.Info("启动订单超时检查定时任务...")
	
	// 每分钟检查一次超时订单
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		// 首次启动立即执行一次
		CheckOrderTimeout()
		
		for range ticker.C {
			CheckOrderTimeout()
		}
	}()
	
	global.Logger.Info("订单超时检查定时任务已启动")
}

// CheckOrderTimeoutWithRedis 使用Redis优化的超时检查（可选实现）
func CheckOrderTimeoutWithRedis() {
	// 使用Redis的过期键通知机制
	// 当订单创建时，设置一个Redis键，过期时间为支付截止时间
	// 监听Redis过期事件，触发订单关闭
	// 这种方式更高效，避免了定时轮询数据库
	
	// 示例代码：
	// 1. 订单创建时：redis.SetEX("order_timeout:{order_id}", 30*60, "1")
	// 2. 监听过期事件：redis.Subscribe("__keyevent@0__:expired")
	// 3. 处理过期事件：解析order_id，关闭对应订单
}

// RollbackOrderInventory 回滚订单库存
func RollbackOrderInventory(orderId int32) error {
	// 查询订单商品
	var orderGoods []model.OrderGoods
	if err := global.DB.Where("order = ?", orderId).Find(&orderGoods).Error; err != nil {
		return err
	}
	
	// 将订单商品信息转换为库存归还格式
	rebackItems := make([]*inventorypb.GoodsInvInfo, 0, len(orderGoods))
	for _, orderGood := range orderGoods {
		rebackItems = append(rebackItems, &inventorypb.GoodsInvInfo{
			GoodsId: orderGood.Goods,
			Num:     orderGood.Nums,
		})
	}
	
	// 批量归还库存
	if len(rebackItems) > 0 {
		return utils.RebackInventory(context.Background(), rebackItems)
	}
	
	return nil
}