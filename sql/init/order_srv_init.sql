-- 订单服务初始化脚本
-- 初始化订单服务相关数据：购物车、订单信息、订单商品

-- 设置字符集
SET NAMES utf8mb4;

-- 禁用外键检查
SET FOREIGN_KEY_CHECKS = 0;

-- 按依赖顺序执行
-- 1. 购物车数据
SOURCE 06_shopping_cart.sql;

-- 2. 订单相关数据
SOURCE 07_orders.sql;

-- 启用外键检查
SET FOREIGN_KEY_CHECKS = 1;

-- 显示结果
SELECT '===========================================' AS '';
SELECT '✅ 订单服务测试数据初始化完成！' AS 'Status';
SELECT '===========================================' AS '';

-- 统计信息
SELECT 'shopping_cart' AS table_name, COUNT(*) AS record_count FROM shopping_cart
UNION ALL
SELECT 'order_info', COUNT(*) FROM order_info
UNION ALL
SELECT 'order_goods', COUNT(*) FROM order_goods;

-- 购物车分析
SELECT '🛒 购物车分析：' AS '';
SELECT CONCAT('购物车总记录：', COUNT(*), ' 条') AS 'Cart Analysis' FROM shopping_cart;
SELECT CONCAT('已选中商品：', COUNT(*), ' 条') AS 'Cart Analysis' FROM shopping_cart WHERE checked = true;
SELECT CONCAT('未选中商品：', COUNT(*), ' 条') AS 'Cart Analysis' FROM shopping_cart WHERE checked = false;

-- 按用户统计购物车
SELECT '👤 用户购物车统计：' AS '';
SELECT 
    CONCAT('用户ID ', user, ': ', COUNT(*), ' 件商品 (', 
           SUM(CASE WHEN checked THEN 1 ELSE 0 END), ' 已选中)') AS 'User Cart Summary'
FROM shopping_cart 
GROUP BY user 
ORDER BY user;

-- 订单状态分析
SELECT '📋 订单状态分析：' AS '';
SELECT status AS 'Order Status', COUNT(*) AS 'Count', 
       CONCAT(ROUND(COUNT(*) * 100.0 / (SELECT COUNT(*) FROM order_info), 2), '%') AS 'Percentage'
FROM order_info 
GROUP BY status 
ORDER BY COUNT(*) DESC;

-- 订单金额分析
SELECT '💰 订单金额分析：' AS '';
SELECT CONCAT('订单总数：', COUNT(*), ' 个') AS 'Order Amount Analysis' FROM order_info;
SELECT CONCAT('总交易额：¥', FORMAT(SUM(order_mount), 2)) AS 'Order Amount Analysis' FROM order_info;
SELECT CONCAT('平均订单金额：¥', FORMAT(AVG(order_mount), 2)) AS 'Order Amount Analysis' FROM order_info;
SELECT CONCAT('最高订单金额：¥', FORMAT(MAX(order_mount), 2)) AS 'Order Amount Analysis' FROM order_info;
SELECT CONCAT('最低订单金额：¥', FORMAT(MIN(order_mount), 2)) AS 'Order Amount Analysis' FROM order_info;

-- 支付方式统计
SELECT '💳 支付方式统计：' AS '';
SELECT 
    pay_type AS 'Payment Method',
    COUNT(*) AS 'Count',
    CONCAT('¥', FORMAT(SUM(order_mount), 2)) AS 'Total Amount'
FROM order_info 
GROUP BY pay_type;

-- 用户订单统计
SELECT '🛍️ 用户订单统计：' AS '';
SELECT 
    CONCAT('用户ID ', user, ': ', COUNT(*), ' 个订单, 总额¥', FORMAT(SUM(order_mount), 2)) AS 'User Order Summary'
FROM order_info 
GROUP BY user 
ORDER BY SUM(order_mount) DESC;

-- 热门商品统计（基于订单）
SELECT '🔥 热门商品统计（订单中）：' AS '';
SELECT 
    CONCAT('商品ID ', goods, ': ', goods_name, ' - 销量 ', SUM(nums), ' 件') AS 'Hot Products'
FROM order_goods 
GROUP BY goods, goods_name 
ORDER BY SUM(nums) DESC 
LIMIT 10;

SELECT '===========================================' AS '';
SELECT '🧪 测试建议：' AS '';
SELECT '1. 购物车操作：增删改查用户购物车' AS 'Test Suggestions';
SELECT '2. 订单创建：使用购物车创建订单' AS 'Test Suggestions';
SELECT '3. 订单查询：按用户、状态查询订单' AS 'Test Suggestions';
SELECT '4. 订单状态流转：支付、发货、完成流程' AS 'Test Suggestions';
SELECT '5. 批量处理测试：一次处理多个购物车商品' AS 'Test Suggestions';
SELECT '6. 分布式锁测试：并发创建订单防重复' AS 'Test Suggestions';
SELECT '===========================================' AS '';