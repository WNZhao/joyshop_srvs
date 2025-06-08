# 库存服务（order_srv）

本模块为 JoyShop 微服务架构中的库存服务，基于 Go 语言实现，采用 gRPC 通信协议，支持 Consul 服务注册与发现，集成 zap 日志、viper 配置、GORM ORM、自动数据库迁移及标准化测试用例。特别实现了高并发场景下的库存管理机制。

## 目录结构

```
order_srv/
├── config/         # 配置文件与配置结构体
│   ├── config-develop.yaml
│   ├── config-sample.yaml
│   └── config.go
├── global/         # 全局变量与全局配置
│   └── global.go
├── handler/        # gRPC 业务处理器（handler）
│   └── inventory.go
├── initialize/     # 初始化相关（配置、数据库、日志、Consul注册等）
│   ├── config.go
│   ├── consul.go
│   ├── db.go
│   ├── logger.go
│   └── redis.go
├── model/          # 数据模型
│   └── inventory.go
├── proto/          # gRPC 协议定义及生成代码
│   ├── inventory.proto
│   ├── inventory.pb.go
│   └── inventory_grpc.pb.go
├── tests/          # 单元测试用例
│   └── inventory_test.go
├── util/           # 工具函数
│   └── config.go
└── main.go         # 服务启动入口
```

## 创建步骤

### 1. 基础目录结构创建
- 创建主要目录：config、global、handler、initialize、model、proto、tests、util
- 复制并修改 .gitignore 文件

### 2. 配置文件设置
- 复制并修改 config-sample.yaml 为库存服务配置
- 创建 config.go 定义配置结构体

### 3. 全局变量设置
- 创建 global.go 定义全局变量和配置

### 4. 初始化模块
- 复制并修改 initialize 目录下的初始化逻辑
- 包括：配置、数据库、日志、Consul 注册、Redis 连接等

### 5. 库存服务定义
- 创建 proto/inventory.proto 定义库存服务接口
- 生成 gRPC 代码

### 6. 数据模型
- 创建 model/inventory.go 定义库存相关模型
- 实现数据库操作相关方法

### 7. 业务处理器
- 创建 handler/inventory.go 实现库存服务接口
- 实现库存相关的业务逻辑，包括并发控制

### 8. 测试用例
- 创建 tests/inventory_test.go 编写测试用例
- 覆盖库存服务的主要功能，包括并发测试

### 9. 主程序
- 创建 main.go 作为服务启动入口
- 实现服务注册和启动逻辑

## 功能说明

### 1. gRPC 库存服务接口

- 接口定义见 `proto/inventory.proto`，主要包括：
  - `SetInventory`：设置商品库存
  - `GetInventory`：获取商品库存
  - `Sell`：库存扣减（支持并发控制）
  - `Reback`：库存归还（支持并发控制）

### 2. 并发控制机制

#### 2.1 Redis分布式锁实现
- **分布式锁机制**：
  - 使用 Redis SETNX 实现分布式锁
  - 锁超时时间：10秒
  - 自动释放锁机制
  - 重试机制（最多3次）

- **锁的获取与释放**：
  - 使用商品ID作为锁的key
  - 确保锁的互斥性
  - 使用 defer 确保锁的释放

- **事务控制**：
  - 结合分布式锁和数据库事务
  - 确保数据一致性
  - 防止并发更新冲突

#### 2.2 乐观锁实现（保留供参考）
- **乐观锁机制**：
  - 使用 `version` 字段进行并发控制
  - 在更新时检查版本号是否匹配
  - 通过 `RowsAffected` 判断更新是否成功

- **重试机制**：
  - 最大重试次数：3次
  - 重试延迟：100ms
  - 自动处理并发冲突

- **事务控制**：
  - 使用数据库事务确保原子性
  - 结合行级锁和乐观锁
  - 精确的并发冲突检测

