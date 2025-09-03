# JoyShop 微服务测试数据说明

## 概述

本目录包含了 JoyShop 微服务系统的完整测试数据集，用于支持各服务的开发、测试和演示。数据设计真实可信，涵盖了电商业务的各个环节。

## 📁 文件结构

```
sql/
├── README.md                           # 本说明文件
├── create_databases.sql                # 🔧 创建微服务数据库
├── init_microservices.sql              # 🚀 微服务架构初始化（推荐）
├── test_microservices_integration.sql # 🧪 微服务联动测试
├── init_all.sql                       # 🚀 单体架构初始化
├── test_integration.sql               # 🧪 单体架构联动测试
├── test_basic_queries.sql             # 🧪 基础功能测试查询
├── user_srv_init.sql        # 👤 用户服务专用初始化
├── goods_srv_init.sql       # 🛍️ 商品服务专用初始化  
├── inventory_srv_init.sql   # 📦 库存服务专用初始化
├── order_srv_init.sql       # 📋 订单服务专用初始化
├── 01_users.sql             # 用户基础数据
├── 02_categories.sql        # 商品分类数据（3级分类）
├── 03_brands.sql           # 品牌数据及分类关联
├── 04_goods.sql            # 商品数据及分类关联
├── 05_inventory.sql        # 库存数据
├── 06_shopping_cart.sql    # 购物车数据
├── 07_orders.sql           # 订单数据
└── 08_banners.sql                     # 轮播图数据
```

## 🚀 快速开始

### 方法1: 微服务架构（推荐生产环境）

```bash
# 1. 创建微服务数据库
mysql -u username -p
mysql> source sql/create_databases.sql;

# 2. 初始化微服务数据
mysql> source sql/init_microservices.sql;

# 3. 运行微服务集成测试
mysql> source sql/test_microservices_integration.sql;
```

### 方法2: 单体架构（开发测试）

```bash
# 1. 连接数据库
mysql -u username -p database_name

# 2. 执行完整初始化
mysql> source sql/init_all.sql;
```

### 方法3: 按服务分别初始化（单体架构）

```bash
# 只初始化用户服务数据
mysql> source sql/user_srv_init.sql;

# 只初始化商品服务数据  
mysql> source sql/goods_srv_init.sql;

# 只初始化库存服务数据
mysql> source sql/inventory_srv_init.sql;

# 只初始化订单服务数据
mysql> source sql/order_srv_init.sql;
```

### 方法4: 手动逐个执行

```bash
mysql> source sql/01_users.sql;
mysql> source sql/02_categories.sql;
mysql> source sql/03_brands.sql;
mysql> source sql/04_goods.sql;
mysql> source sql/05_inventory.sql;
mysql> source sql/06_shopping_cart.sql;
mysql> source sql/07_orders.sql;
mysql> source sql/08_banners.sql;
```

## 🧪 测试验证

### 微服务架构测试
初始化完成后，运行微服务联动测试脚本：

```bash
mysql> source sql/test_microservices_integration.sql;
```

该脚本会检查：
- ✅ 跨数据库连接性测试
- ✅ 跨服务数据ID一致性验证
- ✅ 微服务业务逻辑验证
- ✅ 分布式事务模拟测试
- ✅ 微服务配置验证建议

### 单体架构测试
```bash
mysql> source sql/test_integration.sql;
```

该脚本会检查：
- ✅ 数据一致性（商品-库存-订单关联）
- ✅ 业务逻辑正确性（库存充足性、金额计算）
- ✅ 用户行为路径完整性
- ✅ 数据质量（无异常数据）

## 📊 数据概览

### 用户数据 (15个用户)
| 用户类型 | 数量 | 说明 |
|---------|------|------|
| 管理员 | 1个 | ID=1, 用户名: admin |
| 普通用户 | 12个 | ID=2-12, 不同活跃度 |
| VIP用户 | 2个 | ID=13-14, 高价值用户 |

**登录信息**: 所有用户密码统一为 `123456`

### 商品数据 (20+个商品)
| 分类 | 商品数量 | 价格范围 | 热门品牌 |
|------|----------|----------|----------|
| 电子数码 | 10+ | ¥69 - ¥31,067 | Apple, HUAWEI, 小米 |
| 服装鞋包 | 4+ | ¥79 - ¥1,299 | Nike, Adidas, UNIQLO |
| 家居生活 | 3+ | ¥699 - ¥4,199 | IKEA, 美的, 海尔 |
| 图书文教 | 2+ | ¥19.8 - ¥99 | 人邮社, 机工社 |
| 美妆个护 | 1+ | ¥1,390 | SK-II |

