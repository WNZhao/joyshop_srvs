#!/usr/bin/env python3
"""
修复库存数据，只保留实际表中存在的字段
根据 inventory 表结构：id, goods_id, stock, version, created_at, updated_at, deleted_at, is_deleted
"""

import re

def fix_inventory_sql():
    """修复库存SQL文件"""
    file_path = "/Users/walker/gitsource/github/joyshop_srvs/sql/schemas/05_inventory.sql"
    
    print(f"正在修复库存数据文件: {file_path}")
    
    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # 修复INSERT语句，只保留需要的字段
    # 原来的格式: INSERT INTO inventory (goods_id, stock, sold, frozen, created_at, updated_at)
    # 修改为: INSERT INTO inventory (goods_id, stock, version, created_at, updated_at)
    
    # 替换INSERT语句的字段列表
    content = content.replace(
        "INSERT INTO inventory (goods_id, stock, sold, frozen, created_at, updated_at)",
        "INSERT INTO inventory (goods_id, stock, version, created_at, updated_at)"
    )
    
    # 修复VALUES中的数据
    # 原格式: (goods_id, stock_value, sold_value, frozen_value, NOW(), NOW())
    # 新格式: (goods_id, stock_value, version, NOW(), NOW())
    # 其中version设为0，去掉sold和frozen字段
    
    lines = content.split('\n')
    fixed_lines = []
    
    for line in lines:
        if line.strip().startswith('(') and 'NOW()' in line:
            # 这是数据行，需要修复
            # 使用正则表达式提取字段
            match = re.match(r'\((\d+),\s*(\d+),\s*\d+,\s*\d+,\s*(NOW\(\)),\s*(NOW\(\))\)(.*)', line.strip())
            if match:
                goods_id = match.group(1)
                stock = match.group(2)
                version = "0"  # 版本号设为0
                created_at = match.group(3)
                updated_at = match.group(4)
                suffix = match.group(5)  # 可能是逗号或分号
                
                # 重构行
                new_line = f"({goods_id}, {stock}, {version}, {created_at}, {updated_at}){suffix}"
                fixed_lines.append(new_line)
            else:
                fixed_lines.append(line)
        else:
            fixed_lines.append(line)
    
    content = '\n'.join(fixed_lines)
    
    # 添加SQL设置
    header = """-- 修复后的库存数据
-- 设置字符集和SQL模式
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

"""
    
    footer = """
-- 恢复外键检查
SET FOREIGN_KEY_CHECKS = 1;
"""
    
    if not content.startswith('-- 修复后的库存数据'):
        content = header + content + footer
    
    with open(file_path, 'w', encoding='utf-8') as f:
        f.write(content)
    
    print("库存数据修复完成")

if __name__ == "__main__":
    fix_inventory_sql()