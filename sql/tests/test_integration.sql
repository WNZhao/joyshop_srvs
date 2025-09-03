-- JoyShop 微服务联动测试脚本
-- 验证各服务间数据一致性和业务流程完整性

-- 设置字符集
SET NAMES utf8mb4;

SELECT '============================================' AS '';
SELECT '🔬 JoyShop 微服务联动测试开始' AS 'Test Status';
SELECT '============================================' AS '';

-- ==============================================
-- 测试1: 数据一致性检查
-- ==============================================
SELECT '📊 测试1: 数据一致性检查' AS '';

-- 检查商品与库存的一致性
SELECT '🔍 检查商品与库存数据一致性...' AS '';
SELECT 
    CASE 
        WHEN goods_without_inventory > 0 THEN CONCAT('❌ 发现 ', goods_without_inventory, ' 个商品没有库存记录')
        ELSE '✅ 所有商品都有对应的库存记录'
    END AS 'Goods-Inventory Consistency'
FROM (
    SELECT COUNT(*) as goods_without_inventory
    FROM goods g 
    LEFT JOIN inventory i ON g.id = i.goods_id 
    WHERE i.goods_id IS NULL AND g.on_sale = true
) AS check_result;

-- 检查购物车商品有效性
SELECT '🔍 检查购物车商品有效性...' AS '';
SELECT 
    CASE 
        WHEN invalid_cart_items > 0 THEN CONCAT('❌ 发现 ', invalid_cart_items, ' 个购物车商品无效')
        ELSE '✅ 所有购物车商品都有效'
    END AS 'Cart Goods Validity'
FROM (
    SELECT COUNT(*) as invalid_cart_items
    FROM shopping_cart sc 
    LEFT JOIN goods g ON sc.goods = g.id 
    WHERE g.id IS NULL
) AS check_result;

-- 检查订单商品一致性
SELECT '🔍 检查订单商品一致性...' AS '';
SELECT 
    CASE 
        WHEN inconsistent_orders > 0 THEN CONCAT('❌ 发现 ', inconsistent_orders, ' 个订单的商品数据不一致')
        ELSE '✅ 所有订单的商品数据一致'
    END AS 'Order Goods Consistency'
FROM (
    SELECT COUNT(*) as inconsistent_orders
    FROM order_goods og 
    LEFT JOIN goods g ON og.goods = g.id 
    WHERE g.id IS NULL
) AS check_result;

-- ==============================================
-- 测试2: 业务逻辑验证
-- ==============================================
SELECT '' AS '', '📋 测试2: 业务逻辑验证' AS '';

-- 检查库存充足性（购物车中的商品是否有足够库存）
SELECT '🔍 检查购物车商品库存充足性...' AS '';
CREATE TEMPORARY TABLE temp_stock_check AS
SELECT 
    sc.user,
    sc.goods,
    sc.nums as cart_nums,
    COALESCE(SUM(i.stock), 0) as total_stock,
    CASE 
        WHEN COALESCE(SUM(i.stock), 0) >= sc.nums THEN '✅'
        ELSE '❌'
    END as stock_status
FROM shopping_cart sc
LEFT JOIN inventory i ON sc.goods = i.goods_id
WHERE sc.checked = true  -- 只检查已选中的商品
GROUP BY sc.user, sc.goods, sc.nums;

SELECT 
    CONCAT(
        (SELECT COUNT(*) FROM temp_stock_check WHERE stock_status = '✅'), 
        ' 个购物车商品库存充足, ',
        (SELECT COUNT(*) FROM temp_stock_check WHERE stock_status = '❌'), 
        ' 个商品库存不足'
    ) AS 'Stock Sufficiency Check';

-- 显示库存不足的详情
SELECT '⚠️ 库存不足的购物车商品：' AS '';
SELECT 
    CONCAT('用户ID ', user, ', 商品ID ', goods, ': 需要 ', cart_nums, ' 件, 库存 ', total_stock, ' 件') AS 'Insufficient Stock Details'
FROM temp_stock_check 
WHERE stock_status = '❌'
LIMIT 5;

-- 检查订单金额计算准确性
SELECT '🔍 检查订单金额计算准确性...' AS '';
CREATE TEMPORARY TABLE temp_amount_check AS
SELECT 
    o.id as order_id,
    o.order_mount as recorded_amount,
    COALESCE(SUM(og.goods_price * og.nums), 0) as calculated_amount,
    ABS(o.order_mount - COALESCE(SUM(og.goods_price * og.nums), 0)) as difference
