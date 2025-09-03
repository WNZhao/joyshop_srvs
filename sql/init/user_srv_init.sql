-- 用户服务初始化脚本
-- 仅初始化用户服务相关的数据

-- 设置字符集
SET NAMES utf8mb4;

-- 禁用外键检查
SET FOREIGN_KEY_CHECKS = 0;

-- 执行用户数据初始化
SOURCE 01_users.sql;

-- 启用外键检查  
SET FOREIGN_KEY_CHECKS = 1;

-- 显示结果
SELECT '========================================' AS '';
SELECT '✅ 用户服务测试数据初始化完成！' AS 'Status';
SELECT '========================================' AS '';

-- 统计信息
SELECT 'users' AS table_name, COUNT(*) AS record_count FROM users;

-- 测试用户说明
SELECT '👤 测试用户说明：' AS '';
SELECT 'ID=1: admin (管理员) - admin@joyshop.com' AS 'Test Users';
SELECT 'ID=2: zhangsan (张三) - zhangsan@example.com' AS 'Test Users';  
SELECT 'ID=3: lisi (李四) - lisi@example.com' AS 'Test Users';
SELECT 'ID=4: wangwu (王五) - wangwu@example.com' AS 'Test Users';
SELECT 'ID=13-14: VIP用户 - 高价值客户' AS 'Test Users';
SELECT '密码统一为：123456' AS 'Password';
SELECT '========================================' AS '';