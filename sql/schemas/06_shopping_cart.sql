-- 购物车测试数据
-- 模拟不同用户的购物车状态，包含已选中和未选中的商品

-- 清空现有数据
DELETE FROM shopping_cart;

-- 重置自增ID
ALTER TABLE shopping_cart AUTO_INCREMENT = 1;

-- 插入购物车数据
INSERT INTO shopping_cart (id, user, goods, nums, checked, created_at, updated_at, is_deleted) VALUES
-- 用户2（张三）的购物车 - 正在准备下单的用户
(1, 2, 1, 1, true, NOW(), NOW(), false),   -- iPhone 15 Pro Max，已选中
(2, 2, 8, 1, true, NOW(), NOW(), false),   -- Nike Air Jordan 1，已选中
(3, 2, 10, 2, true, NOW(), NOW(), false),  -- UNIQLO T恤 x2，已选中
(4, 2, 15, 1, false, NOW(), NOW(), false), -- 计算机系统书籍，未选中（在考虑中）

-- 用户3（李四）的购物车 - 女性用户，关注美妆和服装
(5, 3, 17, 1, true, NOW(), NOW(), false),  -- SK-II 神仙水，已选中
(6, 3, 9, 1, true, NOW(), NOW(), false),   -- Adidas 跑步鞋，已选中  
(7, 3, 11, 1, false, NOW(), NOW(), false), -- Levi's 牛仔裤，未选中
(8, 3, 20, 5, true, NOW(), NOW(), false),  -- 晨光中性笔 x5，已选中

-- 用户4（王五）的购物车 - 科技爱好者
(9, 4, 6, 1, true, NOW(), NOW(), false),   -- MacBook Pro，已选中
(10, 4, 14, 1, true, NOW(), NOW(), false), -- Canon 相机，已选中
(11, 4, 16, 1, true, NOW(), NOW(), false), -- Python编程书，已选中
(12, 4, 5, 1, false, NOW(), NOW(), false), -- Xiaomi 14 Ultra，未选中（价格考虑）

-- 用户5（赵六）的购物车 - 家庭用户
(13, 5, 12, 1, true, NOW(), NOW(), false), -- 美的空调，已选中
(14, 5, 13, 1, true, NOW(), NOW(), false), -- 海尔冰箱，已选中
(15, 5, 18, 1, true, NOW(), NOW(), false), -- IKEA 床架，已选中
(16, 5, 19, 1, false, NOW(), NOW(), false), -- 小米电视，未选中（预算考虑）

-- 用户6（孙七）的购物车 - 运动爱好者
(17, 6, 8, 1, true, NOW(), NOW(), false),  -- Nike AJ1，已选中
(18, 6, 9, 2, true, NOW(), NOW(), false),  -- Adidas 跑步鞋 x2（不同颜色），已选中
(19, 6, 10, 3, true, NOW(), NOW(), false), -- UNIQLO T恤 x3，已选中

-- 用户7（周八）的购物车 - 学生用户，预算有限
(20, 7, 2, 1, true, NOW(), NOW(), false),  -- iPhone 14（相对便宜），已选中
(21, 7, 15, 1, true, NOW(), NOW(), false), -- 计算机系统书，已选中
(22, 7, 16, 1, true, NOW(), NOW(), false), -- Python编程书，已选中
(23, 7, 20, 10, true, NOW(), NOW(), false), -- 中性笔 x10，已选中
(24, 7, 10, 1, false, NOW(), NOW(), false), -- T恤，未选中

-- 用户8（吴九）的购物车 - 空购物车用户（测试空状态）
-- 无数据

-- 用户9（郑十）的购物车 - 只浏览不购买的用户
(25, 9, 3, 1, false, NOW(), NOW(), false), -- HUAWEI Mate 60 Pro，未选中
(26, 9, 4, 1, false, NOW(), NOW(), false), -- HUAWEI P60 Pro，未选中
(27, 9, 17, 1, false, NOW(), NOW(), false), -- SK-II，未选中
(28, 9, 14, 1, false, NOW(), NOW(), false), -- Canon 相机，未选中

