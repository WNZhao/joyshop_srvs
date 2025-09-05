#!/usr/bin/env python3
"""
生成订单数据，使用真实的商品信息
确保订单商品表中的数据与商品表保持一致
"""

import random
import pymysql
from datetime import datetime, timedelta

def get_db_connection():
    """获取数据库连接"""
    return pymysql.connect(
        host='127.0.0.1',
        user='root',
        password='walker123',
        charset='utf8mb4'
    )

def get_real_goods_data():
    """从数据库获取真实的商品数据"""
    conn = get_db_connection()
    cursor = conn.cursor()
    
    # 获取商品数据
    cursor.execute("""
        SELECT id, name, goods_front_image, shop_price 
        FROM joyshop_goods.goods 
        WHERE on_sale = true 
        LIMIT 500
    """)
    
    goods = cursor.fetchall()
    cursor.close()
    conn.close()
    
    return goods

def get_user_ids():
    """获取用户ID列表"""
    conn = get_db_connection()
    cursor = conn.cursor()
    
    cursor.execute("SELECT id FROM joyshop_user.user")
    users = [row[0] for row in cursor.fetchall()]
    
    cursor.close()
    conn.close()
    
    return users if users else [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15]

def generate_order_sql():
    """生成订单SQL文件"""
    
    print("正在从数据库获取商品数据...")
    try:
        goods_data = get_real_goods_data()
        user_ids = get_user_ids()
    except Exception as e:
        print(f"获取数据失败: {e}")
        print("使用默认测试数据...")
        # 使用默认数据
        goods_data = [
            (1, 'MacBook Air 13英寸 M2 午夜色', 'placeholder_front_1.jpg', 9183.56),
            (2, 'MacBook Air 13英寸 M2 午夜色', 'placeholder_front_2.jpg', 7129.68),
            (3, 'iPhone 15 Pro Max 512GB 蓝色钛金属', 'placeholder_front_3.jpg', 7730.54),
            (4, 'MacBook Pro 14英寸 M3 深空灰色', 'placeholder_front_4.jpg', 14804.78),
            (5, 'HUAWEI P60 Pro 512GB 翡冷翠', 'placeholder_front_5.jpg', 4767.49),
        ]
        user_ids = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15]
    
    print(f"获取到 {len(goods_data)} 个商品")
    print(f"获取到 {len(user_ids)} 个用户")
    
    # 生成订单数据
    order_info_sql = []
    order_goods_sql = []
    
    # 生成50个订单
    for order_id in range(1, 51):
        user_id = random.choice(user_ids)
        
        # 生成订单基本信息
        order_sn = f'JOY2025{random.randint(1000000000, 9999999999)}'
        order_mount = 0  # 稍后计算
        address = f'测试地址{user_id}-{order_id}'
        signer_name = f'收货人{user_id}'
        singer_mobile = f'138{random.randint(10000000, 99999999)}'
        post = f'订单备注{order_id}'
        pay_type = random.choice(['alipay', 'wechat'])
        status = random.choice(['PAYING', 'TRADE_SUCCESS', 'TRADE_FINISHED', 'TRADE_CLOSED'])
        trade_no = f'trade_{order_id}_{random.randint(100000, 999999)}' if status != 'PAYING' else ''
        
        # 生成支付时间
        if status in ['TRADE_SUCCESS', 'TRADE_FINISHED']:
            days_ago = random.randint(1, 30)
            pay_time = (datetime.now() - timedelta(days=days_ago)).strftime('%Y-%m-%d %H:%M:%S')
        else:
            pay_time = 'NULL'
        
        # 为该订单随机选择1-3个商品
        num_goods = random.randint(1, 3)
        selected_goods = random.sample(goods_data, min(num_goods, len(goods_data)))
        
        # 计算订单总金额
        order_mount = 0
        for goods in selected_goods:
            goods_id, goods_name, goods_image, goods_price = goods
            nums = random.randint(1, 3)
            order_mount += float(goods_price) * nums
            
            # 生成订单商品记录
            order_goods_sql.append(
                f"({order_id}, {goods_id}, '{goods_name}', '{goods_image}', {goods_price:.2f}, {nums})"
            )
        
        # 生成订单信息记录
        pay_time_str = f"'{pay_time}'" if pay_time != 'NULL' else pay_time
        order_info_sql.append(
            f"({user_id}, '{order_sn}', {order_mount:.2f}, '{address}', '{signer_name}', '{singer_mobile}', '{post}', '{pay_type}', '{status}', '{trade_no}', {pay_time_str})"
        )
    
    # 生成最终的SQL文件内容
    sql_content = f"""-- 基于真实商品数据的订单信息
-- 生成时间: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- 清空现有数据
DELETE FROM order_goods;
DELETE FROM order_info;

-- 重置自增ID
ALTER TABLE order_info AUTO_INCREMENT = 1;
ALTER TABLE order_goods AUTO_INCREMENT = 1;

-- 插入订单信息
INSERT INTO order_info (user, order_sn, order_mount, address, signer_name, singer_mobile, post, pay_type, status, trade_no, pay_time) VALUES
{',\n'.join(order_info_sql)};

-- 插入订单商品信息（与实际商品数据保持一致）
INSERT INTO order_goods (`order`, goods, goods_name, goods_image, goods_price, nums) VALUES
{',\n'.join(order_goods_sql)};

-- 恢复外键检查
SET FOREIGN_KEY_CHECKS = 1;

-- 数据统计
SELECT '订单总数' as stat_name, COUNT(*) as count FROM order_info
UNION ALL
SELECT '订单商品记录数', COUNT(*) FROM order_goods
UNION ALL
SELECT '订单总金额', SUM(order_mount) FROM order_info;
"""
    
    # 写入文件
    output_file = '/Users/walker/gitsource/github/joyshop_srvs/sql/schemas/07_orders.sql'
    with open(output_file, 'w', encoding='utf-8') as f:
        f.write(sql_content)
    
    print(f"订单数据生成完成！")
    print(f"- 生成了 {len(order_info_sql)} 个订单")
    print(f"- 生成了 {len(order_goods_sql)} 条订单商品记录")
    print(f"- 数据已保存到: {output_file}")

if __name__ == "__main__":
    generate_order_sql()