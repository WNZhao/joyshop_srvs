#!/usr/bin/env python3
"""
商品数据生成脚本
生成3500个商品数据，包括商品基础信息、分类关联、库存数据等
"""

import random
import json
from datetime import datetime, timedelta

# 基础商品模板
GOODS_TEMPLATES = {
    # 电子数码类 (分类1)
    'electronics': [
        {'name': 'iPhone 15 Pro Max', 'specs': ['128GB', '256GB', '512GB', '1TB'], 'colors': ['深空黑色', '原色钛金属', '白色钛金属', '蓝色钛金属'], 'brand_id': 1, 'base_price': 9999, 'categories': [1, 11, 111]},
        {'name': 'iPhone 14', 'specs': ['128GB', '256GB', '512GB'], 'colors': ['午夜色', '星光色', '蓝色', '紫色', '红色'], 'brand_id': 1, 'base_price': 5999, 'categories': [1, 11, 111]},
        {'name': 'HUAWEI Mate 60 Pro', 'specs': ['256GB', '512GB', '1TB'], 'colors': ['雅川青', '雅丹黑', '白羽砂'], 'brand_id': 2, 'base_price': 6999, 'categories': [1, 11, 112]},
        {'name': 'HUAWEI P60 Pro', 'specs': ['128GB', '256GB', '512GB'], 'colors': ['洛可可白', '翡冷翠', '雅黑'], 'brand_id': 2, 'base_price': 5988, 'categories': [1, 11, 112]},
        {'name': 'Xiaomi 14 Ultra', 'specs': ['12GB+256GB', '16GB+512GB', '16GB+1TB'], 'colors': ['黑色', '白色', '钛金属'], 'brand_id': 3, 'base_price': 5999, 'categories': [1, 11, 112]},
        {'name': 'MacBook Pro', 'specs': ['14英寸 M3', '16英寸 M3 Pro', '16英寸 M3 Max'], 'colors': ['深空灰色', '银色'], 'brand_id': 1, 'base_price': 15999, 'categories': [1, 12, 121]},
        {'name': 'MacBook Air', 'specs': ['13英寸 M2', '15英寸 M2'], 'colors': ['午夜色', '星光色', '银色', '深空灰色'], 'brand_id': 1, 'base_price': 8999, 'categories': [1, 12, 121]},
        {'name': 'ThinkPad X1 Carbon', 'specs': ['Gen10', 'Gen11'], 'colors': ['碳纤维黑', '银色'], 'brand_id': 5, 'base_price': 12999, 'categories': [1, 12, 121]},
    ],
    # 运动鞋帽类 (分类23)
    'sports': [
        {'name': 'Nike Air Jordan 1', 'specs': ['Retro High OG', 'Retro Mid', 'Retro Low'], 'colors': ['黑红配色', '白黑配色', '芝加哥配色'], 'brand_id': 11, 'base_price': 1299, 'categories': [2, 23, 231]},
        {'name': 'Adidas Ultraboost', 'specs': ['22', '23', '24'], 'colors': ['黑白配色', '全白', '全黑', '灰色'], 'brand_id': 12, 'base_price': 1599, 'categories': [2, 23, 231]},
        {'name': 'Nike Air Force 1', 'specs': ['Low', 'Mid', 'High'], 'colors': ['全白', '全黑', '白黑配色'], 'brand_id': 11, 'base_price': 899, 'categories': [2, 23, 231]},
    ],
    # 服装类 (分类21)
    'clothing': [
        {'name': 'UNIQLO 纯棉T恤', 'specs': ['男士', '女士', '儿童'], 'colors': ['白色', '黑色', '灰色', '蓝色', '红色'], 'brand_id': 13, 'base_price': 79, 'categories': [2, 21, 211]},
        {'name': "Levis 501牛仔裤", 'specs': ['经典直筒', '修身版型', '宽松版型'], 'colors': ['深蓝色', '浅蓝色', '黑色', '灰色'], 'brand_id': 16, 'base_price': 699, 'categories': [2, 21, 212]},
        {'name': 'ZARA 连衣裙', 'specs': ['春季款', '夏季款', '秋季款'], 'colors': ['黑色', '白色', '红色', '蓝色', '印花款'], 'brand_id': 14, 'base_price': 299, 'categories': [2, 21, 213]},
    ],
    # 家电类 (分类13)
    'appliance': [
        {'name': '美的变频空调', 'specs': ['1匹', '1.5匹', '2匹', '3匹'], 'colors': ['白色', '银色'], 'brand_id': 23, 'base_price': 2999, 'categories': [1, 13, 131]},
        {'name': '海尔冰箱', 'specs': ['双开门', '三开门', '对开门'], 'colors': ['白色', '银色', '金色'], 'brand_id': 24, 'base_price': 4199, 'categories': [1, 13, 132]},
        {'name': 'Xiaomi电视', 'specs': ['55英寸', '65英寸', '75英寸', '85英寸'], 'colors': ['黑色'], 'brand_id': 3, 'base_price': 2999, 'categories': [1, 13, 133]},
    ],
    # 图书类 (分类31)
    'books': [
        {'name': 'Python编程：从入门到实践', 'specs': ['第1版', '第2版', '第3版'], 'colors': ['标准版'], 'brand_id': 32, 'base_price': 89, 'categories': [3, 31, 311]},
        {'name': '深入理解计算机系统', 'specs': ['原书第3版', '中文版'], 'colors': ['标准版'], 'brand_id': 31, 'base_price': 139, 'categories': [3, 31, 311]},
        {'name': 'Java核心技术', 'specs': ['卷I', '卷II', '套装'], 'colors': ['标准版'], 'brand_id': 31, 'base_price': 159, 'categories': [3, 31, 311]},
    ],
    # 美妆护肤类 (分类24)
    'beauty': [
        {'name': 'SK-II 神仙水', 'specs': ['75ml', '160ml', '230ml', '330ml'], 'colors': ['标准版'], 'brand_id': 53, 'base_price': 1690, 'categories': [2, 24, 241]},
        {'name': "LOreal 面膜", 'specs': ['补水款', '美白款', '抗老款'], 'colors': ['标准版'], 'brand_id': 51, 'base_price': 199, 'categories': [2, 24, 242]},
    ],
    # 家具类 (分类22)
    'furniture': [
        {'name': 'IKEA 床架', 'specs': ['单人床', '双人床', '大号双人床'], 'colors': ['白色', '黑色', '原木色'], 'brand_id': 21, 'base_price': 899, 'categories': [3, 22, 221]},
        {'name': 'IKEA 衣柜', 'specs': ['2门', '3门', '4门'], 'colors': ['白色', '黑色', '原木色'], 'brand_id': 21, 'base_price': 1299, 'categories': [3, 22, 222]},
    ]
}