-- 用户10（陈一）的购物车 - 新用户，少量商品
(29, 10, 10, 1, true, NOW(), NOW(), false), -- UNIQLO T恤，已选中
(30, 10, 20, 2, true, NOW(), NOW(), false), -- 中性笔 x2，已选中

-- 用户11（李二）的购物车 - 服装关注者
(31, 11, 11, 1, true, NOW(), NOW(), false), -- Levi's 牛仔裤，已选中
(32, 11, 10, 2, true, NOW(), NOW(), false), -- UNIQLO T恤 x2，已选中

-- 用户12（王三）的购物车 - 电子产品爱好者
(33, 12, 1, 1, false, NOW(), NOW(), false), -- iPhone 15 Pro Max，未选中（价格高）
(34, 12, 2, 1, true, NOW(), NOW(), false),  -- iPhone 14，已选中
(35, 12, 7, 1, true, NOW(), NOW(), false),  -- ThinkPad X1，已选中

-- VIP用户13（VIP刘总）的购物车 - 高价值用户，多件商品
(36, 13, 1, 2, true, NOW(), NOW(), false),  -- iPhone 15 Pro Max x2，已选中
(37, 13, 6, 1, true, NOW(), NOW(), false),  -- MacBook Pro，已选中
(38, 13, 14, 1, true, NOW(), NOW(), false), -- Canon 相机，已选中
(39, 13, 12, 2, true, NOW(), NOW(), false), -- 美的空调 x2，已选中
(40, 13, 17, 3, true, NOW(), NOW(), false), -- SK-II x3，已选中

-- VIP用户14（VIP张总）的购物车 - 高价值女性用户
(41, 14, 17, 5, true, NOW(), NOW(), false), -- SK-II x5，已选中
(42, 14, 6, 1, true, NOW(), NOW(), false),  -- MacBook Pro，已选中
(43, 14, 9, 2, true, NOW(), NOW(), false),  -- Adidas 跑步鞋 x2，已选中
(44, 14, 13, 1, true, NOW(), NOW(), false); -- 海尔冰箱，已选中

-- 创建一些时间较久的购物车记录（测试购物车清理功能）
INSERT INTO shopping_cart (id, user, goods, nums, checked, created_at, updated_at, is_deleted) VALUES
-- 30天前的购物车记录
(45, 15, 1, 1, false, DATE_SUB(NOW(), INTERVAL 30 DAY), DATE_SUB(NOW(), INTERVAL 30 DAY), false),
(46, 15, 2, 1, false, DATE_SUB(NOW(), INTERVAL 25 DAY), DATE_SUB(NOW(), INTERVAL 25 DAY), false),

-- 7天前的购物车记录
(47, 15, 10, 2, true, DATE_SUB(NOW(), INTERVAL 7 DAY), DATE_SUB(NOW(), INTERVAL 7 DAY), false),
(48, 15, 20, 3, true, DATE_SUB(NOW(), INTERVAL 5 DAY), DATE_SUB(NOW(), INTERVAL 5 DAY), false);

-- 说明：
-- user: 用户ID，对应 users 表中的用户
-- goods: 商品ID，对应 goods 表中的商品  
-- nums: 商品数量，根据商品类型设置合理数量
-- checked: 是否选中，true=已选中（准备下单），false=未选中（观望中）
-- 
-- 购物车场景设计：
-- 1. 有些用户购物车全选，准备下单
-- 2. 有些用户部分选中，在考虑中
-- 3. 有些用户空购物车
-- 4. 有些用户只看不买（全部未选中）
-- 5. VIP用户购物车金额较高，数量较多
-- 6. 普通用户购物车相对简单
-- 7. 包含不同时间的购物车记录，便于测试购物车时效性