# 分布式锁机制说明文档

## 概述

本文档详细说明了订单服务中分布式锁的设计与实现，包括Redis分布式锁的使用场景、实现原理、配置参数和最佳实践。

## 使用场景

### 1. 订单创建锁
- **锁键**: `order_create_lock:{userId}`
- **用途**: 防止用户重复下单
- **锁定时长**: 30秒
- **应用场景**: 用户点击下单按钮时防止重复提交

### 2. 订单更新锁
- **锁键**: `order_update_lock:id:{orderId}` 或 `order_update_lock:sn:{orderSn}`
- **用途**: 防止订单状态并发更新
- **锁定时长**: 10秒
- **应用场景**: 支付回调、订单状态变更时

### 3. 订单删除锁
- **锁键**: `order_delete_lock:{orderId}`
- **用途**: 防止订单并发删除
- **锁定时长**: 15秒
- **应用场景**: 取消订单操作时

## 实现原理

### Redis分布式锁实现

```go
type RedisLock struct {
    key        string
    value      string
    expiration time.Duration
    client     *redis.Client
}

func (r *RedisLock) TryLock(ctx context.Context, maxRetries int, retryInterval time.Duration) (bool, error) {
    for i := 0; i < maxRetries; i++ {
        result, err := r.client.SetNX(ctx, r.key, r.value, r.expiration).Result()
        if err != nil {
            return false, fmt.Errorf("获取锁失败: %w", err)
        }
        
        if result {
            return true, nil
        }
        
        if i < maxRetries-1 {
            time.Sleep(retryInterval)
        }
    }
    
    return false, nil
}
```

### 锁的唯一标识

```go
func generateLockValue() string {
    hostname, _ := os.Hostname()
    pid := os.Getpid()
    timestamp := time.Now().UnixNano()
    return fmt.Sprintf("%s-%d-%d", hostname, pid, timestamp)
}
```

## 配置参数详解

### 订单创建锁配置

| 参数 | 值 | 说明 | 原因 |
|------|----|----- |------|
| 锁定时长 | 30秒 | 订单创建完整流程时间 | 包含跨服务调用和批量插入的最大预期时间 |
| 重试次数 | 3次 | 获取锁的最大重试次数 | 平衡用户体验和系统负载 |
| 重试间隔 | 100ms | 两次重试之间的等待时间 | 避免频繁重试导致的系统压力 |

### 订单更新锁配置

| 参数 | 值 | 说明 | 原因 |
|------|----|----- |------|
| 锁定时长 | 10秒 | 订单状态更新流程时间 | 更新操作相对简单，时间较短 |
| 重试次数 | 3次 | 获取锁的最大重试次数 | 支付回调等场景需要一定容错性 |
| 重试间隔 | 50ms | 两次重试之间的等待时间 | 更新操作时效性要求更高 |

### 订单删除锁配置

| 参数 | 值 | 说明 | 原因 |
|------|----|----- |------|
| 锁定时长 | 15秒 | 订单删除完整流程时间 | 包含订单和订单商品的批量删除 |
| 重试次数 | 3次 | 获取锁的最大重试次数 | 删除操作相对不那么紧急 |
| 重试间隔 | 50ms | 两次重试之间的等待时间 | 平衡响应时间和重试效果 |

## 代码实现示例

### 订单创建锁使用

```go
func (s *OrderServiceServer) OrderCreate(ctx context.Context, req *proto.OrderRequest) (*proto.OrderInfoResponse, error) {
    // 创建分布式锁，防止用户重复下单
    lockKey := fmt.Sprintf("order_create_lock:%d", req.UserId)
    lock := utils.NewRedisLock(lockKey, 30*time.Second)
    
    locked, err := lock.TryLock(ctx, 3, 100*time.Millisecond)
    if err != nil {
        global.Logger.Errorf("获取创建订单锁失败: %v", err)
        return nil, status.Errorf(codes.Internal, "系统忙，请稍后重试")
    }
    if !locked {
        global.Logger.Warn("用户正在创建其他订单，请稍后重试")
        return nil, status.Errorf(codes.ResourceExhausted, "正在处理其他订单，请稍后重试")
    }
    
    // 确保在任何情况下都能释放锁
    defer func() {
        if unlockErr := lock.Unlock(ctx); unlockErr != nil {
            global.Logger.Errorf("释放创建订单锁失败: %v", unlockErr)
        }
    }()
    
    // 订单创建逻辑...
}
```

### 订单状态更新锁使用

```go
func (s *OrderServiceServer) OrderUpdate(ctx context.Context, req *proto.OrderStatus) (*emptypb.Empty, error) {
    // 创建分布式锁，防止并发更新同一订单
    var lockKey string
    if req.Id > 0 {
        lockKey = fmt.Sprintf("order_update_lock:id:%d", req.Id)
    } else {
        lockKey = fmt.Sprintf("order_update_lock:sn:%s", req.OrderSn)
    }
    
    lock := utils.NewRedisLock(lockKey, 10*time.Second)
    locked, err := lock.TryLock(ctx, 3, 50*time.Millisecond)
    if err != nil {
        global.Logger.Errorf("获取订单更新锁失败: %v", err)
        return nil, status.Errorf(codes.Internal, "系统忙，请稍后重试")
    }
    if !locked {
        global.Logger.Warn("订单正在被其他请求更新，请稍后重试")
        return nil, status.Errorf(codes.ResourceExhausted, "订单正在处理中，请稍后重试")
    }
    
    defer func() {
        if unlockErr := lock.Unlock(ctx); unlockErr != nil {
            global.Logger.Errorf("释放订单更新锁失败: %v", unlockErr)
        }
    }()
    
    // 订单更新逻辑...
}
```

