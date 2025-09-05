-- 修复后的用户数据
-- 设置字符集和SQL模式
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- 用户服务测试数据
-- 创建测试用户数据，包含不同角色和真实的用户信息

-- 清空现有数据
DELETE FROM user;

-- 重置自增ID
ALTER TABLE user AUTO_INCREMENT = 1;

-- 插入测试用户数据
INSERT INTO user (id, mobile, email, password, nick_name, user_name, birthday, gender, avatar, role, created_at, updated_at) VALUES
-- 管理员用户
(1, '13800000001', 'admin@joyshop.com', 'e10adc3949ba59abbe56e057f20f883e', '系统管理员', 'admin', '1990-01-01', 'male', 'placeholder_avatar_admin', 2, NOW(), NOW()),

-- 普通用户 - 活跃用户
(2, '13800000002', 'zhangsan@example.com', 'e10adc3949ba59abbe56e057f20f883e', '张三', 'zhangsan', '1995-03-15', 'male', 'placeholder_avatar_zhangsan', 1, NOW(), NOW()),
(3, '13800000003', 'lisi@example.com', 'e10adc3949ba59abbe56e057f20f883e', '李四', 'lisi', '1992-07-20', 'female', 'placeholder_avatar_lisi', 1, NOW(), NOW()),
(4, '13800000004', 'wangwu@example.com', 'e10adc3949ba59abbe56e057f20f883e', '王五', 'wangwu', '1988-11-10', 'male', 'placeholder_avatar_wangwu', 1, NOW(), NOW()),
(5, '13800000005', 'zhaoliu@example.com', 'e10adc3949ba59abbe56e057f20f883e', '赵六', 'zhaoliu', '1993-05-25', 'female', 'placeholder_avatar_zhaoliu', 1, NOW(), NOW()),

-- 普通用户 - 中等活跃度用户
(6, '13800000006', 'sunqi@example.com', 'e10adc3949ba59abbe56e057f20f883e', '孙七', 'sunqi', '1991-09-08', 'male', 'placeholder_avatar_sunqi', 1, NOW(), NOW()),
(7, '13800000007', 'zhouba@example.com', 'e10adc3949ba59abbe56e057f20f883e', '周八', 'zhouba', '1994-12-03', 'female', 'placeholder_avatar_zhouba', 1, NOW(), NOW()),
(8, '13800000008', 'wujiu@example.com', 'e10adc3949ba59abbe56e057f20f883e', '吴九', 'wujiu', '1987-04-17', 'male', 'placeholder_avatar_wujiu', 1, NOW(), NOW()),
(9, '13800000009', 'zhengshi@example.com', 'e10adc3949ba59abbe56e057f20f883e', '郑十', 'zhengshi', '1996-08-22', 'female', 'placeholder_avatar_zhengshi', 1, NOW(), NOW()),

-- 新用户 - 低活跃度
(10, '13800000010', 'chenyi@example.com', 'e10adc3949ba59abbe56e057f20f883e', '陈一', 'chenyi', '1989-02-14', 'male', 'placeholder_avatar_chenyi', 1, NOW(), NOW()),
(11, '13800000011', 'liner@example.com', 'e10adc3949ba59abbe56e057f20f883e', '李二', 'liner', '1993-10-30', 'female', 'placeholder_avatar_liner', 1, NOW(), NOW()),
(12, '13800000012', 'wangsan@example.com', 'e10adc3949ba59abbe56e057f20f883e', '王三', 'wangsan', '1990-06-18', 'male', 'placeholder_avatar_wangsan', 1, NOW(), NOW()),

-- VIP用户 - 高价值用户
(13, '13800000013', 'liuvip@example.com', 'e10adc3949ba59abbe56e057f20f883e', 'VIP刘总', 'liuvip', '1985-03-08', 'male', 'placeholder_avatar_liuvip', 1, NOW(), NOW()),
(14, '13800000014', 'zhangvip@example.com', 'e10adc3949ba59abbe56e057f20f883e', 'VIP张总', 'zhangvip', '1982-11-15', 'female', 'placeholder_avatar_zhangvip', 1, NOW(), NOW()),
(15, '13800000015', 'testuser@example.com', 'e10adc3949ba59abbe56e057f20f883e', '测试用户', 'testuser', '1990-01-01', 'unknown', 'placeholder_avatar_testuser', 1, NOW(), NOW());

-- 说明：
-- 密码统一为 123456 的MD5值：e10adc3949ba59abbe56e057f20f883e
-- 角色：1=普通用户，2=管理员
-- 性别：male=男，female=女，unknown=未知
-- 头像使用 dicebear API 生成的随机头像
-- 用户分为不同活跃度等级，便于测试不同场景
-- 恢复外键检查
SET FOREIGN_KEY_CHECKS = 1;
