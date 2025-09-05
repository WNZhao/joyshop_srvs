#!/usr/bin/env python3
"""
根据order.go模型精确修复订单数据

OrderInfo模型字段映射：
- User -> user (int32)
- OrderSn -> order_sn (varchar(30))  
- PayType -> pay_type (varchar(20))
- Status -> status (varchar(20))
- TradeNo -> trade_no (varchar(100))
- OrderMount -> order_mount (float32)
- PayTime -> pay_time (datetime)
- Address -> address (varchar(100))
- SignerName -> signer_name (varchar(20))
- SingerMobile -> singer_mobile (varchar(11))  # 注意：确实是SingerMobile
- Post -> post (varchar(20))

OrderGoods模型字段映射：
- Order -> order (int32)  # 不是order_id！
- Goods -> goods (int32)
- GoodsName -> goods_name (varchar(100))
- GoodsImage -> goods_image (varchar(200))
- GoodsPrice -> goods_price (float32)
- Nums -> nums (int32)
"""

import re

def fix_order_data():
    """根据模型定义精确修复订单数据"""
    file_path = "/Users/walker/gitsource/github/joyshop_srvs/sql/schemas/07_orders.sql"
    
    print(f"根据order.go模型修复订单数据: {file_path}")
    
    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # 1. 修复OrderInfo的INSERT语句
    # 正确的字段顺序，基于实际数据库表结构
    order_info_insert = """INSERT INTO order_info (id, user, order_sn, order_mount, address, signer_name, singer_mobile, post, pay_type, status, trade_no, pay_time, created_at, updated_at) VALUES"""
    
    # 查找并替换原有的INSERT语句
    content = re.sub(
        r'INSERT INTO order_info.*?VALUES',
        order_info_insert,
        content,
        flags=re.DOTALL
    )
    
    # 2. 修复OrderGoods的INSERT语句  
    # 正确字段：order (不是order_id), goods, goods_name, goods_image, goods_price, nums
    order_goods_insert = """INSERT INTO order_goods (id, order, goods, goods_name, goods_image, goods_price, nums, created_at, updated_at) VALUES"""
    
    content = re.sub(
        r'INSERT INTO order_goods.*?VALUES',
        order_goods_insert,
        content,
        flags=re.DOTALL
    )
    
    # 3. 修复所有order_id为order
    content = content.replace('order_id', 'order')
    
    # 4. 确保singer_mobile字段名正确（根据模型，确实是singer_mobile）
    # 模型中定义的是SingerMobile，在数据库中应该是singer_mobile
    
    # 5. 移除order_time字段（如果存在）
    lines = content.split('\n')
    fixed_lines = []
    
    for line in lines:
        # 处理OrderInfo数据行 - 移除多余的order_time字段
        if line.strip().startswith('(') and 'JOY202' in line:
            # 订单数据行，需要确保字段数量正确
            # 预期格式：(id, user, order_sn, order_mount, address, signer_name, singer_mobile, post, pay_type, status, trade_no, pay_time, created_at, updated_at)
            # 如果有额外的order_time字段，需要移除
            parts = line.split(',')
            if len(parts) > 14:  # 超过14个字段说明有多余字段
                # 移除第12个位置的order_time字段（如果存在）
                new_parts = parts[:11] + parts[12:]
                line = ','.join(new_parts)
            
            fixed_lines.append(line)
        else:
            fixed_lines.append(line)
    
    content = '\n'.join(fixed_lines)
    
    # 6. 添加SQL设置
    header = """-- 根据order.go模型修复的订单数据
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

"""
    
    footer = """
-- 恢复外键检查
SET FOREIGN_KEY_CHECKS = 1;
"""
    
    if not content.startswith('-- 根据order.go模型修复'):
        content = header + content + footer
    
    with open(file_path, 'w', encoding='utf-8') as f:
        f.write(content)
    
    print("订单数据修复完成！")
    print("修复内容：")
    print("- OrderInfo字段：user, order_sn, pay_type, status, trade_no, order_mount, pay_time, address, signer_name, singer_mobile, post")
    print("- OrderGoods字段：order (不是order_id), goods, goods_name, goods_image, goods_price, nums")
    print("- 移除多余的order_time字段")
    print("- 确保字段顺序与数据库表结构匹配")

if __name__ == "__main__":
    fix_order_data()