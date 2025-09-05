#!/usr/bin/env python3
"""
购物车数据生成脚本
为用户生成购物车数据，商品ID引用1-3500范围
"""

import random

def generate_cart_data(max_goods_id=3500):
    """生成购物车数据"""
    cart_items = []
    cart_id = 1
    
    # 为15个用户生成购物车数据
    for user_id in range(1, 16):
        # 随机决定购物车中商品数量 (0-8个商品)
        # 有些用户购物车可能为空
        if random.random() < 0.2:  # 20%概率购物车为空
            continue
            
        items_count = random.randint(1, 8)
        
        # 随机选择商品ID，避免重复
        selected_goods = random.sample(range(1, max_goods_id + 1), items_count)
        
        for goods_id in selected_goods:
            # 随机数量 (1-5个)
            nums = random.randint(1, 5)
            
            # 70%概率选中商品
            checked = random.random() < 0.7
            
            cart_item = {
                'id': cart_id,
                'user': user_id,
                'goods': goods_id,
                'nums': nums,
                'checked': checked
            }
            
            cart_items.append(cart_item)
            cart_id += 1
    
    return cart_items

def generate_cart_sql(cart_items):
    """生成购物车SQL文件内容"""
    sql_content = """-- 购物车测试数据 (扩展版)
-- 模拟不同用户的购物车状态，商品ID范围1-3500

-- 清空现有数据
DELETE FROM shopping_cart;

-- 重置自增ID
ALTER TABLE shopping_cart AUTO_INCREMENT = 1;

-- 插入购物车数据
INSERT INTO shopping_cart (id, user, goods, nums, checked, created_at, updated_at, is_deleted) VALUES\n"""
    
    if not cart_items:
        sql_content += "-- 暂无购物车数据\n"
        return sql_content
    
    # 生成INSERT语句
    cart_values = []
    for item in cart_items:
        checked_str = 'true' if item['checked'] else 'false'
        value = f"({item['id']}, {item['user']}, {item['goods']}, {item['nums']}, {checked_str}, NOW(), NOW(), false)"
        cart_values.append(value)
    
    sql_content += ',\n'.join(cart_values) + ';\n\n'
    
    # 添加统计查询
    sql_content += """-- 购物车数据统计查询
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
"""
    
    return sql_content

if __name__ == '__main__':
    print("开始生成购物车数据...")
    
    # 生成购物车数据
    cart_items = generate_cart_data(3500)
    
    print(f"生成了 {len(cart_items)} 个购物车商品记录")
    
    # 生成SQL文件
    sql_content = generate_cart_sql(cart_items)
    
    # 写入文件
    with open('/Users/walker/gitsource/github/joyshop_srvs/sql/schemas/06_shopping_cart.sql', 'w', encoding='utf-8') as f:
        f.write(sql_content)
    
    print("购物车数据SQL文件已生成: schemas/06_shopping_cart.sql")
    
    if cart_items:
        # 统计信息
        total_items = len(cart_items)
        checked_items = sum(1 for item in cart_items if item['checked'])
        total_quantity = sum(item['nums'] for item in cart_items)
        unique_users = len(set(item['user'] for item in cart_items))
        unique_goods = len(set(item['goods'] for item in cart_items))
        max_goods_id = max(item['goods'] for item in cart_items)
        min_goods_id = min(item['goods'] for item in cart_items)
        
        print(f"\n购物车统计:")
        print(f"- 总商品记录: {total_items}")
        print(f"- 已选中商品: {checked_items} ({checked_items/total_items*100:.1f}%)")
        print(f"- 商品总数量: {total_quantity}")
        print(f"- 涉及用户数: {unique_users}")
        print(f"- 涉及商品数: {unique_goods}")
        print(f"- 商品ID范围: {min_goods_id}-{max_goods_id}")
    else:
        print("所有用户购物车均为空")