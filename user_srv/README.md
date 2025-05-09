# 用户服务（user_srv）

本模块为 JoyShop 微服务架构中的用户服务，基于 Go 语言实现，采用 gRPC 通信协议，支持 Consul 服务注册与发现，集成 zap 日志、viper 配置、GORM ORM、自动数据库迁移及标准化测试用例。

## 目录结构

```
user_srv/
├── config/         # 配置文件与配置结构体
│   ├── config-develop.yaml
│   ├── config-sample.yaml
│   └── config.go
├── global/         # 全局变量与全局配置
│   └── global.go
├── handler/        # gRPC 业务处理器（handler）
│   └── user.go
├── initialize/     # 初始化相关（配置、数据库、日志、Consul注册等）
│   ├── config.go
│   ├── consul.go
│   ├── db.go
│   └── logger.go
├── model/          # 数据模型
│   ├── user.go
│   └── main/
│       ├── main.go
│       └── md5/
│           └── main.go
├── proto/          # gRPC 协议定义及生成代码
│   ├── user.proto
│   ├── user.pb.go
│   └── user_grpc.pb.go
├── tests/          # 单元测试用例
│   └── user_test.go
├── util/           # 工具函数
│   └── config.go
└── main.go         # 服务启动入口
```

## 功能说明

### 1. gRPC 用户服务接口

- 接口定义见 `proto/user.proto`，主要包括：
  - `GetUserList`：分页获取用户列表
  - `GetUserByMobile`：通过手机号查找用户
  - `GetUserById`：通过ID查找用户
  - `CreateUser`：创建新用户（含密码加密）
  - `UpdateUser`：更新用户信息
  - `DeleteUser`：删除用户
  - `CheckPassword`：校验密码

### 2. 配置管理

- 配置结构体定义于 `config/config.go`，通过 viper 统一加载，支持多环境配置（如 `config-develop.yaml`）。
- 配置内容包括服务信息、数据库、Consul、日志等。

### 3. 日志系统

- 集成 zap 日志，初始化于 `initialize/logger.go`，全局统一调用 `zap.S()` 进行日志输出，支持多级别日志。

### 4. 数据库与模型

- 使用 GORM 作为 ORM，模型定义于 `model/user.go`。
- 启动时自动迁移表结构，确保数据库与模型同步。

### 5. Consul 服务注册

- 通过 `initialize/consul.go` 实现 Consul 注册，支持健康检查，便于服务发现与治理。

### 6. 测试用例

- 所有测试用例集中于 `tests/user_test.go`，覆盖用户服务的主要功能。
- 测试用例自动初始化配置、数据库、gRPC 连接，并生成唯一性测试数据，避免因重复导致的失败。
- 推荐使用 `go test ./tests` 运行全部测试。

### 7. 工具与辅助

- `util/` 目录下存放通用工具函数，如配置辅助加载等。

## 启动方式

1. **配置准备**  
   复制 `config/config-sample.yaml` 为 `config-develop.yaml`，根据实际环境修改数据库、Consul 等参数。

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
- **日志与监控**：优先使用 zap 日志，便于后续接入 ELK、Prometheus 等监控系统。
- **接口扩展**：新增 gRPC 接口时，先修改 `proto/user.proto`，再生成代码并实现 handler。
- **测试规范**：所有功能变更需补充/完善 `user_test.go`，保证核心流程可回归。

---

如需详细开发文档或接口说明，请参考各目录下源码及注释。 