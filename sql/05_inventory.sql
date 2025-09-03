-- 库存服务测试数据
-- 为所有商品创建对应的库存信息

-- 清空现有数据
DELETE FROM inventory;

-- 重置自增ID
ALTER TABLE inventory AUTO_INCREMENT = 1;

-- 插入库存数据，对应商品表中的商品
INSERT INTO inventory (id, goods_id, stock, version, created_at, updated_at, is_deleted) VALUES
-- 手机类商品库存（高价值商品，库存相对较少）
(1, 1, 156, 0, NOW(), NOW(), false),  -- iPhone 15 Pro Max - 适中库存
(2, 2, 234, 0, NOW(), NOW(), false),  -- iPhone 14 - 较多库存
(3, 3, 89, 0, NOW(), NOW(), false),   -- HUAWEI Mate 60 Pro - 较少库存（热门）
(4, 4, 167, 0, NOW(), NOW(), false),  -- HUAWEI P60 Pro - 适中库存
(5, 5, 78, 0, NOW(), NOW(), false),   -- Xiaomi 14 Ultra - 较少库存（热门）

-- 电脑类商品库存
(6, 6, 45, 0, NOW(), NOW(), false),   -- MacBook Pro - 较少库存（高价值）
(7, 7, 67, 0, NOW(), NOW(), false),   -- ThinkPad X1 - 适中库存

-- 鞋类商品库存（按尺码分布）
(8, 8, 234, 0, NOW(), NOW(), false),  -- Nike Air Jordan 1 - 多库存（热门款）
(9, 9, 187, 0, NOW(), NOW(), false),  -- Adidas Ultraboost - 适中库存

-- 服装类商品库存（按尺码库存较多）
(10, 10, 456, 0, NOW(), NOW(), false), -- UNIQLO T恤 - 大量库存（基础款）
(11, 11, 123, 0, NOW(), NOW(), false), -- Levi's 牛仔裤 - 适中库存

-- 家电类商品库存
(12, 12, 89, 0, NOW(), NOW(), false),  -- 美的空调 - 季节性商品，适中库存
(13, 13, 76, 0, NOW(), NOW(), false),  -- 海尔冰箱 - 适中库存

-- 相机类商品库存（专业设备）
(14, 14, 23, 0, NOW(), NOW(), false),  -- Canon EOS R6 Mark II - 少量库存（专业设备）

-- 图书类商品库存（按需印刷，库存较多）
(15, 15, 567, 0, NOW(), NOW(), false), -- 深入理解计算机系统 - 教材，大量库存
(16, 16, 789, 0, NOW(), NOW(), false), -- Python编程 - 热门教材，大量库存

-- 美妆类商品库存
(17, 17, 145, 0, NOW(), NOW(), false), -- SK-II 神仙水 - 适中库存（进口商品）

-- 家具类商品库存（大件商品）
(18, 18, 34, 0, NOW(), NOW(), false),  -- IKEA 床架 - 较少库存（大件物流）

-- 电视类商品库存
(19, 19, 67, 0, NOW(), NOW(), false),  -- 小米电视65英寸 - 适中库存

-- 文具类商品库存（日常消耗品）
(20, 20, 2340, 0, NOW(), NOW(), false), -- 晨光中性笔 - 大量库存（消耗品）

-- 添加一些额外的库存数据，模拟多SKU商品
-- 为热门商品添加不同规格的库存记录

-- iPhone 15 Pro Max 不同存储容量（模拟多SKU）
(21, 1, 89, 0, NOW(), NOW(), false),   -- 假设这是128GB版本的库存
(22, 1, 67, 0, NOW(), NOW(), false),   -- 假设这是1TB版本的库存

-- 小米14 Ultra 不同颜色库存
(23, 5, 45, 0, NOW(), NOW(), false),   -- 假设白色版本库存
(24, 5, 34, 0, NOW(), NOW(), false),   -- 假设黑色版本库存

-- Nike AJ1 不同尺码库存分布（模拟尺码库存）
(25, 8, 25, 0, NOW(), NOW(), false),   -- 假设40码库存
(26, 8, 34, 0, NOW(), NOW(), false),   -- 假设41码库存  
(27, 8, 45, 0, NOW(), NOW(), false),   -- 假设42码库存（热门尺码）
(28, 8, 56, 0, NOW(), NOW(), false),   -- 假设43码库存
(29, 8, 23, 0, NOW(), NOW(), false),   -- 假设44码库存

-- UNIQLO T恤不同尺码库存
(30, 10, 89, 0, NOW(), NOW(), false),  -- S码库存
(31, 10, 123, 0, NOW(), NOW(), false), -- M码库存（热门尺码）
(32, 10, 145, 0, NOW(), NOW(), false), -- L码库存（热门尺码）
(33, 10, 78, 0, NOW(), NOW(), false),  -- XL码库存
(34, 10, 45, 0, NOW(), NOW(), false),  -- XXL码库存

-- 一些库存告警的商品（库存不足）
(35, 3, 5, 0, NOW(), NOW(), false),    -- HUAWEI Mate 60 Pro 某个版本库存不足
(36, 5, 3, 0, NOW(), NOW(), false),    -- Xiaomi 14 Ultra 某个颜色库存不足
(37, 14, 2, 0, NOW(), NOW(), false),   -- Canon 相机库存极少

-- 一些库存充足的商品
(38, 15, 999, 0, NOW(), NOW(), false), -- 计算机系统书籍库存充足
(39, 16, 888, 0, NOW(), NOW(), false), -- Python编程书籍库存充足
(40, 20, 5000, 0, NOW(), NOW(), false); -- 中性笔库存非常充足

-- 创建一些历史库存变动数据（可选，用于测试库存历史功能）
-- 注意：inventory_history 表结构需要先确认字段名是否正确
/*
INSERT INTO inventory_history (user, goods, nums, order, state, created_at, updated_at) VALUES
-- 用户2购买iPhone的历史记录
(2, 1, 1, 1001, 2, DATE_SUB(NOW(), INTERVAL 7 DAY), DATE_SUB(NOW(), INTERVAL 7 DAY)),
(2, 2, 1, 1002, 2, DATE_SUB(NOW(), INTERVAL 5 DAY), DATE_SUB(NOW(), INTERVAL 5 DAY)),

-- 用户3购买运动鞋的历史记录  
(3, 8, 1, 1003, 2, DATE_SUB(NOW(), INTERVAL 3 DAY), DATE_SUB(NOW(), INTERVAL 3 DAY)),
(3, 9, 1, 1004, 1, DATE_SUB(NOW(), INTERVAL 1 DAY), DATE_SUB(NOW(), INTERVAL 1 DAY)),

-- 用户4购买家电的历史记录
(4, 12, 1, 1005, 2, DATE_SUB(NOW(), INTERVAL 10 DAY), DATE_SUB(NOW(), INTERVAL 10 DAY)),
(4, 13, 1, 1006, 2, DATE_SUB(NOW(), INTERVAL 8 DAY), DATE_SUB(NOW(), INTERVAL 8 DAY));
*/

-- 说明：
-- stock: 库存数量，根据商品类型设置不同的库存量
--   - 高价值商品（手机、电脑、相机）：相对较少库存
--   - 服装、图书、文具：较多库存  
--   - 家电、家具：中等库存
--   - 消耗品（文具）：大量库存
-- version: 版本号，用于乐观锁控制，初始为0
-- 部分商品设置了多个库存记录，模拟多SKU（不同颜色、尺码、规格）
-- 包含了库存告警测试数据（库存不足）和库存充足的测试数据
-- 库存数据与商品数据保持一致，确保联动测试的完整性