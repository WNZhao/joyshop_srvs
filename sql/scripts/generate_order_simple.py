#!/usr/bin/env python3
"""
生成订单数据，使用数据库中实际的商品信息
不依赖外部库，直接使用mysql命令获取数据
"""

import random
import subprocess
import json
from datetime import datetime, timedelta

def get_goods_from_db():
    """通过mysql命令获取商品数据"""
    cmd = [
        '/usr/local/opt/mysql-client/bin/mysql',
        '-h', '127.0.0.1',
        '-u', 'root',
        '-pwalker123',
        'joyshop_goods',
        '-e', 'SELECT id, name, goods_front_image, shop_price FROM goods WHERE on_sale = true LIMIT 100;',
        '--batch',  # 使用tab分隔的输出
        '--skip-column-names'  # 不显示列名
    ]
    
    try:
        result = subprocess.run(cmd, capture_output=True, text=True)
        if result.returncode == 0:
            goods_data = []
            for line in result.stdout.strip().split('\n'):
                if line:
                    parts = line.split('\t')
                    if len(parts) >= 4:
                        goods_data.append({
                            'id': int(parts[0]),
                            'name': parts[1],
                            'image': parts[2],
                            'price': float(parts[3])
                        })
            return goods_data
    except Exception as e:
        print(f"获取商品数据失败: {e}")
    
    return None

