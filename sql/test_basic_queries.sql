-- 基础功能测试查询
-- 用于验证各个服务的基本CRUD操作

-- 设置字符集
SET NAMES utf8mb4;

SELECT '============================================' AS '';
SELECT '🧪 JoyShop 基础功能测试查询' AS 'Test Title';
SELECT '============================================' AS '';

-- ==============================================
-- 用户服务测试查询
-- ==============================================
SELECT '👤 用户服务测试查询' AS '';

-- 用户登录测试
SELECT '1. 用户登录验证测试：' AS '';
SELECT 
    id, user_name, nick_name, mobile, email, role,
    CASE 
        WHEN password = 'e10adc3949ba59abbe56e057f20f883e' THEN '✅ 密码正确'
        ELSE '❌ 密码错误'
    END as password_check
FROM users 
WHERE user_name = 'zhangsan' AND password = 'e10adc3949ba59abbe56e057f20f883e';

-- 用户列表查询
SELECT '2. 用户列表查询（分页）：' AS '';
SELECT id, user_name, nick_name, mobile, role, created_at
FROM users 
WHERE role = 1  -- 普通用户
ORDER BY created_at DESC 
LIMIT 5;

-- VIP用户查询
SELECT '3. VIP用户查询：' AS '';
SELECT id, user_name, nick_name, mobile, 
       (SELECT COUNT(*) FROM order_info WHERE user = users.id) as order_count,
       (SELECT COALESCE(SUM(order_mount), 0) FROM order_info WHERE user = users.id) as total_amount
FROM users 
WHERE id IN (13, 14);

-- ==============================================
-- 商品服务测试查询  
-- ==============================================
SELECT '' AS '', '🛍️ 商品服务测试查询' AS '';

-- 分类层级查询
SELECT '1. 分类层级结构：' AS '';
SELECT 
    CASE level
        WHEN 1 THEN CONCAT('📁 ', name)
        WHEN 2 THEN CONCAT('  ├─ ', name)
        WHEN 3 THEN CONCAT('    └─ ', name)
    END as category_tree
FROM category 
WHERE level <= 3
ORDER BY parent_id, sort, id;

-- 热销商品查询
SELECT '2. 热销商品查询：' AS '';
SELECT g.id, g.name, g.shop_price, g.market_price, g.click_num, g.fav_num, b.name as brand_name
FROM goods g
JOIN brand b ON g.brand_id = b.id
WHERE g.is_hot = true AND g.on_sale = true
ORDER BY g.click_num DESC, g.fav_num DESC
LIMIT 5;

-- 分类商品统计
SELECT '3. 分类商品统计：' AS '';
SELECT c.name as category_name, COUNT(gc.goods_id) as goods_count
FROM category c
LEFT JOIN goods_category gc ON c.id = gc.category_id
WHERE c.level = 1
GROUP BY c.id, c.name
ORDER BY goods_count DESC;

-- 品牌商品统计
SELECT '4. 品牌商品统计：' AS '';
SELECT b.name as brand_name, COUNT(g.id) as goods_count,
       AVG(g.shop_price) as avg_price
FROM brand b
LEFT JOIN goods g ON b.id = g.brand_id
GROUP BY b.id, b.name
HAVING goods_count > 0
ORDER BY goods_count DESC, avg_price DESC;

-- ==============================================
-- 库存服务测试查询
-- ==============================================
SELECT '' AS '', '📦 库存服务测试查询' AS '';

-- 库存状态概览
SELECT '1. 库存状态概览：' AS '';
SELECT 
    CASE 
        WHEN stock = 0 THEN '🚫 零库存'
        WHEN stock <= 10 THEN '⚠️ 低库存'
        WHEN stock <= 50 THEN '📦 正常库存'
        WHEN stock <= 200 THEN '📈 充足库存'
        ELSE '🏭 大量库存'
    END as stock_level,
    COUNT(*) as count,
    CONCAT(ROUND(COUNT(*) * 100.0 / (SELECT COUNT(*) FROM inventory), 1), '%') as percentage
