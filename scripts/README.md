# 服务管理脚本

## 概述

`service_manager.sh` 是一个用于管理 joyshop_srvs 项目中所有微服务的脚本。它提供了启动、停止、重启、状态查看和日志查看等功能。

## 功能特性

- ✅ **批量管理**：支持同时管理多个服务
- ✅ **单个控制**：支持单独控制某个服务
- ✅ **状态监控**：实时查看服务运行状态
- ✅ **日志查看**：实时查看服务日志
- ✅ **进程管理**：自动管理 PID 文件
- ✅ **端口检查**：启动前检查端口占用
- ✅ **优雅停止**：支持优雅停止和强制停止
- ✅ **易于扩展**：添加新服务只需修改配置

## 使用方法

### 基本语法

```bash
./scripts/service_manager.sh [命令] [服务名]
```

### 可用命令

| 命令      | 说明     | 示例                                          |
| --------- | -------- | --------------------------------------------- |
| `start`   | 启动服务 | `./scripts/service_manager.sh start`          |
| `stop`    | 停止服务 | `./scripts/service_manager.sh stop`           |
| `restart` | 重启服务 | `./scripts/service_manager.sh restart`        |
| `status`  | 查看状态 | `./scripts/service_manager.sh status`         |
| `logs`    | 查看日志 | `./scripts/service_manager.sh logs goods_srv` |
| `clean`   | 清理文件 | `./scripts/service_manager.sh clean`          |
| `help`    | 显示帮助 | `./scripts/service_manager.sh help`           |

### 服务列表

| 服务名          | 端口  | 说明     |
| --------------- | ----- | -------- |
| `goods_srv`     | 50051 | 商品服务 |
| `inventory_srv` | 50052 | 库存服务 |
| `order_srv`     | 50053 | 订单服务 |
| `user_srv`      | 50054 | 用户服务 |

## 使用示例

### 启动所有服务

```bash
./scripts/service_manager.sh start
```

### 启动单个服务

```bash
./scripts/service_manager.sh start goods_srv
```

### 停止所有服务

```bash
./scripts/service_manager.sh stop
```

### 停止单个服务

```bash
./scripts/service_manager.sh stop user_srv
```

### 重启所有服务

```bash
./scripts/service_manager.sh restart
```

### 重启单个服务

```bash
./scripts/service_manager.sh restart order_srv
```

### 查看所有服务状态

```bash
./scripts/service_manager.sh status
```

### 查看单个服务状态

```bash
./scripts/service_manager.sh status inventory_srv
```

### 查看服务日志

```bash
./scripts/service_manager.sh logs goods_srv
```

### 清理日志和 PID 文件

```bash
./scripts/service_manager.sh clean
```

## 目录结构

脚本运行时会自动创建以下目录：

```
joyshop_srvs/
├── logs/           # 服务日志文件
│   ├── goods_srv.log
│   ├── inventory_srv.log
│   ├── order_srv.log
│   └── user_srv.log
├── pids/           # PID 文件
│   ├── goods_srv.pid
│   ├── inventory_srv.pid
│   ├── order_srv.pid
│   └── user_srv.pid
└── scripts/
    └── service_manager.sh
```

## 扩展新服务

要添加新的服务，只需修改脚本中的 `SERVICES` 配置：

```bash
# 在脚本中找到这一行
declare -A SERVICES=(
    ["goods_srv"]="50051"
    ["inventory_srv"]="50052"
    ["order_srv"]="50053"
    ["user_srv"]="50054"
    ["new_service"]="50055"  # 添加新服务
)
```

## 注意事项

1. **端口冲突**：确保新服务的端口不被其他服务占用
2. **目录结构**：新服务必须遵循 `service_name/main.go` 的目录结构
3. **依赖关系**：启动顺序可能需要考虑服务间的依赖关系
4. **权限问题**：确保脚本有执行权限 `chmod +x scripts/service_manager.sh`

## 故障排除

### 服务启动失败

1. 检查端口是否被占用：`lsof -i :端口号`
2. 检查服务目录是否存在
3. 查看服务日志：`./scripts/service_manager.sh logs 服务名`

### 服务停止失败

1. 检查 PID 文件是否存在
2. 手动杀死进程：`kill -9 PID`
3. 清理 PID 文件：`./scripts/service_manager.sh clean`

### 权限问题

```bash
chmod +x scripts/service_manager.sh
```

## 开发建议

1. **日志管理**：定期清理日志文件避免磁盘空间不足
2. **监控集成**：可以集成到监控系统中
3. **自动化**：可以结合 CI/CD 流程使用
4. **备份**：重要配置建议定期备份
