-- 清理所有测试脏数据脚本
-- 执行前请确认数据库备份

-- =============================================
-- 清理 joyshop_goods 数据库
-- =============================================
USE joyshop_goods;
SET FOREIGN_KEY_CHECKS = 0;

-- 清空商品相关表
TRUNCATE TABLE goods_category;
TRUNCATE TABLE category_brand;
TRUNCATE TABLE goods;
TRUNCATE TABLE category;
TRUNCATE TABLE brand;
TRUNCATE TABLE banner;

SET FOREIGN_KEY_CHECKS = 1;

SELECT '✅ joyshop_goods 数据已清理' AS status;

-- =============================================
-- 清理 joyshop_inventory 数据库
-- =============================================
USE joyshop_inventory;
SET FOREIGN_KEY_CHECKS = 0;

-- 清空库存表
TRUNCATE TABLE inventory;

SET FOREIGN_KEY_CHECKS = 1;

SELECT '✅ joyshop_inventory 数据已清理' AS status;

-- =============================================
-- 清理 joyshop_order 数据库
-- =============================================
USE joyshop_order;
SET FOREIGN_KEY_CHECKS = 0;

-- 清空订单相关表
TRUNCATE TABLE order_goods;
TRUNCATE TABLE order_info;
TRUNCATE TABLE shopping_cart;

SET FOREIGN_KEY_CHECKS = 1;

SELECT '✅ joyshop_order 数据已清理' AS status;

-- =============================================
-- 清理 joyshop_user 数据库
-- =============================================
USE joyshop_user;
SET FOREIGN_KEY_CHECKS = 0;

-- 清空用户表
TRUNCATE TABLE users;

SET FOREIGN_KEY_CHECKS = 1;

SELECT '✅ joyshop_user 数据已清理' AS status;

-- =============================================
-- 显示清理结果
-- =============================================
SELECT '============================================' AS '';
SELECT '🧹 所有测试数据清理完成！' AS 'Clean Status';
SELECT '============================================' AS '';
SELECT '准备导入新数据...' AS 'Next Step';