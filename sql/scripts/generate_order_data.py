#!/usr/bin/env python3
"""
订单数据生成脚本  
生成800个订单，包含订单信息和订单商品详情
订单商品从1-3500的商品ID中随机选择
"""

import random
from datetime import datetime, timedelta

# 订单状态定义
ORDER_STATUS = {
    'UNPAID': 1,      # 待付款
    'PAID': 2,        # 已付款待发货  
    'SHIPPED': 3,     # 已发货待收货
    'COMPLETED': 4,   # 已完成
    'CANCELLED': 5,   # 已取消
    'REFUNDING': 6,   # 退款中
    'REFUNDED': 7     # 已退款
}

# 支付方式
PAY_TYPE_CHOICES = [1, 2, 3]  # 1:支付宝 2:微信 3:银行卡

def generate_order_data(order_count=800, max_goods_id=3500):
    """生成订单数据"""
    orders = []
    order_goods = []
    
    # 生成随机时间（最近3个月）
    end_date = datetime.now()
    start_date = end_date - timedelta(days=90)
    
    for order_id in range(1, order_count + 1):
        # 随机选择用户ID (1-15)
        user_id = random.randint(1, 15)
        
        # 随机生成订单时间
        order_time = start_date + timedelta(
            seconds=random.randint(0, int((end_date - start_date).total_seconds()))
        )
        
        # 随机选择订单状态 (偏向已完成状态)
        status_weights = [5, 10, 15, 50, 8, 7, 5]  # 对应各状态的权重
        status = random.choices(list(ORDER_STATUS.values()), weights=status_weights)[0]
        
        # 随机生成订单商品 (1-5个商品)
        goods_count = random.randint(1, 5)
        selected_goods = random.sample(range(1, max_goods_id + 1), goods_count)
        
        # 计算订单总金额
        order_mount = 0
        order_goods_items = []
        
        for goods_id in selected_goods:
            # 随机生成商品价格 (基于商品ID范围)
            if goods_id <= 500:  # 电子产品等高价商品
                price = random.uniform(1000, 20000)
            elif goods_id <= 1500:  # 中等价格商品
                price = random.uniform(100, 2000)  
            else:  # 低价商品
                price = random.uniform(20, 500)
            
            price = round(price, 2)
            nums = random.randint(1, 3)  # 购买数量
            
            order_goods_item = {
                'order_id': order_id,
                'goods_id': goods_id,
                'goods_name': f'商品_{goods_id}',  # 简化商品名
                'goods_image': f'placeholder_front_{goods_id}.jpg',
                'goods_price': price,
                'nums': nums,
                'created_at': order_time,
                'updated_at': order_time
            }
            
            order_goods_items.append(order_goods_item)
            order_mount += price * nums
        
        # 生成订单信息
        order = {
            'id': order_id,
            'user_id': user_id,
            'order_sn': f'JOY{order_time.strftime("%Y%m%d")}{order_id:06d}',
            'order_mount': round(order_mount, 2),
            'address': f'测试地址{user_id}-{order_id}',
            'singer_name': f'收货人{user_id}',
            'singer_mobile': f'138{user_id:04d}{random.randint(1000, 9999)}',
            'post': f'{random.randint(100000, 999999)}',
            'pay_type': random.choice(PAY_TYPE_CHOICES),
            'status': status,
            'trade_no': f'trade_{order_id}_{random.randint(100000, 999999)}' if status >= 2 else '',
            'order_time': order_time,
            'pay_time': order_time + timedelta(minutes=random.randint(5, 120)) if status >= 2 else None,
            'created_at': order_time,
            'updated_at': order_time + timedelta(hours=random.randint(0, 24))
        }
        
        orders.append(order)
        order_goods.extend(order_goods_items)
    
    return orders, order_goods