def generate_goods_data(target_count=3500):
    """生成商品数据"""
    goods_data = []
    goods_category_data = []
    current_id = 1
    
    # 计算每个类别需要生成的数量
    template_keys = list(GOODS_TEMPLATES.keys())
    items_per_category = target_count // len(template_keys)
    
    for category_name, templates in GOODS_TEMPLATES.items():
        category_count = 0
        
        while category_count < items_per_category and current_id <= target_count:
            template = random.choice(templates)
            
            # 生成商品变体
            spec = random.choice(template['specs'])
            color = random.choice(template['colors'])
            
            # 构建商品名称
            if template['colors'] == ['标准版']:
                goods_name = f"{template['name']} {spec}"
            else:
                goods_name = f"{template['name']} {spec} {color}"
            
            # 计算价格（基础价格 +/- 20%的随机浮动）
            price_variation = random.uniform(0.8, 1.2)
            market_price = round(template['base_price'] * price_variation, 2)
            shop_price = round(market_price * random.uniform(0.85, 0.95), 2)  # 打85-95折
            
            # 生成商品数据
            goods_item = {
                'id': current_id,
                'brand_id': template['brand_id'],
                'on_sale': random.choice([True, True, True, False]),  # 75%概率上架
                'ship_free': random.choice([True, False]),
                'is_new': random.choice([True, False]),
                'is_hot': random.choice([True, False, False]),  # 33%概率热门
                'name': goods_name,
                'goods_sn': f'GD{current_id:06d}',
                'click_num': random.randint(50, 10000),
                'fav_num': random.randint(10, 1000),
                'market_price': market_price,
                'shop_price': shop_price,
                'goods_brief': f'{goods_name}，优质商品，值得信赖',
                'images': f'["placeholder_image_{current_id}_1.jpg","placeholder_image_{current_id}_2.jpg","placeholder_image_{current_id}_3.jpg"]',
                'desc_images': f'["placeholder_desc_{current_id}_1.jpg","placeholder_desc_{current_id}_2.jpg"]',
                'goods_front_image': f'placeholder_front_{current_id}.jpg',
                'status': 1
            }
            
            goods_data.append(goods_item)
            
            # 生成商品分类关联
            for category_id in template['categories']:
                goods_category_data.append({
                    'goods_id': current_id,
                    'category_id': category_id
                })
            
            current_id += 1
            category_count += 1
    
    return goods_data, goods_category_data

