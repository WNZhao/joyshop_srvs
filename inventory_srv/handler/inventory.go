package handler

import (
	"context"
	"fmt"
	"sync"
	"time"

	"inventory_srv/global"
	"inventory_srv/model"
	"inventory_srv/proto"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type InventoryServer struct {
	proto.UnimplementedInventoryServiceServer
	mu sync.Mutex
}

// SetInventory 设置库存
func (s *InventoryServer) SetInventory(ctx context.Context, req *proto.GoodsInvInfo) (*emptypb.Empty, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var inv model.Inventory
	result := global.DB.Where("goods_id = ?", req.GoodsId).First(&inv)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// 如果记录不存在，创建新记录
			inv = model.Inventory{
				GoodsID: req.GoodsId,
				Stock:   req.Num,
				Version: 0,
			}
			if err := global.DB.Create(&inv).Error; err != nil {
				zap.S().Errorf("创建库存记录失败: %v", err)
				return nil, status.Error(codes.Internal, "创建库存记录失败")
			}
		} else {
			zap.S().Errorf("查询库存记录失败: %v", result.Error)
			return nil, status.Error(codes.Internal, "查询库存记录失败")
		}
	} else {
		// 更新现有记录
		inv.Stock = req.Num
		if err := global.DB.Save(&inv).Error; err != nil {
			zap.S().Errorf("更新库存记录失败: %v", err)
			return nil, status.Error(codes.Internal, "更新库存记录失败")
		}
	}

	return &emptypb.Empty{}, nil
}

// GetInventory 获取库存
func (s *InventoryServer) GetInventory(ctx context.Context, req *proto.GoodsInvInfo) (*proto.GoodsInvInfo, error) {
	var inv model.Inventory
	result := global.DB.Where("goods_id = ?", req.GoodsId).First(&inv)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return &proto.GoodsInvInfo{
				GoodsId: req.GoodsId,
				Num:     0,
			}, nil
		}
		zap.S().Errorf("查询库存记录失败: %v", result.Error)
		return nil, status.Error(codes.Internal, "查询库存记录失败")
	}

	return &proto.GoodsInvInfo{
		GoodsId: inv.GoodsID,
		Num:     inv.Stock,
	}, nil
}

// 获取分布式锁
func (s *InventoryServer) acquireLock(ctx context.Context, key string, timeout time.Duration) (bool, error) {
	return global.RedisClient.SetNX(ctx, key, "1", timeout).Result()
}

// 释放分布式锁
func (s *InventoryServer) releaseLock(ctx context.Context, key string) error {
	return global.RedisClient.Del(ctx, key).Err()
}

// Sell 库存扣减 - Redis分布式锁实现
func (s *InventoryServer) Sell(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	const lockTimeout = 10 * time.Second
	const maxRetries = 3
	const retryDelay = 100 * time.Millisecond

	for retry := 0; retry < maxRetries; retry++ {
		// 获取分布式锁
		lockKey := fmt.Sprintf("inventory:lock:%d", req.GoodsInvInfo[0].GoodsId)
		acquired, err := s.acquireLock(ctx, lockKey, lockTimeout)
		if err != nil {
			zap.S().Errorf("获取分布式锁失败: %v", err)
			return nil, status.Error(codes.Internal, "获取分布式锁失败")
		}
		if !acquired {
			if retry < maxRetries-1 {
				time.Sleep(retryDelay)
				continue
			}
			return nil, status.Error(codes.ResourceExhausted, "获取锁失败，请重试")
		}

		// 确保在函数返回时释放锁
		defer s.releaseLock(ctx, lockKey)

		// 开启事务
		tx := global.DB.Begin()
		if tx.Error != nil {
			return nil, status.Error(codes.Internal, "开启事务失败")
		}

		success := true
		// 遍历所有商品
		for _, goodsInfo := range req.GoodsInvInfo {
			var inv model.Inventory
			result := tx.Where("goods_id = ?", goodsInfo.GoodsId).First(&inv)
			if result.Error != nil {
				tx.Rollback()
				if result.Error == gorm.ErrRecordNotFound {
					return nil, status.Error(codes.NotFound, "商品库存不存在")
				}
				return nil, status.Error(codes.Internal, "查询库存失败")
			}

			// 检查库存是否充足
			if inv.Stock < goodsInfo.Num {
				tx.Rollback()
				return nil, status.Error(codes.ResourceExhausted, "库存不足")
			}

			// 更新库存
			inv.Stock -= goodsInfo.Num
			inv.Version += 1 // 增加版本号
			if err := tx.Save(&inv).Error; err != nil {
				tx.Rollback()
				success = false
				break
			}
		}

		if success {
			if err := tx.Commit().Error; err != nil {
				return nil, status.Error(codes.Internal, "提交事务失败")
			}
			return &emptypb.Empty{}, nil
		}

		if retry < maxRetries-1 {
			time.Sleep(retryDelay)
			continue
		}
	}

	return nil, status.Error(codes.Internal, "库存更新失败，请重试")
}

// Reback 库存归还 - Redis分布式锁实现
func (s *InventoryServer) Reback(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	const lockTimeout = 10 * time.Second
	const maxRetries = 3
	const retryDelay = 100 * time.Millisecond

	for retry := 0; retry < maxRetries; retry++ {
		// 获取分布式锁
		lockKey := fmt.Sprintf("inventory:lock:%d", req.GoodsInvInfo[0].GoodsId)
		acquired, err := s.acquireLock(ctx, lockKey, lockTimeout)
		if err != nil {
			zap.S().Errorf("获取分布式锁失败: %v", err)
			return nil, status.Error(codes.Internal, "获取分布式锁失败")
		}
		if !acquired {
			if retry < maxRetries-1 {
				time.Sleep(retryDelay)
				continue
			}
			return nil, status.Error(codes.ResourceExhausted, "获取锁失败，请重试")
		}

		// 确保在函数返回时释放锁
		defer s.releaseLock(ctx, lockKey)

		// 开启事务
		tx := global.DB.Begin()
		if tx.Error != nil {
			return nil, status.Error(codes.Internal, "开启事务失败")
		}

		success := true
		// 遍历所有商品
		for _, goodsInfo := range req.GoodsInvInfo {
			var inv model.Inventory
			result := tx.Where("goods_id = ?", goodsInfo.GoodsId).First(&inv)
			if result.Error != nil {
				tx.Rollback()
				if result.Error == gorm.ErrRecordNotFound {
					return nil, status.Error(codes.NotFound, "商品库存不存在")
				}
				return nil, status.Error(codes.Internal, "查询库存失败")
			}

			// 更新库存
			inv.Stock += goodsInfo.Num
			inv.Version += 1 // 增加版本号
			if err := tx.Save(&inv).Error; err != nil {
				tx.Rollback()
				success = false
				break
			}
		}

		if success {
			if err := tx.Commit().Error; err != nil {
				return nil, status.Error(codes.Internal, "提交事务失败")
			}
			return &emptypb.Empty{}, nil
		}

		if retry < maxRetries-1 {
			time.Sleep(retryDelay)
			continue
		}
	}

	return nil, status.Error(codes.Internal, "库存更新失败，请重试")
}
