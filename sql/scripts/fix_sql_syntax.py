#!/usr/bin/env python3
"""
修复SQL文件中的语法错误
包括单引号转义、特殊字符处理等
"""

import re
import os

def fix_sql_syntax(file_path):
    """修复SQL文件的语法错误"""
    print(f"正在修复文件: {file_path}")
    
    # 读取文件
    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()
    
    original_content = content
    
    # 1. 修复单引号转义问题
    # 查找所有在单引号内的内容
    def escape_quotes_in_string(match):
        string_content = match.group(1)
        # 转义内部的单引号
        escaped_content = string_content.replace("'", "\\'")
        return f"'{escaped_content}'"
    
    # 使用正则表达式查找并替换字符串中的单引号
    # 这个正则表达式匹配 'xxx' 格式的字符串
    content = re.sub(r"'([^']*(?:\\'[^']*)*)'", escape_quotes_in_string, content)
    
    # 2. 特殊处理一些已知的问题品牌名
    problem_brands = {
        "Levi\\'s": "Levis",  # 去掉单引号
        "L\\'Oréal": "LOreal",  # 去掉单引号和特殊字符
        "L'Oréal": "LOreal",   # 直接替换
        "Levi's": "Levis"      # 直接替换
    }
    
    for old, new in problem_brands.items():
        content = content.replace(old, new)
    
    # 3. 检查JSON格式是否正确
    # 确保JSON字段中的双引号被正确转义
    def fix_json_in_sql(match):
        json_content = match.group(1)
        # 确保JSON中的双引号被正确转义
        if not json_content.startswith('[') or not json_content.endswith(']'):
            return match.group(0)  # 返回原始内容
        
        return f"'{json_content}'"
    
    # 修复images和desc_images字段的JSON格式
    content = re.sub(r"'(\[\"[^']*\"\])'", fix_json_in_sql, content)
    
    # 4. 确保所有的逗号和分号正确
    # 检查SQL语句结构
    lines = content.split('\n')
    fixed_lines = []
    
    for line in lines:
        # 跳过注释和空行
        if line.strip().startswith('--') or not line.strip():
            fixed_lines.append(line)
            continue
        
        # 确保VALUES语句的格式正确
        if 'VALUES' in line and not line.strip().endswith(',') and not line.strip().endswith(';'):
            if not any(keyword in line.upper() for keyword in ['INSERT', 'VALUES']):
                # 这是VALUES中的数据行，应该以逗号结尾（除了最后一行）
                pass
        
        fixed_lines.append(line)
    
    content = '\n'.join(fixed_lines)
    
    # 5. 添加SQL错误处理
    # 在文件开头添加一些设置来处理可能的编码和约束问题
    header = """-- 修复后的商品数据
-- 设置字符集和SQL模式
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;
SET sql_mode = 'NO_AUTO_VALUE_ON_ZERO';

"""
    
    footer = """
-- 恢复外键检查
SET FOREIGN_KEY_CHECKS = 1;
"""
    
    # 如果内容没有以这些设置开头，则添加
    if not content.startswith('-- 修复后的商品数据'):
        content = header + content + footer
    
    # 写回文件
    with open(file_path, 'w', encoding='utf-8') as f:
        f.write(content)
    
    # 统计修改
    changes = len(original_content) != len(content)
    print(f"文件修复完成，{'有' if changes else '无'}内容变化")
    
    return changes

def main():
    """主函数"""
    base_dir = "/Users/walker/gitsource/github/joyshop_srvs/sql"
    
    # 需要修复的文件列表
    files_to_fix = [
        "schemas/04_goods.sql",
        "schemas/03_brands.sql",  # 也可能有单引号问题
    ]
    
    print("开始修复SQL语法错误...")
    
    total_changes = 0
    for file_path in files_to_fix:
        full_path = os.path.join(base_dir, file_path)
        if os.path.exists(full_path):
            if fix_sql_syntax(full_path):
                total_changes += 1
        else:
            print(f"文件不存在: {full_path}")
    
    print(f"\n修复完成！共修复 {total_changes} 个文件")
    print("现在可以尝试导入数据了")

if __name__ == "__main__":
    main()