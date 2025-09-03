-- JoyShop 微服务数据库创建脚本
-- 为微服务架构创建独立的数据库

-- ==============================================
-- 创建微服务独立数据库
-- ==============================================

-- 设置字符集
SET NAMES utf8mb4;

SELECT '🗄️ 开始创建微服务数据库...' AS 'Status';

-- 1. 创建用户服务数据库
DROP DATABASE IF EXISTS joyshop_user;
CREATE DATABASE joyshop_user 
CHARACTER SET utf8mb4 
COLLATE utf8mb4_unicode_ci 
COMMENT '用户服务数据库';

SELECT '✅ 创建用户服务数据库: joyshop_user' AS 'Database Created';

-- 2. 创建商品服务数据库
DROP DATABASE IF EXISTS joyshop_goods;
CREATE DATABASE joyshop_goods 
CHARACTER SET utf8mb4 
COLLATE utf8mb4_unicode_ci 
COMMENT '商品服务数据库';

SELECT '✅ 创建商品服务数据库: joyshop_goods' AS 'Database Created';

-- 3. 创建库存服务数据库
DROP DATABASE IF EXISTS joyshop_inventory;
CREATE DATABASE joyshop_inventory 
CHARACTER SET utf8mb4 
COLLATE utf8mb4_unicode_ci 
COMMENT '库存服务数据库';

SELECT '✅ 创建库存服务数据库: joyshop_inventory' AS 'Database Created';

-- 4. 创建订单服务数据库
DROP DATABASE IF EXISTS joyshop_order;
CREATE DATABASE joyshop_order 
CHARACTER SET utf8mb4 
COLLATE utf8mb4_unicode_ci 
COMMENT '订单服务数据库';

SELECT '✅ 创建订单服务数据库: joyshop_order' AS 'Database Created';

-- 5. （可选）创建单体架构数据库（用于开发测试）
DROP DATABASE IF EXISTS joyshop_all;
CREATE DATABASE joyshop_all 
CHARACTER SET utf8mb4 
COLLATE utf8mb4_unicode_ci 
COMMENT '单体架构数据库（开发测试用）';

SELECT '✅ 创建单体架构数据库: joyshop_all' AS 'Database Created';

SELECT '' AS '', '============================================' AS '';
SELECT '🎉 数据库创建完成！' AS 'Status';
SELECT '============================================' AS '';

-- 显示创建的数据库
SELECT '📊 已创建的数据库：' AS '';
SHOW DATABASES LIKE 'joyshop_%';

SELECT '' AS '', '📝 数据库用途说明：' AS '';
SELECT 'joyshop_user     - 用户服务数据库' AS 'Database Purpose';
SELECT 'joyshop_goods    - 商品服务数据库' AS 'Database Purpose';
SELECT 'joyshop_inventory- 库存服务数据库' AS 'Database Purpose';
SELECT 'joyshop_order    - 订单服务数据库' AS 'Database Purpose';
SELECT 'joyshop_all      - 单体架构数据库（可选）' AS 'Database Purpose';

SELECT '' AS '', '🚀 下一步操作：' AS '';
SELECT '微服务模式: mysql> source sql/init_microservices.sql;' AS 'Next Steps';
SELECT '单体模式:   mysql> use joyshop_all; source sql/init_all.sql;' AS 'Next Steps';

SELECT '============================================' AS '';