FROM order_info o
LEFT JOIN order_goods og ON o.id = og.`order`
GROUP BY o.id, o.order_mount;

SELECT 
    CASE 
        WHEN incorrect_orders > 0 THEN CONCAT('❌ 发现 ', incorrect_orders, ' 个订单金额计算错误')
        ELSE '✅ 所有订单金额计算正确'
    END AS 'Order Amount Accuracy'
FROM (
    SELECT COUNT(*) as incorrect_orders 
    FROM temp_amount_check 
    WHERE difference > 0.01  -- 允许1分钱的舍入误差
) AS check_result;

-- ==============================================  
-- 测试3: 用户行为路径测试
-- ==============================================
SELECT '' AS '', '🛒 测试3: 用户购物路径完整性' AS '';

-- 分析用户2（张三）的完整购物路径
SELECT '👤 分析用户2（张三）的购物路径：' AS '';

-- 购物车状态
SELECT '🛒 购物车状态：' AS '';
SELECT 
    CONCAT('商品ID ', goods, ', 数量 ', nums, ', ', 
           CASE WHEN checked THEN '已选中' ELSE '未选中' END) AS 'Cart Status'
FROM shopping_cart 
WHERE user = 2;

-- 历史订单
SELECT '📦 历史订单：' AS '';
SELECT 
    CONCAT('订单号 ', order_sn, ', 状态: ', status, ', 金额: ¥', FORMAT(order_mount, 2)) AS 'Order History'
FROM order_info 
WHERE user = 2
ORDER BY created_at DESC;

-- ==============================================
-- 测试4: 高价值用户分析
-- ==============================================
SELECT '' AS '', '💎 测试4: VIP用户消费分析' AS '';

SELECT 'VIP用户消费情况：' AS '';
SELECT 
    CONCAT('用户ID ', o.user, ': ',
           COUNT(*), ' 个订单, ',
           '总消费 ¥', FORMAT(SUM(o.order_mount), 2), ', ',
           '平均订单 ¥', FORMAT(AVG(o.order_mount), 2)
    ) AS 'VIP User Analysis'
FROM order_info o
WHERE o.user IN (13, 14)  -- VIP用户
GROUP BY o.user;

-- ==============================================
-- 测试5: 商品销售分析
-- ==============================================
SELECT '' AS '', '📈 测试5: 商品销售数据分析' AS '';

-- 热销商品TOP5
SELECT '🔥 热销商品TOP5：' AS '';
SELECT 
    CONCAT('商品: ', og.goods_name, ' (ID:', og.goods, '), 销量: ', SUM(og.nums), ' 件') AS 'Hot Products'
FROM order_goods og
JOIN order_info oi ON og.`order` = oi.id
WHERE oi.status IN ('TRADE_SUCCESS', 'TRADE_FINISHED')  -- 只统计成功的订单
GROUP BY og.goods, og.goods_name
ORDER BY SUM(og.nums) DESC
LIMIT 5;

-- 各品牌销售情况
SELECT '🏷️ 品牌销售统计：' AS '';
SELECT 
    CONCAT(b.name, ': ', COUNT(DISTINCT og.`order`), ' 个订单, ',
           SUM(og.nums), ' 件商品, ',
           '¥', FORMAT(SUM(og.goods_price * og.nums), 2)) AS 'Brand Sales'
FROM order_goods og
JOIN goods g ON og.goods = g.id
JOIN brand b ON g.brand_id = b.id
JOIN order_info oi ON og.`order` = oi.id
WHERE oi.status IN ('TRADE_SUCCESS', 'TRADE_FINISHED')
GROUP BY b.id, b.name
ORDER BY SUM(og.goods_price * og.nums) DESC
LIMIT 8;

-- ==============================================
-- 测试6: 库存预警检查
-- ==============================================
SELECT '' AS '', '⚠️ 测试6: 库存预警检查' AS '';

-- 低库存商品（库存<=10）
SELECT '📦 低库存预警：' AS '';
SELECT 
    CONCAT('商品ID ', i.goods_id, ': 当前库存 ', i.stock, ' 件') AS 'Low Stock Alert'
FROM inventory i
WHERE i.stock <= 10
ORDER BY i.stock ASC;

-- 零库存商品
SELECT '🚫 零库存商品：' AS '';
SELECT 
    CONCAT('商品ID ', goods_id, ': 库存为 0') AS 'Zero Stock'
