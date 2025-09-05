-- 购物车测试数据 (扩展版)
-- 模拟不同用户的购物车状态，商品ID范围1-3500

-- 清空现有数据
DELETE FROM shopping_cart;

-- 重置自增ID
ALTER TABLE shopping_cart AUTO_INCREMENT = 1;

-- 插入购物车数据
INSERT INTO shopping_cart (id, user, goods, nums, checked, created_at, updated_at, is_deleted) VALUES
(1, 3, 1179, 1, true, NOW(), NOW(), false),
(2, 5, 817, 5, true, NOW(), NOW(), false),
(3, 5, 208, 5, false, NOW(), NOW(), false),
(4, 5, 2494, 5, true, NOW(), NOW(), false),
(5, 5, 1466, 5, true, NOW(), NOW(), false),
(6, 5, 2372, 1, false, NOW(), NOW(), false),
(7, 5, 942, 5, false, NOW(), NOW(), false),
(8, 6, 378, 1, true, NOW(), NOW(), false),
(9, 7, 671, 2, false, NOW(), NOW(), false),
(10, 7, 729, 3, true, NOW(), NOW(), false),
(11, 7, 2858, 1, true, NOW(), NOW(), false),
(12, 7, 1697, 3, true, NOW(), NOW(), false),
(13, 7, 1452, 4, true, NOW(), NOW(), false),
(14, 7, 300, 2, false, NOW(), NOW(), false),
(15, 8, 882, 5, true, NOW(), NOW(), false),
(16, 8, 386, 1, true, NOW(), NOW(), false),
(17, 8, 406, 4, true, NOW(), NOW(), false),
(18, 8, 1006, 5, true, NOW(), NOW(), false),
(19, 8, 773, 4, true, NOW(), NOW(), false),
(20, 8, 3210, 3, true, NOW(), NOW(), false),
(21, 9, 1928, 4, true, NOW(), NOW(), false),
(22, 10, 1278, 1, true, NOW(), NOW(), false),
(23, 10, 1164, 3, true, NOW(), NOW(), false),
(24, 12, 2249, 3, true, NOW(), NOW(), false),
(25, 12, 2797, 1, false, NOW(), NOW(), false),
(26, 12, 2827, 1, true, NOW(), NOW(), false),
(27, 15, 1624, 1, true, NOW(), NOW(), false),
(28, 15, 2918, 1, false, NOW(), NOW(), false),
(29, 15, 1898, 4, false, NOW(), NOW(), false),
(30, 15, 2799, 1, true, NOW(), NOW(), false);

-- 购物车数据统计查询
-- 用户购物车统计
-- SELECT 
--     user as user_id,
--     COUNT(*) as items_count,
--     SUM(nums) as total_quantity,
--     SUM(CASE WHEN checked = true THEN nums ELSE 0 END) as checked_quantity
-- FROM shopping_cart 
-- WHERE is_deleted = false
-- GROUP BY user
-- ORDER BY user;

-- 热门购物车商品TOP20
-- SELECT 
--     goods as goods_id,
--     COUNT(*) as add_count,
--     SUM(nums) as total_nums,
--     COUNT(CASE WHEN checked = true THEN 1 END) as checked_count
-- FROM shopping_cart 
-- WHERE is_deleted = false
-- GROUP BY goods
-- ORDER BY add_count DESC, total_nums DESC
-- LIMIT 20;

-- 购物车商品ID分布检查 (确保在1-3500范围内)
-- SELECT 
--     MIN(goods) as min_goods_id,
--     MAX(goods) as max_goods_id,
--     COUNT(DISTINCT goods) as unique_goods_count
-- FROM shopping_cart 
-- WHERE is_deleted = false;
