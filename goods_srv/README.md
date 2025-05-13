# 商品服务（goods_srv）

本模块为 JoyShop 微服务架构中的商品服务，基于 Go 语言实现，采用 gRPC 通信协议，支持 Consul 服务注册与发现，集成 zap 日志、viper 配置、GORM ORM、自动数据库迁移及标准化测试用例。

## 目录结构

```
goods_srv/
├── config/         # 配置文件与配置结构体
│   ├── config-develop.yaml
│   ├── config-sample.yaml
│   └── config.go
├── global/         # 全局变量与全局配置
│   └── global.go
├── handler/        # gRPC 业务处理器（handler）
│   └── goods.go
├── initialize/     # 初始化相关（配置、数据库、日志、Consul注册等）
│   ├── config.go
│   ├── consul.go
│   ├── db.go
│   └── logger.go
├── model/          # 数据模型
│   ├── goods.go
│   └── main/
│       ├── main.go
│       └── md5/
│           └── main.go
├── proto/          # gRPC 协议定义及生成代码
│   ├── goods.proto
│   ├── goods.pb.go
│   └── goods_grpc.pb.go
├── tests/          # 单元测试用例
│   └── goods_test.go
├── util/           # 工具函数
│   └── config.go
└── main.go         # 服务启动入口
```

## 创建步骤

### 1. 基础目录结构创建
- 创建主要目录：config、global、handler、initialize、model、proto、tests、util
- 复制并修改 .gitignore 文件

### 2. 配置文件设置
- 复制并修改 config-sample.yaml 为商品服务配置
- 创建 config.go 定义配置结构体

### 3. 全局变量设置
- 创建 global.go 定义全局变量和配置

### 4. 初始化模块
- 复制并修改 initialize 目录下的初始化逻辑
- 包括：配置、数据库、日志、Consul 注册等

### 5. 商品服务定义
- 创建 proto/goods.proto 定义商品服务接口
- 生成 gRPC 代码

### 6. 数据模型
- 创建 model/goods.go 定义商品相关模型
- 实现数据库操作相关方法

### 7. 业务处理器
- 创建 handler/goods.go 实现商品服务接口
- 实现商品相关的业务逻辑

### 8. 测试用例
- 创建 tests/goods_test.go 编写测试用例
- 覆盖商品服务的主要功能

### 9. 主程序
- 创建 main.go 作为服务启动入口
- 实现服务注册和启动逻辑

## 功能说明

### 1. gRPC 商品服务接口

- 接口定义见 `proto/goods.proto`，主要包括：
  - `GetGoodsList`：分页获取商品列表
  - `GetGoodsById`：通过ID查找商品
  - `CreateGoods`：创建新商品
  - `UpdateGoods`：更新商品信息
  - `DeleteGoods`：删除商品
  - `GetGoodsByCategory`：按分类获取商品
  - `SearchGoods`：搜索商品

### 2. 配置管理

- 配置结构体定义于 `config/config.go`，通过 viper 统一加载，支持多环境配置。
- 配置内容包括服务信息、数据库、Consul、日志等。

### 3. 日志系统

- 集成 zap 日志，初始化于 `initialize/logger.go`，全局统一调用 `zap.S()` 进行日志输出。

### 4. 数据库与模型

- 使用 GORM 作为 ORM，模型定义于 `model/goods.go`。
- 启动时自动迁移表结构，确保数据库与模型同步。

### 5. Consul 服务注册

- 通过 `initialize/consul.go` 实现 Consul 注册，支持健康检查。

### 6. 测试用例

- 所有测试用例集中于 `tests/goods_test.go`，覆盖商品服务的主要功能。

## 启动方式

1. **配置准备**  
   复制 `config/config-sample.yaml` 为 `config-develop.yaml`，根据实际环境修改配置。

2. **数据库准备**  
   确保数据库可用，服务启动时会自动迁移表结构。

3. **启动服务**  
   ```bash
   go run main.go
   ```

4. **运行测试**  
   ```bash
   go test ./tests
   ```

## 维护建议

- **配置变更**：统一修改 `config/` 下的 yaml 文件及结构体。
- **日志与监控**：优先使用 zap 日志，便于后续接入监控系统。
- **接口扩展**：新增 gRPC 接口时，先修改 `proto/goods.proto`，再生成代码并实现 handler。
- **测试规范**：所有功能变更需补充/完善 `goods_test.go`，保证核心流程可回归。

## 更新日志

### 2024-03-21
- 初始化项目结构
- 创建基础目录
- 编写 README.md 