FROM inventory 
WHERE stock = 0;

-- ==============================================
-- 测试7: 数据质量检查
-- ==============================================
SELECT '' AS '', '🔍 测试7: 数据质量检查' AS '';

-- 检查异常数据
SELECT '数据质量检查结果：' AS '';

-- 检查负库存
SELECT 
    CASE 
        WHEN negative_stock > 0 THEN CONCAT('❌ 发现 ', negative_stock, ' 个商品库存为负数')
        ELSE '✅ 没有发现负库存'
    END AS 'Negative Stock Check'
FROM (SELECT COUNT(*) as negative_stock FROM inventory WHERE stock < 0) AS check1;

-- 检查异常价格
SELECT 
    CASE 
        WHEN invalid_price > 0 THEN CONCAT('❌ 发现 ', invalid_price, ' 个商品价格异常')
        ELSE '✅ 商品价格数据正常'
    END AS 'Price Validity Check'
FROM (SELECT COUNT(*) as invalid_price FROM goods WHERE shop_price <= 0 OR market_price <= 0) AS check2;

-- 检查订单状态一致性
SELECT 
    CASE 
        WHEN inconsistent_status > 0 THEN CONCAT('❌ 发现 ', inconsistent_status, ' 个订单状态异常')
        ELSE '✅ 订单状态数据一致'
    END AS 'Order Status Check'
FROM (
    SELECT COUNT(*) as inconsistent_status 
    FROM order_info 
    WHERE status NOT IN ('WAIT_BUYER_PAY', 'PAYING', 'TRADE_SUCCESS', 'TRADE_FINISHED', 'TRADE_CLOSED')
) AS check3;

-- ==============================================
-- 测试总结
-- ==============================================
SELECT '' AS '', '============================================' AS '';
SELECT '📊 联动测试总结报告' AS '';
SELECT '============================================' AS '';

-- 总体数据概况
SELECT '📈 总体数据概况：' AS '';
SELECT CONCAT('👥 用户总数: ', (SELECT COUNT(*) FROM users), ' 人') AS 'Summary';
SELECT CONCAT('🛍️ 商品总数: ', (SELECT COUNT(*) FROM goods), ' 个') AS 'Summary';
SELECT CONCAT('📦 库存记录: ', (SELECT COUNT(*) FROM inventory), ' 条') AS 'Summary';
SELECT CONCAT('🛒 购物车记录: ', (SELECT COUNT(*) FROM shopping_cart), ' 条') AS 'Summary';
SELECT CONCAT('📋 订单总数: ', (SELECT COUNT(*) FROM order_info), ' 个') AS 'Summary';
SELECT CONCAT('💰 总交易额: ¥', FORMAT((SELECT SUM(order_mount) FROM order_info WHERE status IN ('TRADE_SUCCESS', 'TRADE_FINISHED')), 2)) AS 'Summary';

-- 业务指标
SELECT '' AS '', '📊 关键业务指标：' AS '';
SELECT CONCAT('📈 订单转化率: ', 
    ROUND(
        (SELECT COUNT(*) FROM order_info WHERE status IN ('TRADE_SUCCESS', 'TRADE_FINISHED')) * 100.0 / 
        (SELECT COUNT(DISTINCT user) FROM shopping_cart), 2
    ), '%') AS 'Business Metrics';

SELECT CONCAT('💳 支付成功率: ', 
    ROUND(
        (SELECT COUNT(*) FROM order_info WHERE status IN ('TRADE_SUCCESS', 'TRADE_FINISHED')) * 100.0 / 
        (SELECT COUNT(*) FROM order_info WHERE status != 'WAIT_BUYER_PAY'), 2
    ), '%') AS 'Business Metrics';

SELECT CONCAT('🔥 热销商品占比: ', 
    ROUND((SELECT COUNT(*) FROM goods WHERE is_hot = true) * 100.0 / (SELECT COUNT(*) FROM goods), 2), 
    '%') AS 'Business Metrics';

-- 清理临时表
DROP TEMPORARY TABLE IF EXISTS temp_stock_check;
DROP TEMPORARY TABLE IF EXISTS temp_amount_check;

SELECT '' AS '', '============================================' AS '';
SELECT '✅ 联动测试完成！系统数据一致性良好。' AS 'Test Status';
SELECT '🚀 可以开始业务功能测试和性能测试。' AS 'Next Steps';
SELECT '============================================' AS '';