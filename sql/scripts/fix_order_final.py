#!/usr/bin/env python3
"""
基于实际表结构修复订单数据

实际表结构：
order_info: id, created_at, updated_at, deleted_at, is_deleted, user, order_sn, pay_type, status, trade_no, order_mount, pay_time, address, signer_name, singer_mobile, post
order_goods: id, created_at, updated_at, deleted_at, is_deleted, order, goods, goods_name, goods_image, goods_price, nums
"""

def fix_order_sql():
    """修复订单SQL文件"""
    file_path = "/Users/walker/gitsource/github/joyshop_srvs/sql/schemas/07_orders.sql"
    
    print(f"正在基于实际表结构修复订单数据: {file_path}")
    
    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # 简单的修复：删除有问题的字段，重建INSERT语句
    # 1. 找到并替换order_info的INSERT语句
    content = content.replace(
        'INSERT INTO order_info (id, user, order_sn, order_mount, address, singer_name, singer_mobile, post, pay_type, status, trade_no, order_time, pay_time, created_at, updated_at)',
        'INSERT INTO order_info (id, user, order_sn, order_mount, address, signer_name, singer_mobile, post, pay_type, status, trade_no, pay_time, created_at, updated_at)'
    )
    
    # 2. 修复字段名错误
    content = content.replace('singer_name', 'signer_name')
    
    # 3. 移除不存在的order_time字段的数据
    lines = content.split('\n')
    fixed_lines = []
    
    for line in lines:
        if line.strip().startswith('(') and 'JOY202' in line:  # 这是订单数据行
            # 移除order_time字段 (第12个字段)
            # 原格式: (id, user, order_sn, order_mount, address, signer_name, singer_mobile, post, pay_type, status, trade_no, order_time, pay_time, created_at, updated_at)
            # 新格式: (id, user, order_sn, order_mount, address, signer_name, singer_mobile, post, pay_type, status, trade_no, pay_time, created_at, updated_at)
            
            # 使用正则表达式移除order_time字段
            import re
            # 找到第11个逗号后到第12个逗号的内容并删除
            parts = line.split(',')
            if len(parts) >= 15:  # 确保有足够的字段
                # 移除第12个字段 (order_time)
                new_parts = parts[:11] + parts[12:]  # 跳过第12个字段 (索引11)
                line = ','.join(new_parts)
            
            fixed_lines.append(line)
        else:
            fixed_lines.append(line)
    
    content = '\n'.join(fixed_lines)
    
    # 添加SQL设置
    header = """-- 修复后的订单数据
-- 基于实际表结构修复
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
    print("修复内容：")
    print("- 移除不存在的order_time字段")
    print("- 修复signer_name字段名")
    print("- 调整字段顺序匹配实际表结构")

if __name__ == "__main__":
    fix_order_sql()