FROM inventory
GROUP BY 
    CASE 
        WHEN stock = 0 THEN '🚫 零库存'
        WHEN stock <= 10 THEN '⚠️ 低库存'
        WHEN stock <= 50 THEN '📦 正常库存'
        WHEN stock <= 200 THEN '📈 充足库存'
        ELSE '🏭 大量库存'
    END
ORDER BY count DESC;

-- 特定商品库存查询
SELECT '2. 高价值商品库存查询：' AS '';
SELECT i.goods_id, g.name as goods_name, g.shop_price, i.stock, i.version,
       CASE 
           WHEN i.stock > 50 THEN '✅ 库存充足'
           WHEN i.stock > 10 THEN '⚠️ 库存偏少' 
           ELSE '❌ 库存不足'
       END as stock_status
FROM inventory i
JOIN goods g ON i.goods_id = g.id
WHERE g.shop_price > 5000  -- 高价值商品
ORDER BY g.shop_price DESC;

-- 库存周转分析（基于订单数据）
SELECT '3. 商品销售与库存分析：' AS '';
SELECT g.id, g.name, 
       COALESCE(SUM(i.stock), 0) as total_stock,
       COALESCE(sales.sold_qty, 0) as sold_qty,
       CASE 
           WHEN COALESCE(SUM(i.stock), 0) > 0 AND COALESCE(sales.sold_qty, 0) > 0 
           THEN ROUND(COALESCE(sales.sold_qty, 0) * 1.0 / SUM(i.stock) * 100, 2)
           ELSE 0
       END as turnover_rate
FROM goods g
LEFT JOIN inventory i ON g.id = i.goods_id
LEFT JOIN (
    SELECT og.goods, SUM(og.nums) as sold_qty
    FROM order_goods og
    JOIN order_info oi ON og.`order` = oi.id
    WHERE oi.status IN ('TRADE_SUCCESS', 'TRADE_FINISHED')
    GROUP BY og.goods
) sales ON g.id = sales.goods
GROUP BY g.id, g.name, sales.sold_qty
ORDER BY turnover_rate DESC
LIMIT 10;

-- ==============================================
-- 订单服务测试查询
-- ==============================================
SELECT '' AS '', '📋 订单服务测试查询' AS '';

-- 购物车状态分析
SELECT '1. 购物车状态分析：' AS '';
SELECT 
    COUNT(DISTINCT user) as total_users_with_cart,
    COUNT(*) as total_cart_items,
    SUM(CASE WHEN checked THEN 1 ELSE 0 END) as checked_items,
    SUM(CASE WHEN checked THEN 0 ELSE 1 END) as unchecked_items,
    SUM(nums) as total_quantity
FROM shopping_cart;

-- 用户购物车详情（示例用户）
SELECT '2. 用户购物车详情（用户ID=2）：' AS '';
SELECT sc.goods, g.name as goods_name, g.shop_price, sc.nums,
       ROUND(g.shop_price * sc.nums, 2) as subtotal,
       CASE WHEN sc.checked THEN '✅ 已选' ELSE '⏸️ 未选' END as status
FROM shopping_cart sc
JOIN goods g ON sc.goods = g.id
WHERE sc.user = 2
ORDER BY sc.checked DESC, sc.created_at DESC;

-- 订单状态统计
SELECT '3. 订单状态统计：' AS '';
SELECT 
    status,
    COUNT(*) as order_count,
    CONCAT('¥', FORMAT(SUM(order_mount), 2)) as total_amount,
    CONCAT('¥', FORMAT(AVG(order_mount), 2)) as avg_amount
FROM order_info
GROUP BY status
ORDER BY COUNT(*) DESC;

-- 订单详情查询（示例订单）
SELECT '4. 订单详情查询（订单ID=1）：' AS '';
SELECT 
    '订单基本信息：' as info_type,
    CONCAT('订单号: ', order_sn, ', 用户ID: ', user, ', 状态: ', status, ', 金额: ¥', order_mount) as details
