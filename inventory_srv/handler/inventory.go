package handler

import (
	"context"
	"fmt"
	"sync"
	"time"

	"inventory_srv/global"
	"inventory_srv/model"
	"inventory_srv/proto"
	"inventory_srv/utils"

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


// Sell 库存扣减 - 改进的Redis分布式锁实现
func (s *InventoryServer) Sell(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	if len(req.GoodsInvInfo) == 0 {
		return nil, status.Error(codes.InvalidArgument, "商品信息不能为空")
	}

	zap.S().Infof("开始扣减库存，商品数量: %d", len(req.GoodsInvInfo))

	// 提取所有商品ID
	goodsIds := make([]int32, len(req.GoodsInvInfo))
	for i, goodsInfo := range req.GoodsInvInfo {
		goodsIds[i] = goodsInfo.GoodsId
		zap.S().Infof("扣减商品: ID=%d, 数量=%d", goodsInfo.GoodsId, goodsInfo.Num)
	}

	// 创建批量锁管理器
	lockManager := utils.NewBatchLockManager(goodsIds, 15*time.Second)
	
	// 获取所有商品的锁
	if err := lockManager.LockAll(ctx, 3, 100*time.Millisecond); err != nil {
		zap.S().Errorf("获取批量锁失败: %v", err)
		return nil, status.Error(codes.Internal, "系统忙，请稍后重试")
	}

	// 确保释放所有锁
	defer lockManager.UnlockAll(ctx)

	// 开启事务
	tx := global.DB.Begin()
	if tx.Error != nil {
		zap.S().Errorf("开启事务失败: %v", tx.Error)
		return nil, status.Error(codes.Internal, "开启事务失败")
	}
	
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			zap.S().Errorf("库存扣减过程中发生panic: %v", r)
		}
	}()

	// 第一阶段：检查所有商品库存是否充足
	inventories := make(map[int32]*model.Inventory)
	for _, goodsInfo := range req.GoodsInvInfo {
		var inv model.Inventory
		result := tx.Where("goods_id = ?", goodsInfo.GoodsId).First(&inv)
		if result.Error != nil {
			tx.Rollback()
			if result.Error == gorm.ErrRecordNotFound {
				zap.S().Errorf("商品库存不存在，商品ID: %d", goodsInfo.GoodsId)
				return nil, status.Error(codes.NotFound, fmt.Sprintf("商品%d库存不存在", goodsInfo.GoodsId))
			}
			zap.S().Errorf("查询库存失败: %v", result.Error)
			return nil, status.Error(codes.Internal, "查询库存失败")
		}

		// 检查库存是否充足
		if inv.Stock < goodsInfo.Num {
			tx.Rollback()
			zap.S().Warnf("库存不足，商品ID: %d，当前库存: %d，需要: %d", 
				goodsInfo.GoodsId, inv.Stock, goodsInfo.Num)
			return nil, status.Error(codes.ResourceExhausted, 
				fmt.Sprintf("商品%d库存不足，当前库存%d，需要%d", goodsInfo.GoodsId, inv.Stock, goodsInfo.Num))
		}

		inventories[goodsInfo.GoodsId] = &inv
		zap.S().Infof("商品%d库存检查通过，当前库存: %d，扣减: %d", 
			goodsInfo.GoodsId, inv.Stock, goodsInfo.Num)
	}

	// 第二阶段：执行库存扣减
	for _, goodsInfo := range req.GoodsInvInfo {
		inv := inventories[goodsInfo.GoodsId]
		oldStock := inv.Stock
		
		// 使用乐观锁更新库存
		result := tx.Model(inv).
			Where("goods_id = ? AND version = ?", inv.GoodsID, inv.Version).
			Updates(map[string]interface{}{
				"stock":   inv.Stock - goodsInfo.Num,
				"version": inv.Version + 1,
			})

		if result.Error != nil {
			tx.Rollback()
			zap.S().Errorf("更新库存失败: %v", result.Error)
			return nil, status.Error(codes.Internal, "更新库存失败")
		}

		if result.RowsAffected == 0 {
			tx.Rollback()
			zap.S().Warnf("乐观锁冲突，商品ID: %d", goodsInfo.GoodsId)
			return nil, status.Error(codes.Aborted, "库存并发冲突，请重试")
		}

		zap.S().Infof("成功扣减库存，商品ID: %d，原库存: %d，扣减: %d，新库存: %d", 
			goodsInfo.GoodsId, oldStock, goodsInfo.Num, oldStock-goodsInfo.Num)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		zap.S().Errorf("提交事务失败: %v", err)
		return nil, status.Error(codes.Internal, "提交事务失败")
	}

	zap.S().Infof("库存扣减成功完成，共处理%d个商品", len(req.GoodsInvInfo))
	return &emptypb.Empty{}, nil
}

