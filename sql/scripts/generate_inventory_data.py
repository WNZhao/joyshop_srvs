#!/usr/bin/env python3
"""
库存数据生成脚本
为3500个商品生成对应的库存数据
"""

import random

def generate_inventory_data(goods_count=3500):
    """生成库存数据"""
    inventory_data = []
    
    for goods_id in range(1, goods_count + 1):
        # 根据商品ID范围确定库存策略
        if goods_id <= 500:  # 热门商品，高库存
            stock = random.randint(200, 500)
            sold = random.randint(50, int(stock * 0.4))
        elif goods_id <= 2000:  # 普通商品，中等库存
            stock = random.randint(80, 200)
            sold = random.randint(10, int(stock * 0.3))
        else:  # 长尾商品，低库存
            stock = random.randint(20, 100)
            sold = random.randint(0, int(stock * 0.2))
        
        # 冻结数量（购物车、待支付订单等）
        frozen = random.randint(0, min(10, stock - sold))
        
        inventory_item = {
            'goods_id': goods_id,
            'stock': stock,
            'sold': sold,
            'frozen': frozen
        }
        
        inventory_data.append(inventory_item)
    
    return inventory_data

def generate_inventory_sql(inventory_data):
    """生成库存SQL文件内容"""
    sql_content = """-- 库存测试数据 (扩展版)
-- 为3500个商品创建对应的库存数据

-- 清空现有数据
DELETE FROM inventory;

-- 重置自增ID
ALTER TABLE inventory AUTO_INCREMENT = 1;

-- 插入库存数据
INSERT INTO inventory (goods_id, stock, sold, frozen, created_at, updated_at) VALUES\n"""
    
    # 生成INSERT语句
    inventory_values = []
    for inventory in inventory_data:
        value = f"({inventory['goods_id']}, {inventory['stock']}, {inventory['sold']}, {inventory['frozen']}, NOW(), NOW())"
        inventory_values.append(value)
    
    sql_content += ',\n'.join(inventory_values) + ';\n\n'
    
    # 添加一些统计查询
    sql_content += """-- 库存数据统计查询
-- SELECT 
--     COUNT(*) as total_goods,
--     SUM(stock) as total_stock,
--     SUM(sold) as total_sold,
--     SUM(frozen) as total_frozen,
--     AVG(stock) as avg_stock,
--     MAX(stock) as max_stock,
--     MIN(stock) as min_stock
-- FROM inventory;

-- 按库存量分组统计
-- SELECT 
--     CASE 
--         WHEN stock >= 200 THEN '高库存(>=200)'
--         WHEN stock >= 80 THEN '中库存(80-199)'  
--         ELSE '低库存(<80)'
--     END as stock_level,
--     COUNT(*) as goods_count
-- FROM inventory 
-- GROUP BY 
--     CASE 
--         WHEN stock >= 200 THEN '高库存(>=200)'
--         WHEN stock >= 80 THEN '中库存(80-199)'
--         ELSE '低库存(<80)'
--     END;
"""
    
    return sql_content

if __name__ == '__main__':
    print("开始生成库存数据...")
    
    # 生成库存数据
    inventory_data = generate_inventory_data(3500)
    
    print(f"生成了 {len(inventory_data)} 条库存记录")
    
    # 生成SQL文件
    sql_content = generate_inventory_sql(inventory_data)
    
    # 写入文件
    with open('/Users/walker/gitsource/github/joyshop_srvs/sql/schemas/05_inventory.sql', 'w', encoding='utf-8') as f:
        f.write(sql_content)
    
    print("库存数据SQL文件已生成: schemas/05_inventory.sql")
    
    # 统计信息
    total_stock = sum(item['stock'] for item in inventory_data)
    total_sold = sum(item['sold'] for item in inventory_data)
    avg_stock = total_stock / len(inventory_data)
    
    print(f"\n库存统计:")
    print(f"- 总库存量: {total_stock:,}")
    print(f"- 总销售量: {total_sold:,}")
    print(f"- 平均库存: {avg_stock:.1f}")
    print(f"- 库存覆盖率: 100% (3500/3500商品)")