FROM order_info WHERE id = 1
UNION ALL
SELECT 
    '收货信息：' as info_type,
    CONCAT('收货人: ', signer_name, ', 电话: ', singer_mobile, ', 地址: ', address) as details
FROM order_info WHERE id = 1
UNION ALL
SELECT 
    '商品信息：' as info_type,
    CONCAT(goods_name, ' × ', nums, ', 单价: ¥', goods_price, ', 小计: ¥', goods_price * nums) as details
FROM order_goods WHERE `order` = 1;

-- 用户消费排行
SELECT '5. 用户消费排行TOP5：' AS '';
SELECT 
    u.id, u.nick_name, 
    COUNT(o.id) as order_count,
    CONCAT('¥', FORMAT(SUM(o.order_mount), 2)) as total_spent,
    CONCAT('¥', FORMAT(AVG(o.order_mount), 2)) as avg_order
FROM users u
JOIN order_info o ON u.id = o.user
WHERE o.status IN ('TRADE_SUCCESS', 'TRADE_FINISHED')
GROUP BY u.id, u.nick_name
ORDER BY SUM(o.order_mount) DESC
LIMIT 5;

-- ==============================================
-- 跨服务联合查询测试
-- ==============================================
SELECT '' AS '', '🔗 跨服务联合查询测试' AS '';

-- 完整的商品信息（商品+库存+销量）
SELECT '1. 商品完整信息查询：' AS '';
SELECT g.id, g.name, b.name as brand, g.shop_price,
       COALESCE(SUM(i.stock), 0) as total_stock,
       COALESCE(sales.total_sold, 0) as total_sold,
       COALESCE(cart.in_cart, 0) as in_cart_qty
FROM goods g
LEFT JOIN brand b ON g.brand_id = b.id
LEFT JOIN inventory i ON g.id = i.goods_id
LEFT JOIN (
    SELECT og.goods, SUM(og.nums) as total_sold
    FROM order_goods og
    JOIN order_info oi ON og.`order` = oi.id
    WHERE oi.status IN ('TRADE_SUCCESS', 'TRADE_FINISHED')
    GROUP BY og.goods
) sales ON g.id = sales.goods
LEFT JOIN (
    SELECT goods, SUM(nums) as in_cart
    FROM shopping_cart 
    WHERE checked = true
    GROUP BY goods
) cart ON g.id = cart.goods
WHERE g.on_sale = true
GROUP BY g.id, g.name, b.name, g.shop_price, sales.total_sold, cart.in_cart
ORDER BY g.shop_price DESC
LIMIT 8;

-- 用户完整画像
SELECT '2. 用户完整画像（用户ID=2）：' AS '';
SELECT 
    '基本信息' as info_category,
    CONCAT('用户名: ', user_name, ', 昵称: ', nick_name, ', 手机: ', mobile) as info_detail
FROM users WHERE id = 2
UNION ALL
SELECT 
    '购物车状态' as info_category,
    CONCAT('商品种类: ', COUNT(DISTINCT goods), ', 总数量: ', SUM(nums), ', 已选中: ', SUM(CASE WHEN checked THEN nums ELSE 0 END)) as info_detail
FROM shopping_cart WHERE user = 2
UNION ALL
SELECT 
    '消费记录' as info_category,
    CONCAT('订单数量: ', COUNT(*), ', 总消费: ¥', FORMAT(SUM(order_mount), 2), ', 平均订单: ¥', FORMAT(AVG(order_mount), 2)) as info_detail
FROM order_info WHERE user = 2 AND status IN ('TRADE_SUCCESS', 'TRADE_FINISHED');

SELECT '' AS '', '============================================' AS '';
SELECT '✅ 基础功能测试查询完成！' AS 'Status';
SELECT '💡 以上查询展示了各服务的核心功能和数据关联' AS 'Summary';
SELECT '🔧 可以基于这些查询开发相应的API接口' AS 'Next Steps';
SELECT '============================================' AS '';