def generate_order_sql(orders, order_goods):
    """生成订单SQL文件内容"""
    sql_content = """-- 订单服务测试数据 (扩展版)
-- 创建800个订单信息和对应的订单商品数据，商品ID范围1-3500

-- 清空现有数据
DELETE FROM order_goods;
DELETE FROM order_info;

-- 重置自增ID
ALTER TABLE order_info AUTO_INCREMENT = 1;
ALTER TABLE order_goods AUTO_INCREMENT = 1;

-- 插入订单信息
INSERT INTO order_info (id, user_id, order_sn, order_mount, address, singer_name, singer_mobile, post, pay_type, status, trade_no, order_time, pay_time, created_at, updated_at) VALUES\n"""
    
    # 生成订单信息INSERT语句
    order_values = []
    for order in orders:
        pay_time_str = f"'{order['pay_time'].strftime('%Y-%m-%d %H:%M:%S')}'" if order['pay_time'] else 'NULL'
        trade_no_str = f"'{order['trade_no']}'" if order['trade_no'] else "''"
        
        value = f"({order['id']}, {order['user_id']}, '{order['order_sn']}', {order['order_mount']}, '{order['address']}', '{order['singer_name']}', '{order['singer_mobile']}', '{order['post']}', {order['pay_type']}, {order['status']}, {trade_no_str}, '{order['order_time'].strftime('%Y-%m-%d %H:%M:%S')}', {pay_time_str}, '{order['created_at'].strftime('%Y-%m-%d %H:%M:%S')}', '{order['updated_at'].strftime('%Y-%m-%d %H:%M:%S')}')"
        order_values.append(value)
    
    sql_content += ',\n'.join(order_values) + ';\n\n'
    
    # 生成订单商品INSERT语句
    sql_content += "-- 插入订单商品\nINSERT INTO order_goods (id, order_id, goods_id, goods_name, goods_image, goods_price, nums, created_at, updated_at, is_commented) VALUES\n"
    
    goods_values = []
    for i, goods in enumerate(order_goods, 1):
        value = f"({i}, {goods['order_id']}, {goods['goods_id']}, '{goods['goods_name']}', '{goods['goods_image']}', {goods['goods_price']}, {goods['nums']}, '{goods['created_at'].strftime('%Y-%m-%d %H:%M:%S')}', '{goods['updated_at'].strftime('%Y-%m-%d %H:%M:%S')}', false)"
        goods_values.append(value)
    
    sql_content += ',\n'.join(goods_values) + ';\n\n'
    
    # 添加统计查询
    sql_content += """-- 订单数据统计查询
-- 按状态统计订单数量
-- SELECT 
--     CASE status
--         WHEN 1 THEN '待付款'
--         WHEN 2 THEN '已付款待发货'
--         WHEN 3 THEN '已发货待收货'  
--         WHEN 4 THEN '已完成'
--         WHEN 5 THEN '已取消'
--         WHEN 6 THEN '退款中'
--         WHEN 7 THEN '已退款'
--     END as status_name,
--     COUNT(*) as order_count,
--     SUM(order_mount) as total_amount
-- FROM order_info 
-- GROUP BY status
-- ORDER BY status;

-- 商品销售排行 (订单商品表)
-- SELECT 
--     goods_id,
--     goods_name,
--     SUM(nums) as total_sold,
--     COUNT(DISTINCT order_id) as order_count,
--     SUM(goods_price * nums) as total_revenue
-- FROM order_goods 
-- GROUP BY goods_id, goods_name
-- ORDER BY total_sold DESC
-- LIMIT 20;

-- 用户订单统计
-- SELECT 
--     user_id,
--     COUNT(*) as order_count,
--     SUM(order_mount) as total_spent,
--     AVG(order_mount) as avg_order_amount
-- FROM order_info 
-- GROUP BY user_id
-- ORDER BY total_spent DESC;
"""
    
    return sql_content

if __name__ == '__main__':
    print("开始生成订单数据...")
    
    # 生成订单数据
    orders, order_goods = generate_order_data(800, 3500)
    
    print(f"生成了 {len(orders)} 个订单")
    print(f"生成了 {len(order_goods)} 个订单商品记录")
    
    # 生成SQL文件
    sql_content = generate_order_sql(orders, order_goods)
    
    # 写入文件
    with open('/Users/walker/gitsource/github/joyshop_srvs/sql/schemas/07_orders.sql', 'w', encoding='utf-8') as f:
        f.write(sql_content)
    
    print("订单数据SQL文件已生成: schemas/07_orders.sql")
    
    # 统计信息
    total_amount = sum(order['order_mount'] for order in orders)
    avg_amount = total_amount / len(orders)
    max_goods_id = max(goods['goods_id'] for goods in order_goods)
    min_goods_id = min(goods['goods_id'] for goods in order_goods)
    
    print(f"\n订单统计:")
    print(f"- 订单总金额: ¥{total_amount:,.2f}")
    print(f"- 平均订单金额: ¥{avg_amount:.2f}") 
    print(f"- 商品ID覆盖范围: {min_goods_id}-{max_goods_id}")
    print(f"- 平均每订单商品数: {len(order_goods)/len(orders):.1f}")
    
    # 状态分布统计
    status_count = {}
    for order in orders:
        status = order['status']
        status_count[status] = status_count.get(status, 0) + 1
    
    status_names = {1: '待付款', 2: '已付款', 3: '已发货', 4: '已完成', 5: '已取消', 6: '退款中', 7: '已退款'}
    print(f"\n订单状态分布:")
    for status, count in sorted(status_count.items()):
        print(f"- {status_names.get(status, f'状态{status}')}: {count} 订单")