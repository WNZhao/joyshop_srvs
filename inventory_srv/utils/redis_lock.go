package utils

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"inventory_srv/global"
)

type RedisLock struct {
	key        string
	value      string
	expiration time.Duration
}

// NewRedisLock 创建新的Redis分布式锁
func NewRedisLock(key string, expiration time.Duration) *RedisLock {
	return &RedisLock{
		key:        key,
		value:      generateLockValue(),
		expiration: expiration,
	}
}

// Lock 获取锁
func (l *RedisLock) Lock(ctx context.Context) (bool, error) {
	result, err := global.RedisClient.SetNX(ctx, l.key, l.value, l.expiration).Result()
	if err != nil {
		return false, fmt.Errorf("获取分布式锁失败: %w", err)
	}
	return result, nil
}

// TryLock 尝试获取锁，带重试机制
func (l *RedisLock) TryLock(ctx context.Context, retryCount int, retryInterval time.Duration) (bool, error) {
	for i := 0; i < retryCount; i++ {
		locked, err := l.Lock(ctx)
		if err != nil {
			return false, err
		}
		if locked {
			return true, nil
		}
		
		// 如果没有获取到锁，等待一段时间后重试
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		case <-time.After(retryInterval):
			continue
		}
	}
	return false, fmt.Errorf("尝试 %d 次后仍无法获取锁", retryCount)
}

// Unlock 释放锁
func (l *RedisLock) Unlock(ctx context.Context) error {
	// 使用Lua脚本确保只有持有锁的客户端才能释放锁
	luaScript := `
		if redis.call("GET", KEYS[1]) == ARGV[1] then
			return redis.call("DEL", KEYS[1])
		else
			return 0
		end
	`
	
	result, err := global.RedisClient.Eval(ctx, luaScript, []string{l.key}, l.value).Result()
	if err != nil {
		return fmt.Errorf("释放分布式锁失败: %w", err)
	}
	
	if result.(int64) == 0 {
		return fmt.Errorf("无法释放锁，可能已被其他客户端释放")
	}
	
	return nil
}

// generateLockValue 生成锁的唯一值
func generateLockValue() string {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		// 如果随机数生成失败，使用时间戳作为备选方案
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(bytes)
}