def generate_sql_file(goods_data, goods_category_data):
    """生成SQL文件内容"""
    sql_content = """-- 商品测试数据 (扩展版)
-- 创建3500个商品数据，包含商品信息、价格、图片等

-- 清空现有数据
DELETE FROM goods_category;
DELETE FROM goods;

-- 重置自增ID  
ALTER TABLE goods AUTO_INCREMENT = 1;

-- 插入商品数据
INSERT INTO goods (id, brand_id, on_sale, ship_free, is_new, is_hot, name, goods_sn, click_num, fav_num, market_price, shop_price, goods_brief, images, desc_images, goods_front_image, status, created_at, updated_at) VALUES\n"""
    
    # 生成商品INSERT语句
    goods_values = []
    for i, goods in enumerate(goods_data):
        on_sale = 'true' if goods['on_sale'] else 'false'
        ship_free = 'true' if goods['ship_free'] else 'false' 
        is_new = 'true' if goods['is_new'] else 'false'
        is_hot = 'true' if goods['is_hot'] else 'false'
        
        value = f"({goods['id']}, {goods['brand_id']}, {on_sale}, {ship_free}, {is_new}, {is_hot}, '{goods['name']}', '{goods['goods_sn']}', {goods['click_num']}, {goods['fav_num']}, {goods['market_price']}, {goods['shop_price']}, '{goods['goods_brief']}', '{goods['images']}', '{goods['desc_images']}', '{goods['goods_front_image']}', {goods['status']}, NOW(), NOW())"
        goods_values.append(value)
    
    sql_content += ',\n'.join(goods_values) + ';\n\n'
    
    # 生成分类关联INSERT语句
    sql_content += "-- 建立商品与分类的关联关系\nINSERT INTO goods_category (goods_id, category_id) VALUES\n"
    
    category_values = []
    for relation in goods_category_data:
        value = f"({relation['goods_id']}, {relation['category_id']})"
        category_values.append(value)
    
    sql_content += ',\n'.join(category_values) + ';\n'
    
    return sql_content

if __name__ == '__main__':
    print("开始生成商品数据...")
    
    # 生成数据
    goods_data, goods_category_data = generate_goods_data(3500)
    
    print(f"生成了 {len(goods_data)} 个商品")
    print(f"生成了 {len(goods_category_data)} 个分类关联")
    
    # 生成SQL文件
    sql_content = generate_sql_file(goods_data, goods_category_data)
    
    # 写入文件
    with open('/Users/walker/gitsource/github/joyshop_srvs/sql/schemas/04_goods.sql', 'w', encoding='utf-8') as f:
        f.write(sql_content)
    
    print("商品数据SQL文件已生成: schemas/04_goods.sql")
    print("请执行该文件以导入数据")