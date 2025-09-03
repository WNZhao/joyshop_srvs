package utils

import (
	"context"
	"fmt"
	"sort"
	"time"

	"go.uber.org/zap"
)

type BatchLockManager struct {
	locks   map[string]*RedisLock
	goodsIds []int32
}

// NewBatchLockManager 创建批量锁管理器
func NewBatchLockManager(goodsIds []int32, expiration time.Duration) *BatchLockManager {
	// 对商品ID排序，避免死锁
	sortedIds := make([]int32, len(goodsIds))
	copy(sortedIds, goodsIds)
	sort.Slice(sortedIds, func(i, j int) bool {
		return sortedIds[i] < sortedIds[j]
	})

	locks := make(map[string]*RedisLock)
	for _, goodsId := range sortedIds {
		lockKey := fmt.Sprintf("inventory:lock:%d", goodsId)
		locks[lockKey] = NewRedisLock(lockKey, expiration)
	}

	return &BatchLockManager{
		locks:   locks,
		goodsIds: sortedIds,
	}
}

// LockAll 获取所有商品的锁
func (blm *BatchLockManager) LockAll(ctx context.Context, retryCount int, retryInterval time.Duration) error {
	// 按照排序后的顺序获取锁，避免死锁
	acquiredLocks := make([]string, 0)
	
	defer func() {
		// 如果获取锁失败，释放已获取的锁
		if len(acquiredLocks) < len(blm.locks) {
			blm.unlockKeys(ctx, acquiredLocks)
		}
	}()

	for _, goodsId := range blm.goodsIds {
		lockKey := fmt.Sprintf("inventory:lock:%d", goodsId)
		lock := blm.locks[lockKey]
		
		locked, err := lock.TryLock(ctx, retryCount, retryInterval)
		if err != nil {
			return fmt.Errorf("获取商品%d的锁失败: %w", goodsId, err)
		}
		if !locked {
			return fmt.Errorf("无法获取商品%d的锁", goodsId)
		}
		
		acquiredLocks = append(acquiredLocks, lockKey)
		zap.S().Debugf("成功获取商品%d的锁", goodsId)
	}

	return nil
}

// UnlockAll 释放所有锁
func (blm *BatchLockManager) UnlockAll(ctx context.Context) {
	lockKeys := make([]string, 0, len(blm.locks))
	for key := range blm.locks {
		lockKeys = append(lockKeys, key)
	}
	blm.unlockKeys(ctx, lockKeys)
}

// unlockKeys 释放指定的锁
func (blm *BatchLockManager) unlockKeys(ctx context.Context, keys []string) {
	for _, key := range keys {
		if lock, exists := blm.locks[key]; exists {
			if err := lock.Unlock(ctx); err != nil {
				zap.S().Errorf("释放锁%s失败: %v", key, err)
			} else {
				zap.S().Debugf("成功释放锁%s", key)
			}
		}
	}
}