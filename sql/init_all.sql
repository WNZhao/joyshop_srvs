-- JoyShop 微服务联动测试数据初始化脚本（多数据库版本）
-- 按依赖顺序执行所有初始化SQL，确保数据一致性

-- ==============================================
-- 说明：本脚本用于初始化所有微服务的测试数据
-- 执行顺序很重要，必须按照依赖关系执行
-- 适用于微服务架构的多数据库部署
-- ==============================================

-- ⚠️ 使用说明：
-- 1. 微服务模式：每个服务使用独立数据库
-- 2. 单体模式：所有数据在同一个数据库（开发测试用）
-- 3. 请根据实际部署模式选择对应的初始化脚本

-- 禁用外键检查，加快导入速度
SET FOREIGN_KEY_CHECKS = 0;

-- 设置字符集
SET NAMES utf8mb4;

-- 1. 初始化用户数据（基础数据，其他服务会引用用户ID）
SOURCE 01_users.sql;

-- 2. 初始化商品分类数据（商品依赖分类）
SOURCE 02_categories.sql;

-- 3. 初始化品牌数据（商品依赖品牌）
SOURCE 03_brands.sql;

-- 4. 初始化商品数据（库存、购物车、订单依赖商品）
SOURCE 04_goods.sql;

-- 5. 初始化库存数据（订单创建时会检查库存）
SOURCE 05_inventory.sql;

-- 6. 初始化购物车数据（订单创建会清理购物车）
SOURCE 06_shopping_cart.sql;

-- 7. 初始化订单数据（依赖用户、商品等数据）
SOURCE 07_orders.sql;

-- 8. 初始化轮播图数据（独立数据）
SOURCE 08_banners.sql;

-- 启用外键检查
SET FOREIGN_KEY_CHECKS = 1;

-- 显示初始化完成信息
SELECT '===========================================' AS '';
SELECT '✅ JoyShop 测试数据初始化完成！' AS 'Status';
SELECT '===========================================' AS '';

-- 显示数据统计
SELECT 'users' AS table_name, COUNT(*) AS record_count FROM users
UNION ALL
SELECT 'category', COUNT(*) FROM category
UNION ALL  
SELECT 'brand', COUNT(*) FROM brand
UNION ALL
SELECT 'goods', COUNT(*) FROM goods
UNION ALL
SELECT 'inventory', COUNT(*) FROM inventory  
UNION ALL
SELECT 'shopping_cart', COUNT(*) FROM shopping_cart
UNION ALL
SELECT 'order_info', COUNT(*) FROM order_info
UNION ALL
SELECT 'order_goods', COUNT(*) FROM order_goods
UNION ALL
SELECT 'banner', COUNT(*) FROM banner;

SELECT '===========================================' AS '';
SELECT '📊 数据统计完成，可以开始联动测试！' AS 'Status';
SELECT '===========================================' AS '';

-- 推荐的测试场景：
SELECT '🧪 推荐测试场景：' AS '';
SELECT '1. 用户登录：使用用户ID 2-15，密码 123456' AS 'Test Case';
SELECT '2. 商品浏览：查看不同分类下的商品' AS 'Test Case';
SELECT '3. 购物车操作：查看用户2-14的购物车' AS 'Test Case';
SELECT '4. 订单创建：使用已选中商品的购物车创建订单' AS 'Test Case';
SELECT '5. 库存检查：验证订单创建后的库存变化' AS 'Test Case';
SELECT '6. 订单状态流转：测试支付、发货、完成等状态' AS 'Test Case';
SELECT '===========================================' AS '';