### 库存分布
- 📦 **总库存记录**: 40+ 条
- 📈 **库存范围**: 2 - 5,000 件
- ⚠️ **低库存商品**: 5+ 个（库存 ≤ 10）
- 🔥 **充足库存**: 消耗品类库存充足

### 订单状态分布
| 状态 | 数量 | 说明 |
|------|------|------|
| TRADE_FINISHED | 2个 | 已完成订单 |
| TRADE_SUCCESS | 4个 | 支付成功订单 |
| WAIT_BUYER_PAY | 2个 | 待支付订单 |
| PAYING | 1个 | 支付中订单 |
| TRADE_CLOSED | 1个 | 已关闭订单 |

## 🎯 测试场景

### 场景1: 新用户注册登录
```
用户: testuser (ID=15)
密码: 123456
邮箱: testuser@example.com
```

### 场景2: 商品浏览购买流程
```
1. 浏览分类: /category/1 (电子数码)
2. 查看商品: iPhone 15 Pro Max (ID=1)
3. 检查库存: 156 件可用
4. 加入购物车: 数量 1
5. 创建订单: 测试批量处理
```

### 场景3: 购物车批量操作
```
用户: 张三 (ID=2)
购物车状态: 3件已选中，1件未选中
测试: 批量结算、库存检查、订单创建
```

### 场景4: VIP用户大额订单
```
用户: VIP刘总 (ID=13)  
购物车: 多件高价值商品
订单金额: ¥45,876
测试: 大额订单处理、VIP服务
```

### 场景5: 库存不足处理
```
商品: HUAWEI Mate 60 Pro (库存仅5件)
测试: 库存告警、超卖防护、库存回滚
```

### 场景6: 订单状态流转
```
订单流程: 待支付 → 支付中 → 支付成功 → 交易完成
测试: 状态机验证、分布式锁、事务一致性
```

## 🔧 开发建议

### 1. 微服务架构开发
- **数据库隔离**: 每个服务连接独立数据库，避免跨库直接查询
- **API调用**: 服务间通过API调用，不直接访问其他服务数据库  
- **数据一致性**: 使用分布式事务或最终一致性保证数据同步
- **配置管理**: 各服务配置独立的数据库连接信息

### 2. 接口测试
- **用户服务**: 连接joyshop_user，使用ID 2-15测试用户相关接口
- **商品服务**: 连接joyshop_goods，测试分类查询、商品搜索、商品详情
- **库存服务**: 连接joyshop_inventory，测试库存查询、扣减、回滚操作
- **订单服务**: 连接joyshop_order，测试购物车操作、订单CRUD、状态流转

### 3. 性能测试
- **并发下单**: 使用多个用户同时创建订单
- **库存压测**: 模拟高并发库存扣减
- **批量处理**: 测试大购物车的批量插入性能

### 4. 异常测试
- **库存不足**: 尝试购买超过库存数量的商品
- **重复下单**: 测试分布式锁防重复机制
- **数据一致性**: 模拟服务异常时的数据回滚

## 📈 监控指标

### 关键业务指标
- **订单转化率**: ~66.7% (基于测试数据)
- **支付成功率**: ~87.5% 
- **平均订单金额**: ¥13,456
- **热销商品占比**: 35%

### 系统性能指标
- **响应时间**: 订单创建 < 2s
- **并发处理**: 支持 100+ TPS
- **数据一致性**: 99.9%

## ⚠️ 注意事项

1. **数据依赖**: 按照文档顺序执行初始化，确保外键关系正确
2. **字符编码**: 使用 UTF8MB4 编码支持emoji和特殊字符
3. **外键约束**: 脚本会临时禁用外键检查，提高导入速度
4. **数据清理**: 每次执行会清空现有数据，请备份重要数据
5. **图片链接**: 商品图片使用模拟链接，实际项目需要替换

## 🤝 贡献指南

如需添加新的测试数据或修改现有数据：

1. **保持数据真实性**: 使用真实的商品名称、价格、品牌
2. **维护关联关系**: 确保商品-库存-订单数据一致
3. **更新统计信息**: 修改数据后更新README中的统计数字
4. **测试验证**: 运行 `test_integration.sql` 验证数据完整性

## 📞 支持

如有问题或建议，请：
- 📧 提交 Issue 到项目仓库
- 📚 查看项目文档
- 🔍 运行联动测试脚本排查问题

---

🎉 **开始愉快的测试吧！** 这套数据将帮助您快速验证微服务间的协作和业务逻辑的正确性。