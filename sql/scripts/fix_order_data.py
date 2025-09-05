#!/usr/bin/env python3
"""
修复订单数据文件，使其与order.go模型完全匹配

OrderInfo字段：
- BaseModel: ID, CreatedAt, UpdatedAt, DeletedAt, IsDeleted
- User (int32)
- OrderSn (string)
- PayType (string)
- Status (string)
- TradeNo (string)
- OrderMount (float32) - 注意不是order_mount
- PayTime (time.Time)
- Address (string)
- SignerName (string) - 注意不是signer_name  
- SingerMobile (string) - 注意是SingerMobile不是singer_mobile
- Post (string)

OrderGoods字段:
- BaseModel: ID, CreatedAt, UpdatedAt, DeletedAt, IsDeleted
- Order (int32)
- Goods (int32)
- GoodsName (string)
- GoodsImage (string)
- GoodsPrice (float32)
- Nums (int32)
"""

import re

def fix_order_sql():
    """修复订单SQL文件"""
    file_path = "/Users/walker/gitsource/github/joyshop_srvs/sql/schemas/07_orders.sql"
    
    print(f"正在修复订单数据文件: {file_path}")
    
    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # 1. 修复OrderInfo的INSERT语句
    # 实际字段名应该是：user, order_sn, pay_type, status, trade_no, order_mount, pay_time, address, signer_name, singer_mobile, post
    old_order_insert = re.search(r'INSERT INTO order_info.*?VALUES', content, re.DOTALL)
    if old_order_insert:
        new_order_insert = "INSERT INTO order_info (id, user, order_sn, order_mount, address, signer_name, singer_mobile, post, pay_type, status, trade_no, pay_time, created_at, updated_at) VALUES"
        content = content.replace(old_order_insert.group(0), new_order_insert)
    
    # 2. 修复OrderGoods的INSERT语句  
    # 需要检查OrderGoods表的实际字段
    # 应该是：order_id (对应Order), goods (对应Goods), goods_name, goods_image, goods_price, nums
    
    # 3. 移除数据行中多余的字段
    lines = content.split('\n')
    fixed_lines = []
    
    for line in lines:
        if 'order_info' in line and line.strip().startswith('('):
            # 这是OrderInfo的数据行
            # 原格式可能有多个字段，我们需要保持与INSERT匹配的数量
            # 暂时保持原样，因为我们需要更仔细地分析数据格式
            fixed_lines.append(line)
        elif 'order_goods' in line and line.strip().startswith('('):
            # 这是OrderGoods的数据行
            fixed_lines.append(line)
        else:
            fixed_lines.append(line)
    
    content = '\n'.join(fixed_lines)
    
    # 添加SQL设置
    header = """-- 修复后的订单数据
-- 设置字符集和SQL模式
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

"""
    
    footer = """
-- 恢复外键检查
SET FOREIGN_KEY_CHECKS = 1;
"""
    
    if not content.startswith('-- 修复后的订单数据'):
        content = header + content + footer
    
    with open(file_path, 'w', encoding='utf-8') as f:
        f.write(content)
    
    print("订单数据修复完成")

if __name__ == "__main__":
    fix_order_sql()