// Reback 库存归还 - 改进的Redis分布式锁实现
func (s *InventoryServer) Reback(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	if len(req.GoodsInvInfo) == 0 {
		return nil, status.Error(codes.InvalidArgument, "商品信息不能为空")
	}

	zap.S().Infof("开始归还库存，商品数量: %d", len(req.GoodsInvInfo))

	// 提取所有商品ID
	goodsIds := make([]int32, len(req.GoodsInvInfo))
	for i, goodsInfo := range req.GoodsInvInfo {
		goodsIds[i] = goodsInfo.GoodsId
		zap.S().Infof("归还商品: ID=%d, 数量=%d", goodsInfo.GoodsId, goodsInfo.Num)
	}

	// 创建批量锁管理器
	lockManager := utils.NewBatchLockManager(goodsIds, 15*time.Second)
	
	// 获取所有商品的锁
	if err := lockManager.LockAll(ctx, 3, 100*time.Millisecond); err != nil {
		zap.S().Errorf("获取批量锁失败: %v", err)
		return nil, status.Error(codes.Internal, "系统忙，请稍后重试")
	}

	// 确保释放所有锁
	defer lockManager.UnlockAll(ctx)

	// 开启事务
	tx := global.DB.Begin()
	if tx.Error != nil {
		zap.S().Errorf("开启事务失败: %v", tx.Error)
		return nil, status.Error(codes.Internal, "开启事务失败")
	}
	
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			zap.S().Errorf("库存归还过程中发生panic: %v", r)
		}
	}()

	// 查询并归还库存
	for _, goodsInfo := range req.GoodsInvInfo {
		var inv model.Inventory
		result := tx.Where("goods_id = ?", goodsInfo.GoodsId).First(&inv)
		if result.Error != nil {
			tx.Rollback()
			if result.Error == gorm.ErrRecordNotFound {
				zap.S().Errorf("商品库存不存在，商品ID: %d", goodsInfo.GoodsId)
				return nil, status.Error(codes.NotFound, fmt.Sprintf("商品%d库存不存在", goodsInfo.GoodsId))
			}
			zap.S().Errorf("查询库存失败: %v", result.Error)
			return nil, status.Error(codes.Internal, "查询库存失败")
		}

		oldStock := inv.Stock
		
		// 使用乐观锁更新库存
		result = tx.Model(&inv).
			Where("goods_id = ? AND version = ?", inv.GoodsID, inv.Version).
			Updates(map[string]interface{}{
				"stock":   inv.Stock + goodsInfo.Num,
				"version": inv.Version + 1,
			})

		if result.Error != nil {
			tx.Rollback()
			zap.S().Errorf("归还库存失败: %v", result.Error)
			return nil, status.Error(codes.Internal, "归还库存失败")
		}

		if result.RowsAffected == 0 {
			tx.Rollback()
			zap.S().Warnf("乐观锁冲突，商品ID: %d", goodsInfo.GoodsId)
			return nil, status.Error(codes.Aborted, "库存并发冲突，请重试")
		}

		zap.S().Infof("成功归还库存，商品ID: %d，原库存: %d，归还: %d，新库存: %d", 
			goodsInfo.GoodsId, oldStock, goodsInfo.Num, oldStock+goodsInfo.Num)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		zap.S().Errorf("提交事务失败: %v", err)
		return nil, status.Error(codes.Internal, "提交事务失败")
	}

	zap.S().Infof("库存归还成功完成，共处理%d个商品", len(req.GoodsInvInfo))
	return &emptypb.Empty{}, nil
}
