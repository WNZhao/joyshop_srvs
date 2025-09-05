#!/usr/bin/env python3
"""
数据完整性验证脚本
检查所有表之间的关联关系是否正确，特别是扩展到3500商品后的数据一致性
"""

def generate_validation_sql():
    """生成数据完整性验证SQL"""
    sql_content = """-- 数据完整性验证查询
-- 检查商品数据扩展后的各表关联关系

-- ========== 基础数据量统计 ==========
SELECT '基础数据量统计' as check_type;

-- 各表记录数量
SELECT 'users' as table_name, COUNT(*) as record_count FROM users
UNION ALL
SELECT 'brands' as table_name, COUNT(*) as record_count FROM brand  
UNION ALL
SELECT 'categories' as table_name, COUNT(*) as record_count FROM category
UNION ALL
SELECT 'goods' as table_name, COUNT(*) as record_count FROM goods
UNION ALL
SELECT 'inventory' as table_name, COUNT(*) as record_count FROM inventory
UNION ALL
SELECT 'goods_category' as table_name, COUNT(*) as record_count FROM goods_category
UNION ALL
SELECT 'order_info' as table_name, COUNT(*) as record_count FROM order_info
UNION ALL
SELECT 'order_goods' as table_name, COUNT(*) as record_count FROM order_goods
UNION ALL
SELECT 'shopping_cart' as table_name, COUNT(*) as record_count FROM shopping_cart
UNION ALL
SELECT 'banner' as table_name, COUNT(*) as record_count FROM banner
ORDER BY record_count DESC;

-- ========== 商品数据完整性检查 ==========
SELECT '商品数据完整性检查' as check_type;

-- 检查商品ID范围
SELECT 
    'goods_id_range' as check_item,
    MIN(id) as min_id,
    MAX(id) as max_id,
    COUNT(*) as total_count,
    CASE 
        WHEN COUNT(*) >= 3000 AND COUNT(*) <= 5000 THEN '✓ 商品数量符合要求'
        ELSE '✗ 商品数量不符合要求'
    END as status
FROM goods;

-- ========== 库存数据完整性检查 ==========  
SELECT '库存数据完整性检查' as check_type;

-- 检查库存覆盖率
SELECT 
    'inventory_coverage' as check_item,
    (SELECT COUNT(*) FROM goods) as goods_count,
    (SELECT COUNT(*) FROM inventory) as inventory_count,
    CASE 
        WHEN (SELECT COUNT(*) FROM goods) = (SELECT COUNT(*) FROM inventory) 
        THEN '✓ 库存覆盖率100%'
        ELSE CONCAT('✗ 库存缺失 ', (SELECT COUNT(*) FROM goods) - (SELECT COUNT(*) FROM inventory), ' 个商品')
    END as status;

-- 检查孤立的库存记录(商品不存在)
SELECT 
    'orphaned_inventory' as check_item,
    COUNT(*) as orphaned_count,
    CASE 
        WHEN COUNT(*) = 0 THEN '✓ 无孤立库存记录'
        ELSE CONCAT('✗ 发现 ', COUNT(*), ' 个孤立库存记录')
    END as status
FROM inventory i 
LEFT JOIN goods g ON i.goods_id = g.id 
WHERE g.id IS NULL;

-- ========== 商品分类关联检查 ==========
SELECT '商品分类关联检查' as check_type;

-- 检查商品分类关联覆盖率
SELECT 
    'category_coverage' as check_item,
    (SELECT COUNT(*) FROM goods) as goods_count,
    (SELECT COUNT(DISTINCT goods_id) FROM goods_category) as categorized_goods,
    CASE 
        WHEN (SELECT COUNT(*) FROM goods) = (SELECT COUNT(DISTINCT goods_id) FROM goods_category)
        THEN '✓ 所有商品都有分类'
        ELSE CONCAT('✗ ', (SELECT COUNT(*) FROM goods) - (SELECT COUNT(DISTINCT goods_id) FROM goods_category), ' 个商品无分类')
    END as status;

-- 检查孤立的商品分类关联
SELECT 
    'orphaned_goods_category' as check_item,
    COUNT(*) as orphaned_count,
    CASE 
        WHEN COUNT(*) = 0 THEN '✓ 无孤立分类关联'
        ELSE CONCAT('✗ 发现 ', COUNT(*), ' 个孤立分类关联')
    END as status
FROM goods_category gc 
LEFT JOIN goods g ON gc.goods_id = g.id 
WHERE g.id IS NULL;

-- ========== 订单数据完整性检查 ==========
SELECT '订单数据完整性检查' as check_type;

-- 检查订单商品引用范围
SELECT 
    'order_goods_range' as check_item,
    MIN(goods_id) as min_goods_id,
    MAX(goods_id) as max_goods_id,
    COUNT(DISTINCT goods_id) as unique_goods_count,
    CASE 
        WHEN MAX(goods_id) <= 3500 AND MIN(goods_id) >= 1 
        THEN '✓ 订单商品ID范围正确'
        ELSE '✗ 订单商品ID超出范围'
    END as status
FROM order_goods;

-- 检查孤立的订单商品记录  
SELECT 
    'orphaned_order_goods' as check_item,
    COUNT(*) as orphaned_count,
    CASE 
        WHEN COUNT(*) = 0 THEN '✓ 无孤立订单商品记录'
        ELSE CONCAT('✗ 发现 ', COUNT(*), ' 个孤立订单商品记录')
    END as status
FROM order_goods og 
LEFT JOIN goods g ON og.goods_id = g.id 
WHERE g.id IS NULL;

-- 检查订单商品与订单信息关联
SELECT 
    'order_goods_info_link' as check_item,
    COUNT(*) as orphaned_count,
    CASE 
        WHEN COUNT(*) = 0 THEN '✓ 所有订单商品都有对应订单信息'
        ELSE CONCAT('✗ 发现 ', COUNT(*), ' 个订单商品缺少订单信息')
    END as status
FROM order_goods og 
LEFT JOIN order_info oi ON og.order_id = oi.id 
WHERE oi.id IS NULL;

-- ========== 购物车数据完整性检查 ==========
SELECT '购物车数据完整性检查' as check_type;

-- 检查购物车商品引用范围
SELECT 
    'cart_goods_range' as check_item,
    MIN(goods) as min_goods_id,
    MAX(goods) as max_goods_id,
    COUNT(DISTINCT goods) as unique_goods_count,
    CASE 
        WHEN MAX(goods) <= 3500 AND MIN(goods) >= 1 
        THEN '✓ 购物车商品ID范围正确'
        ELSE '✗ 购物车商品ID超出范围'
    END as status
FROM shopping_cart 
WHERE is_deleted = false;

-- 检查孤立的购物车记录
SELECT 
    'orphaned_cart_goods' as check_item,
    COUNT(*) as orphaned_count,
    CASE 
        WHEN COUNT(*) = 0 THEN '✓ 无孤立购物车记录'
        ELSE CONCAT('✗ 发现 ', COUNT(*), ' 个孤立购物车记录')
    END as status
FROM shopping_cart sc 
LEFT JOIN goods g ON sc.goods = g.id 
WHERE sc.is_deleted = false AND g.id IS NULL;

-- ========== 品牌关联检查 ==========
SELECT '品牌关联检查' as check_type;

-- 检查孤立的商品品牌关联
SELECT 
    'orphaned_goods_brand' as check_item,
    COUNT(*) as orphaned_count,
    CASE 
        WHEN COUNT(*) = 0 THEN '✓ 无孤立商品品牌关联'
        ELSE CONCAT('✗ 发现 ', COUNT(*), ' 个孤立商品品牌关联')
    END as status
FROM goods g 
LEFT JOIN brand b ON g.brand_id = b.id 
WHERE b.id IS NULL;

-- ========== 数据分布统计 ==========
SELECT '数据分布统计' as check_type;

-- 商品价格分布
SELECT 
    'goods_price_distribution' as check_item,
    CONCAT('¥', MIN(shop_price), ' - ¥', MAX(shop_price)) as price_range,
    CONCAT('¥', ROUND(AVG(shop_price), 2)) as avg_price,
    COUNT(*) as total_goods
FROM goods;

-- 库存分布统计
SELECT 
    'inventory_distribution' as check_item,
    CASE 
        WHEN stock >= 200 THEN '高库存(>=200)'
        WHEN stock >= 80 THEN '中库存(80-199)'
        ELSE '低库存(<80)'
    END as stock_level,
    COUNT(*) as goods_count,
    CONCAT(ROUND(COUNT(*) * 100.0 / (SELECT COUNT(*) FROM inventory), 1), '%') as percentage
FROM inventory 
GROUP BY 
    CASE 
        WHEN stock >= 200 THEN '高库存(>=200)'
        WHEN stock >= 80 THEN '中库存(80-199)'
        ELSE '低库存(<80)'
    END
ORDER BY 
    CASE 
        WHEN stock >= 200 THEN 1
        WHEN stock >= 80 THEN 2
        ELSE 3
    END;

-- 订单状态分布
SELECT 
    'order_status_distribution' as check_item,
    CASE status
        WHEN 1 THEN '待付款'
        WHEN 2 THEN '已付款待发货'
        WHEN 3 THEN '已发货待收货'
        WHEN 4 THEN '已完成'
        WHEN 5 THEN '已取消'
        WHEN 6 THEN '退款中'
        WHEN 7 THEN '已退款'
        ELSE CONCAT('未知状态:', status)
    END as status_name,
    COUNT(*) as order_count,
    CONCAT('¥', FORMAT(SUM(order_mount), 2)) as total_amount
FROM order_info 
GROUP BY status
ORDER BY status;

-- ========== 性能相关检查 ==========
SELECT '性能相关检查' as check_type;

-- 检查是否需要优化的大表查询
SELECT 
    'large_table_performance' as check_item,
    'goods' as table_name,
    COUNT(*) as record_count,
    CASE 
        WHEN COUNT(*) > 1000 THEN '建议添加适当索引以优化查询性能'
        ELSE '数据量适中'
    END as recommendation
FROM goods
UNION ALL
SELECT 
    'large_table_performance' as check_item,
    'order_goods' as table_name,
    COUNT(*) as record_count,
    CASE 
        WHEN COUNT(*) > 1000 THEN '建议添加适当索引以优化查询性能'
        ELSE '数据量适中'
    END as recommendation
FROM order_goods;

-- ========== 总结 ==========
SELECT '数据完整性检查总结' as summary,
       '所有检查项目请参考上述结果' as note,
       '如发现任何 ✗ 标记的问题，请及时修复' as action_required;
"""
    
    return sql_content

if __name__ == '__main__':
    print("生成数据完整性验证SQL...")
    
    # 生成验证SQL
    sql_content = generate_validation_sql()
    
    # 写入文件
    with open('/Users/walker/gitsource/github/joyshop_srvs/sql/scripts/validate_data_integrity.sql', 'w', encoding='utf-8') as f:
        f.write(sql_content)
    
    print("数据验证SQL文件已生成: scripts/validate_data_integrity.sql")
    print("\n使用方法:")
    print("1. 在MySQL中执行所有数据初始化SQL文件")
    print("2. 执行验证SQL: source scripts/validate_data_integrity.sql")
    print("3. 检查结果中是否有 ✗ 标记的问题")
    
    print("\n推荐执行顺序:")
    print("1. init/create_databases.sql")
    print("2. schemas/01_users.sql")
    print("3. schemas/02_categories.sql") 
    print("4. schemas/03_brands.sql")
    print("5. schemas/04_goods.sql (3500商品)")
    print("6. schemas/05_inventory.sql (3500库存)")
    print("7. schemas/06_shopping_cart.sql")
    print("8. schemas/07_orders.sql (800订单)")
    print("9. schemas/08_banners.sql")
    print("10. scripts/validate_data_integrity.sql (验证)")