## 异常处理机制

### 1. 锁获取失败
```go
if err != nil {
    global.Logger.Errorf("获取锁失败: %v", err)
    return nil, status.Errorf(codes.Internal, "系统忙，请稍后重试")
}
```

### 2. 锁被占用
```go
if !locked {
    global.Logger.Warn("锁被占用，请稍后重试")
    return nil, status.Errorf(codes.ResourceExhausted, "正在处理中，请稍后重试")
}
```

### 3. 锁释放失败
```go
defer func() {
    if unlockErr := lock.Unlock(ctx); unlockErr != nil {
        global.Logger.Errorf("释放锁失败: %v", unlockErr)
        // 记录告警，但不影响主流程
    }
}()
```

## 监控与告警

### 关键指标

1. **锁获取成功率**
   - 指标名: `lock_acquire_success_rate`
   - 计算公式: 成功获取锁次数 / 总尝试次数
   - 告警阈值: < 95%

2. **锁获取耗时**
   - 指标名: `lock_acquire_duration`
   - 统计维度: P50, P90, P99
   - 告警阈值: P99 > 500ms

3. **锁持有时长**
   - 指标名: `lock_hold_duration`
   - 用途: 检测是否存在长时间持有锁的情况
   - 告警阈值: > 配置时长的80%

4. **锁释放失败率**
   - 指标名: `lock_release_failure_rate`
   - 计算公式: 释放锁失败次数 / 总释放次数
   - 告警阈值: > 1%

### 日志记录

```go
// 锁获取成功
global.Logger.Infof("成功获取分布式锁: %s, 耗时: %dms", lockKey, duration.Milliseconds())

// 锁获取失败
global.Logger.Warnf("获取分布式锁失败: %s, 重试次数: %d", lockKey, retryCount)

// 锁释放
global.Logger.Debugf("释放分布式锁: %s, 持有时长: %dms", lockKey, holdDuration.Milliseconds())
```

## 性能优化

### 1. 连接池优化
```go
// Redis连接池配置
&redis.Options{
    PoolSize:     10,           // 连接池大小
    MinIdleConns: 5,            // 最小空闲连接数
    MaxRetries:   3,            // 最大重试次数
    DialTimeout:  5 * time.Second,
    ReadTimeout:  3 * time.Second,
    WriteTimeout: 3 * time.Second,
}
```

### 2. 锁粒度优化
- **用户级别锁定**: 不同用户的操作不会相互影响
- **订单级别锁定**: 同一订单的不同操作串行化
- **避免全局锁**: 防止系统整体性能下降

### 3. 超时时间调优
- **创建操作**: 30秒（包含复杂的跨服务调用）
- **更新操作**: 10秒（相对简单的数据库操作）
- **删除操作**: 15秒（批量删除需要一定时间）

## 容错设计

### 1. 锁超时处理
```go
ctx, cancel := context.WithTimeout(context.Background(), 35*time.Second)
defer cancel()
```

### 2. 网络异常处理
```go
if isNetworkError(err) {
    // 网络异常时的重试逻辑
    return retryWithBackoff(operation, maxRetries)
}
```

### 3. Redis不可用处理
```go
if isRedisUnavailable(err) {
    // 降级处理：记录日志但继续执行
    global.Logger.Errorf("Redis不可用，跳过分布式锁: %v", err)
    // 可以考虑使用本地锁或直接执行
}
```

## 最佳实践

### 1. 锁键设计
- **有意义的前缀**: 如 `order_create_lock:`
- **包含关键标识**: 如用户ID、订单ID
- **避免冲突**: 不同业务使用不同前缀

### 2. 超时时间设置
- **充足但不过长**: 确保业务完成但避免死锁
- **考虑网络延迟**: 包含跨服务调用的网络时间
- **监控实际耗时**: 根据监控数据调整参数

### 3. 错误处理
- **区分错误类型**: 网络错误、业务错误、系统错误
- **合理的重试策略**: 避免无谓重试导致系统压力
- **详细的日志记录**: 便于问题排查和性能优化

### 4. 资源释放
- **使用defer**: 确保锁的释放
- **异常处理**: 即使释放失败也要记录日志
- **超时保护**: 避免死锁情况

## 故障排查指南

### 常见问题

1. **锁获取超时**
   - 检查Redis连接状态
   - 查看锁的持有时间是否过长
   - 确认是否存在死锁情况

2. **锁释放失败**
   - 检查网络连通性
   - 确认锁的唯一标识是否正确
   - 查看Redis服务状态

3. **性能问题**
   - 监控锁获取耗时
   - 检查重试次数和间隔设置
   - 分析锁竞争情况

### 排查工具

```bash
# 查看Redis中的锁信息
redis-cli keys "order_*_lock:*"

# 检查锁的TTL
redis-cli ttl "order_create_lock:12345"

# 查看锁的值
redis-cli get "order_create_lock:12345"
```

## 相关文件

- `order_srv/utils/redis_lock.go`: 分布式锁实现
- `order_srv/handler/order.go`: 锁使用示例
- `order_srv/global/global.go`: Redis客户端配置