-- JoyShop 微服务架构多数据库初始化脚本
-- 用于微服务架构中每个服务使用独立数据库的场景

-- ==============================================
-- 使用说明：
-- 1. 每个微服务连接自己的独立数据库
-- 2. 数据库名称：joyshop_user、joyshop_goods、joyshop_inventory、joyshop_order
-- 3. 按服务顺序逐个初始化，确保跨服务引用的ID一致
-- ==============================================

-- 设置字符集
SET NAMES utf8mb4;

SELECT '============================================' AS '';
SELECT '🚀 开始初始化微服务多数据库架构' AS 'Status';
SELECT '============================================' AS '';

-- ==============================================
-- 1. 初始化用户服务数据库 (joyshop_user)
-- ==============================================
SELECT '👤 正在初始化用户服务数据库...' AS '';

-- 切换到用户数据库
USE joyshop_user;

-- 禁用外键检查
SET FOREIGN_KEY_CHECKS = 0;

-- 导入用户数据
SOURCE 01_users.sql;

-- 启用外键检查
SET FOREIGN_KEY_CHECKS = 1;

-- 统计用户数据
SELECT '用户服务初始化完成：' AS '';
SELECT 'users' AS table_name, COUNT(*) AS record_count FROM users;

-- ==============================================
-- 2. 初始化商品服务数据库 (joyshop_goods)
-- ==============================================
SELECT '' AS '', '🛍️ 正在初始化商品服务数据库...' AS '';

-- 切换到商品数据库
USE joyshop_goods;

-- 禁用外键检查
SET FOREIGN_KEY_CHECKS = 0;

-- 导入商品相关数据（分类、品牌、商品、轮播图）
SOURCE 02_categories.sql;
SOURCE 03_brands.sql;
SOURCE 04_goods.sql;
SOURCE 08_banners.sql;

-- 启用外键检查
SET FOREIGN_KEY_CHECKS = 1;

-- 统计商品数据
SELECT '商品服务初始化完成：' AS '';
SELECT 'category' AS table_name, COUNT(*) AS record_count FROM category
UNION ALL
SELECT 'brand', COUNT(*) FROM brand
UNION ALL
SELECT 'goods', COUNT(*) FROM goods
UNION ALL
SELECT 'banner', COUNT(*) FROM banner;

-- ==============================================
-- 3. 初始化库存服务数据库 (joyshop_inventory)
-- ==============================================
SELECT '' AS '', '📦 正在初始化库存服务数据库...' AS '';

-- 切换到库存数据库
USE joyshop_inventory;

-- 禁用外键检查
SET FOREIGN_KEY_CHECKS = 0;

-- 导入库存数据
SOURCE 05_inventory.sql;

-- 启用外键检查
SET FOREIGN_KEY_CHECKS = 1;

-- 统计库存数据
SELECT '库存服务初始化完成：' AS '';
SELECT 'inventory' AS table_name, COUNT(*) AS record_count FROM inventory;

-- ==============================================
-- 4. 初始化订单服务数据库 (joyshop_order)
-- ==============================================
SELECT '' AS '', '📋 正在初始化订单服务数据库...' AS '';

-- 切换到订单数据库
USE joyshop_order;

-- 禁用外键检查
SET FOREIGN_KEY_CHECKS = 0;

-- 导入订单相关数据（购物车、订单）
SOURCE 06_shopping_cart.sql;
SOURCE 07_orders.sql;

-- 启用外键检查
SET FOREIGN_KEY_CHECKS = 1;

-- 统计订单数据
SELECT '订单服务初始化完成：' AS '';
SELECT 'shopping_cart' AS table_name, COUNT(*) AS record_count FROM shopping_cart
UNION ALL
SELECT 'order_info', COUNT(*) FROM order_info
UNION ALL
SELECT 'order_goods', COUNT(*) FROM order_goods;

-- ==============================================
-- 初始化完成总结
-- ==============================================
SELECT '' AS '', '============================================' AS '';
SELECT '✅ 微服务多数据库初始化完成！' AS 'Status';
SELECT '============================================' AS '';

-- 数据库使用说明
SELECT '📊 数据库分布说明：' AS '';
SELECT 'joyshop_user: 用户相关数据' AS 'Database Distribution';
SELECT 'joyshop_goods: 商品、分类、品牌、轮播图数据' AS 'Database Distribution';
SELECT 'joyshop_inventory: 库存数据' AS 'Database Distribution';
SELECT 'joyshop_order: 购物车、订单数据' AS 'Database Distribution';

SELECT '' AS '', '🔧 服务配置建议：' AS '';
SELECT 'user_srv -> joyshop_user' AS 'Service Config';
SELECT 'goods_srv -> joyshop_goods' AS 'Service Config';
SELECT 'inventory_srv -> joyshop_inventory' AS 'Service Config';
SELECT 'order_srv -> joyshop_order' AS 'Service Config';

SELECT '' AS '', '⚠️ 重要提醒：' AS '';
SELECT '1. 微服务间通过API调用，不直接访问其他服务数据库' AS 'Important Notes';
SELECT '2. 用户ID、商品ID等需要保证跨服务一致性' AS 'Important Notes';
SELECT '3. 建议使用分布式事务管理跨服务数据一致性' AS 'Important Notes';
SELECT '4. 测试时注意验证跨服务数据关联的正确性' AS 'Important Notes';

SELECT '' AS '', '🧪 测试验证建议：' AS '';
SELECT '1. 验证用户服务：登录、查询用户信息' AS 'Test Suggestions';
SELECT '2. 验证商品服务：查询商品、分类、品牌' AS 'Test Suggestions';
SELECT '3. 验证库存服务：查询库存、扣减库存' AS 'Test Suggestions';
SELECT '4. 验证订单服务：购物车操作、创建订单' AS 'Test Suggestions';
SELECT '5. 验证跨服务调用：下单时的库存检查和扣减' AS 'Test Suggestions';

SELECT '============================================' AS '';
SELECT '🎉 可以开始微服务联动测试了！' AS 'Status';
SELECT '============================================' AS '';