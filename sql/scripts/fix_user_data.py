#!/usr/bin/env python3
"""
修复用户数据文件，使其与user.go模型完全匹配
实际表结构：id, created_at, updated_at, deleted_at, is_deleted, mobile, email, password, 
nick_name, user_name, birthday, gender, avatar, role
"""

import re

def fix_user_sql():
    """修复用户SQL文件"""
    file_path = "/Users/walker/gitsource/github/joyshop_srvs/sql/schemas/01_users.sql"
    
    print(f"正在修复用户数据文件: {file_path}")
    
    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # 修复INSERT语句，使用正确的字段顺序
    # 实际表结构字段顺序：id, mobile, email, password, nick_name, user_name, birthday, gender, avatar, role, created_at, updated_at
    old_insert = "INSERT INTO user (id, mobile, email, password, nick_name, user_name, birthday, gender, avatar, role, created_at, updated_at, is_deleted)"
    new_insert = "INSERT INTO user (id, mobile, email, password, nick_name, user_name, birthday, gender, avatar, role, created_at, updated_at)"
    
    content = content.replace(old_insert, new_insert)
    
    # 修复数据行，移除is_deleted字段（最后一个字段）
    lines = content.split('\n')
    fixed_lines = []
    
    for line in lines:
        if line.strip().startswith('(') and 'NOW()' in line:
            # 这是数据行，需要修复 - 移除最后的 is_deleted 字段
            # 格式: (id, mobile, email, password, nick_name, user_name, birthday, gender, avatar, role, NOW(), NOW(), false)
            # 改为: (id, mobile, email, password, nick_name, user_name, birthday, gender, avatar, role, NOW(), NOW())
            
            # 使用正则表达式移除最后的 , false 或 , true
            line = re.sub(r',\s*(false|true)\s*\)', ')', line)
            fixed_lines.append(line)
        else:
            fixed_lines.append(line)
    
    content = '\n'.join(fixed_lines)
    
    # 添加SQL设置
    header = """-- 修复后的用户数据
-- 设置字符集和SQL模式
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

"""
    
    footer = """
-- 恢复外键检查
SET FOREIGN_KEY_CHECKS = 1;
"""
    
    if not content.startswith('-- 修复后的用户数据'):
        content = header + content + footer
    
    with open(file_path, 'w', encoding='utf-8') as f:
        f.write(content)
    
    print("用户数据修复完成")
    print("修复内容：")
    print("- 移除 is_deleted 字段")
    print("- 调整字段顺序匹配实际表结构")
    print("- 添加SQL设置")

if __name__ == "__main__":
    fix_user_sql()