- **代码实现示例**：
  ```go
  // Sell 库存扣减 - 乐观锁实现
  func (s *InventoryServer) Sell(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
      const maxRetries = 3
      const retryDelay = 100 * time.Millisecond

      for retry := 0; retry < maxRetries; retry++ {
          // 开启事务
          tx := global.DB.Begin()
          if tx.Error != nil {
              return nil, status.Error(codes.Internal, "开启事务失败")
          }

          success := true
          // 遍历所有商品
          for _, goodsInfo := range req.GoodsInvInfo {
              var inv model.Inventory
              // 使用行锁查询
              result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
                  Where("goods_id = ?", goodsInfo.GoodsId).
                  First(&inv)

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

              // 使用乐观锁更新库存
              updateResult := tx.Model(&model.Inventory{}).
                  Select("Stock", "Version").
                  Where("goods_id = ? AND version = ?", goodsInfo.GoodsId, inv.Version).
                  Updates(map[string]interface{}{
                      "stock":   inv.Stock - goodsInfo.Num,
                      "version": inv.Version + 1,
                  })

              if updateResult.Error != nil {
                  tx.Rollback()
                  success = false
                  break
              }

              if updateResult.RowsAffected == 0 {
                  tx.Rollback()
                  success = false
                  break
              }
          }

          if success {
              // 提交事务
              if err := tx.Commit().Error; err != nil {
                  return nil, status.Error(codes.Internal, "提交事务失败")
              }
              return &emptypb.Empty{}, nil
          }

          // 如果失败且还有重试次数，等待后重试
          if retry < maxRetries-1 {
              time.Sleep(retryDelay)
              continue
          }
      }

      return nil, status.Error(codes.Internal, "库存更新失败，请重试")
  }

  // Reback 库存归还 - 乐观锁实现
  func (s *InventoryServer) Reback(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
      const maxRetries = 3
      const retryDelay = 100 * time.Millisecond

      for retry := 0; retry < maxRetries; retry++ {
          // 开启事务
          tx := global.DB.Begin()
          if tx.Error != nil {
              return nil, status.Error(codes.Internal, "开启事务失败")
          }

          success := true
          // 遍历所有商品
          for _, goodsInfo := range req.GoodsInvInfo {
              var inv model.Inventory
              // 使用行锁查询
              result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
                  Where("goods_id = ?", goodsInfo.GoodsId).
                  First(&inv)

              if result.Error != nil {
                  tx.Rollback()
                  if result.Error == gorm.ErrRecordNotFound {
                      return nil, status.Error(codes.NotFound, "商品库存不存在")
                  }
                  return nil, status.Error(codes.Internal, "查询库存失败")
              }

              // 使用乐观锁更新库存
              updateResult := tx.Model(&model.Inventory{}).
                  Select("Stock", "Version").
                  Where("goods_id = ? AND version = ?", goodsInfo.GoodsId, inv.Version).
                  Updates(map[string]interface{}{
                      "stock":   inv.Stock + goodsInfo.Num,
                      "version": inv.Version + 1,
                  })

              if updateResult.Error != nil {
                  tx.Rollback()
                  success = false
                  break
              }

              if updateResult.RowsAffected == 0 {
                  tx.Rollback()
                  success = false
                  break
              }
          }

          if success {
              // 提交事务
              if err := tx.Commit().Error; err != nil {
                  return nil, status.Error(codes.Internal, "提交事务失败")
              }
              return &emptypb.Empty{}, nil
          }

          // 如果失败且还有重试次数，等待后重试
          if retry < maxRetries-1 {
              time.Sleep(retryDelay)
              continue
          }
      }

      return nil, status.Error(codes.Internal, "库存更新失败，请重试")
  }
  ```

- **实现要点说明**：
  1. **版本控制**：
     - 使用 `version` 字段记录数据版本
     - 每次更新时版本号加1
     - 通过版本号判断数据是否被修改

  2. **行级锁**：
     - 使用 `clause.Locking{Strength: "UPDATE"}` 实现行级锁
     - 防止并发读取时的数据不一致
     - 确保事务的隔离性

  3. **事务管理**：
     - 使用数据库事务确保原子性
     - 出错时自动回滚
     - 成功时提交事务

  4. **重试机制**：
     - 最多重试3次
     - 重试间隔100ms
     - 自动处理并发冲突

  5. **错误处理**：
     - 详细的错误类型判断
     - 清晰的错误信息
     - 完善的日志记录

- **优缺点分析**：
  1. **优点**：
     - 实现简单，无需额外组件
     - 性能开销较小
     - 适合并发量不大的场景

  2. **缺点**：
     - 不适合分布式环境
     - 重试机制可能增加数据库压力
     - 在高并发场景下性能可能下降

### 3. 配置管理

- 配置结构体定义于 `config/config.go`，通过 viper 统一加载，支持多环境配置。
- 配置内容包括服务信息、数据库、Consul、Redis、日志等。

### 4. 日志系统

- 集成 zap 日志，初始化于 `initialize/logger.go`，全局统一调用 `zap.S()` 进行日志输出。

### 5. 数据库与模型

- 使用 GORM 作为 ORM，模型定义于 `model/inventory.go`。
- 启动时自动迁移表结构，确保数据库与模型同步。

### 6. Consul 服务注册

- 通过 `initialize/consul.go` 实现 Consul 注册，支持健康检查。

### 7. Redis 连接

- 通过 `initialize/redis.go` 实现 Redis 连接
- 支持连接池配置
- 支持超时控制
- 支持健康检查

### 8. 测试用例

- 所有测试用例集中于 `tests/inventory_test.go`，覆盖库存服务的主要功能。
- 包含并发测试场景，验证库存管理的正确性。

## 启动方式

1. **配置准备**  
   复制 `config/config-sample.yaml` 为 `config-develop.yaml`，根据实际环境修改配置。

2. **数据库准备**  
   确保数据库可用，服务启动时会自动迁移表结构。

3. **Redis准备**  
   确保 Redis 服务可用，用于分布式锁实现。

4. **启动服务**  
   ```bash
   go run main.go
   ```

5. **运行测试**  
   ```bash
   go test ./tests
   ```

## 维护建议

- **配置变更**：统一修改 `config/` 下的 yaml 文件及结构体。
- **日志与监控**：优先使用 zap 日志，便于后续接入监控系统。
- **接口扩展**：新增 gRPC 接口时，先修改 `proto/inventory.proto`，再生成代码并实现 handler。
- **测试规范**：所有功能变更需补充/完善 `inventory_test.go`，保证核心流程可回归。
- **并发控制**：修改库存相关逻辑时，确保保持并发控制机制的正确性。
- **Redis维护**：定期检查 Redis 连接状态，确保分布式锁的可靠性。

## 更新日志

### 2024-03-21
- 初始化项目结构
- 创建基础目录
- 实现库存服务基础功能
- 添加 Redis 分布式锁实现
- 优化库存扣减和归还逻辑
- 编写 README.md 