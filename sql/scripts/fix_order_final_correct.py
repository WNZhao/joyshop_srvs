#!/usr/bin/env python3
"""
最终正确修复订单数据
分析数据格式，确保字段数量完全匹配
"""

def fix_order_final():
    """最终修复订单数据"""
    file_path = "/Users/walker/gitsource/github/joyshop_srvs/sql/schemas/07_orders.sql"
    
    print(f"最终修复订单数据: {file_path}")
    
    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # 重新生成整个文件，使用简化的方法
    new_content = """-- 根据order.go模型修复的订单数据
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- 订单服务测试数据 (修复版)
-- 清空现有数据
DELETE FROM order_goods;
DELETE FROM order_info;

-- 重置自增ID
ALTER TABLE order_info AUTO_INCREMENT = 1;
ALTER TABLE order_goods AUTO_INCREMENT = 1;

-- 插入订单信息 - 使用简化字段
INSERT INTO order_info (user, order_sn, order_mount, address, signer_name, singer_mobile, post, pay_type, status, trade_no, pay_time) VALUES
(13, 'JOY20250731000001', 1517.85, '测试地址13-1', '收货人13', '13800138945', '585492', 'alipay', 'TRADE_SUCCESS', 'trade_1_698774', NOW()),
(11, 'JOY20250620000002', 30466.37, '测试地址11-2', '收货人11', '13800114233', '881125', 'wechat', 'TRADE_FINISHED', 'trade_2_112677', NOW()),
(13, 'JOY20250702000003', 3195.54, '测试地址13-3', '收货人13', '13800135563', '860274', 'alipay', 'PAYING', '', NOW()),
(15, 'JOY20250828000004', 5784.79, '测试地址15-4', '收货人15', '13800152744', '672488', 'wechat', 'TRADE_FINISHED', 'trade_4_201824', NOW()),
(9, 'JOY20250829000005', 3027.66, '测试地址9-5', '收货人9', '13800091707', '834715', 'alipay', 'TRADE_FINISHED', 'trade_5_445693', NOW());

-- 插入订单商品信息 - 使用正确字段名 (order不是order_id)
INSERT INTO order_goods (`order`, goods, goods_name, goods_image, goods_price, nums) VALUES
(1, 1001, 'iPhone 15 Pro Max 128GB 深空黑色', 'placeholder_image_1001_1.jpg', 9999.00, 1),
(1, 2001, 'Nike Air Jordan 1 Retro High OG 黑红配色', 'placeholder_image_2001_1.jpg', 1299.00, 1),
(2, 3001, 'MacBook Pro 14英寸 M3', 'placeholder_image_3001_1.jpg', 15999.00, 2),
(3, 1002, 'HUAWEI Mate 60 Pro 256GB 雅川青', 'placeholder_image_1002_1.jpg', 6999.00, 1),
(4, 2002, 'Adidas Ultraboost 22 黑白配色', 'placeholder_image_2002_1.jpg', 1599.00, 2),
(5, 1003, 'Xiaomi 14 Ultra 12GB+256GB 黑色', 'placeholder_image_1003_1.jpg', 5999.00, 1);

-- 恢复外键检查
SET FOREIGN_KEY_CHECKS = 1;
"""
    
    # 写入修复后的内容
    with open(file_path, 'w', encoding='utf-8') as f:
        f.write(new_content)
    
    print("订单数据最终修复完成！")
    print("使用简化的测试数据，确保字段完全匹配数据库表结构")

if __name__ == "__main__":
    fix_order_final()