def generate_order_sql_simple():
    """生成订单SQL文件（简化版）"""
    
    print("正在从数据库获取商品数据...")
    goods_data = get_goods_from_db()
    
    if not goods_data:
        print("无法获取商品数据，使用默认测试数据...")
        # 使用一些默认的商品数据（基于之前生成的数据）
        goods_data = [
            {'id': 1, 'name': 'MacBook Air 13英寸 M2 午夜色', 'image': 'placeholder_front_1.jpg', 'price': 9183.56},
            {'id': 2, 'name': 'MacBook Air 13英寸 M2 午夜色', 'image': 'placeholder_front_2.jpg', 'price': 7129.68},
            {'id': 3, 'name': 'iPhone 15 Pro Max 512GB 蓝色钛金属', 'image': 'placeholder_front_3.jpg', 'price': 7730.54},
            {'id': 4, 'name': 'MacBook Pro 14英寸 M3 深空灰色', 'image': 'placeholder_front_4.jpg', 'price': 14804.78},
            {'id': 5, 'name': 'HUAWEI P60 Pro 512GB 翡冷翠', 'image': 'placeholder_front_5.jpg', 'price': 4767.49},
            {'id': 10, 'name': 'iPhone 15 Pro Max 1TB 蓝色钛金属', 'image': 'placeholder_front_10.jpg', 'price': 9502.23},
            {'id': 15, 'name': 'Xiaomi 14 Ultra 16GB+1TB 钛金属', 'image': 'placeholder_front_15.jpg', 'price': 6989.00},
            {'id': 20, 'name': 'ThinkPad X1 Carbon Gen10 碳纤维黑', 'image': 'placeholder_front_20.jpg', 'price': 12999.00},
            {'id': 25, 'name': 'Nike Air Jordan 1 Retro High OG 黑红配色', 'image': 'placeholder_front_25.jpg', 'price': 1299.00},
            {'id': 30, 'name': 'Adidas Ultraboost 23 全白', 'image': 'placeholder_front_30.jpg', 'price': 1599.00}
        ]
    
    print(f"获取到 {len(goods_data)} 个商品")
    
    # 用户ID列表（基于之前导入的用户数据）
    user_ids = list(range(1, 16))  # 1-15号用户
    
    # 生成订单数据
    order_info_sql = []
    order_goods_sql = []
    
    # 生成100个订单
    for order_id in range(1, 101):
        user_id = random.choice(user_ids)
        
        # 生成订单基本信息
        order_date = datetime.now() - timedelta(days=random.randint(0, 60))
        order_sn = f'JOY{order_date.strftime("%Y%m%d")}{order_id:06d}'
        address = f'测试地址{user_id}-{order_id}'
        signer_name = f'收货人{user_id}'
        singer_mobile = f'138{random.randint(10000000, 99999999):08d}'
        post = f'订单{order_id}'
        pay_type = random.choice(['alipay', 'wechat'])
        status_choices = ['PAYING', 'TRADE_SUCCESS', 'TRADE_FINISHED', 'TRADE_CLOSED', 'WAIT_BUYER_PAY']
        status = random.choice(status_choices)
        
        # 根据状态生成交易号和支付时间
        if status in ['TRADE_SUCCESS', 'TRADE_FINISHED']:
            trade_no = f'trade_{order_id}_{random.randint(100000, 999999)}'
            pay_time = (order_date + timedelta(hours=random.randint(1, 24))).strftime('%Y-%m-%d %H:%M:%S')
        else:
            trade_no = ''
            pay_time = 'NULL'
        
        # 为该订单随机选择1-5个商品
        num_goods = random.randint(1, min(5, len(goods_data)))
        selected_goods = random.sample(goods_data, num_goods)
        
        # 计算订单总金额
        order_mount = 0
        for goods in selected_goods:
            nums = random.randint(1, 3)
            order_mount += goods['price'] * nums
            
            # 确保商品名称中的单引号被正确转义
            goods_name = goods['name'].replace("'", "\\'")
            
            # 生成订单商品记录
            order_goods_sql.append(
                f"({order_id}, {goods['id']}, '{goods_name}', '{goods['image']}', {goods['price']:.2f}, {nums})"
            )
        
        # 生成订单信息记录
        pay_time_str = f"'{pay_time}'" if pay_time != 'NULL' else pay_time
        order_info_sql.append(
            f"({user_id}, '{order_sn}', {order_mount:.2f}, '{address}', '{signer_name}', '{singer_mobile}', '{post}', '{pay_type}', '{status}', '{trade_no}', {pay_time_str})"
        )
    
    # 生成最终的SQL文件内容
    sql_content = f"""-- 基于真实商品数据的订单信息
-- 生成时间: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}
-- 说明：订单商品表中的商品信息与goods表保持一致（冗余存储）

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- 清空现有数据
DELETE FROM order_goods;
DELETE FROM order_info;

-- 重置自增ID
ALTER TABLE order_info AUTO_INCREMENT = 1;
ALTER TABLE order_goods AUTO_INCREMENT = 1;

-- 插入订单信息（100个订单）
INSERT INTO order_info (user, order_sn, order_mount, address, signer_name, singer_mobile, post, pay_type, status, trade_no, pay_time) VALUES
{',\n'.join(order_info_sql)};

-- 插入订单商品信息（与实际商品数据保持一致的冗余数据）
INSERT INTO order_goods (`order`, goods, goods_name, goods_image, goods_price, nums) VALUES
{',\n'.join(order_goods_sql)};

-- 恢复外键检查
SET FOREIGN_KEY_CHECKS = 1;

-- 数据统计
SELECT '✅ 订单数据导入完成' as status;
SELECT '订单总数' as stat_name, COUNT(*) as count FROM order_info
UNION ALL
SELECT '订单商品记录数', COUNT(*) FROM order_goods
UNION ALL
SELECT '平均订单金额', ROUND(AVG(order_mount), 2) FROM order_info;
"""
    
    # 写入文件
    output_file = '/Users/walker/gitsource/github/joyshop_srvs/sql/schemas/07_orders.sql'
    with open(output_file, 'w', encoding='utf-8') as f:
        f.write(sql_content)
    
    print(f"订单数据生成完成！")
    print(f"- 生成了 {len(order_info_sql)} 个订单")
    print(f"- 生成了 {len(order_goods_sql)} 条订单商品记录")
    print(f"- 数据已保存到: schemas/07_orders.sql")
    print("\n订单数据特点：")
    print("- 订单商品信息与goods表中的数据保持一致")
    print("- 包含多种订单状态（待支付、已完成、已关闭等）")
    print("- 每个订单包含1-5个商品")
    print("- 订单时间分布在最近60天内")

if __name__ == "__main__":
    generate_order_sql_simple()