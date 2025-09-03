-- 用户服务测试数据
-- 创建测试用户数据，包含不同角色和真实的用户信息

-- 清空现有数据
DELETE FROM users;

-- 重置自增ID
ALTER TABLE users AUTO_INCREMENT = 1;

-- 插入测试用户数据
INSERT INTO users (id, mobile, email, password, nick_name, user_name, birthday, gender, avatar, role, created_at, updated_at, is_deleted) VALUES
-- 管理员用户
(1, '13800000001', 'admin@joyshop.com', 'e10adc3949ba59abbe56e057f20f883e', '系统管理员', 'admin', '1990-01-01', 'male', 'https://api.dicebear.com/7.x/avataaars/svg?seed=admin', 2, NOW(), NOW(), false),

-- 普通用户 - 活跃用户
(2, '13800000002', 'zhangsan@example.com', 'e10adc3949ba59abbe56e057f20f883e', '张三', 'zhangsan', '1995-03-15', 'male', 'https://api.dicebear.com/7.x/avataaars/svg?seed=zhangsan', 1, NOW(), NOW(), false),
(3, '13800000003', 'lisi@example.com', 'e10adc3949ba59abbe56e057f20f883e', '李四', 'lisi', '1992-07-20', 'female', 'https://api.dicebear.com/7.x/avataaars/svg?seed=lisi', 1, NOW(), NOW(), false),
(4, '13800000004', 'wangwu@example.com', 'e10adc3949ba59abbe56e057f20f883e', '王五', 'wangwu', '1988-11-10', 'male', 'https://api.dicebear.com/7.x/avataaars/svg?seed=wangwu', 1, NOW(), NOW(), false),
(5, '13800000005', 'zhaoliu@example.com', 'e10adc3949ba59abbe56e057f20f883e', '赵六', 'zhaoliu', '1993-05-25', 'female', 'https://api.dicebear.com/7.x/avataaars/svg?seed=zhaoliu', 1, NOW(), NOW(), false),

-- 普通用户 - 中等活跃度用户
(6, '13800000006', 'sunqi@example.com', 'e10adc3949ba59abbe56e057f20f883e', '孙七', 'sunqi', '1991-09-08', 'male', 'https://api.dicebear.com/7.x/avataaars/svg?seed=sunqi', 1, NOW(), NOW(), false),
(7, '13800000007', 'zhouba@example.com', 'e10adc3949ba59abbe56e057f20f883e', '周八', 'zhouba', '1994-12-03', 'female', 'https://api.dicebear.com/7.x/avataaars/svg?seed=zhouba', 1, NOW(), NOW(), false),
(8, '13800000008', 'wujiu@example.com', 'e10adc3949ba59abbe56e057f20f883e', '吴九', 'wujiu', '1987-04-17', 'male', 'https://api.dicebear.com/7.x/avataaars/svg?seed=wujiu', 1, NOW(), NOW(), false),
(9, '13800000009', 'zhengshi@example.com', 'e10adc3949ba59abbe56e057f20f883e', '郑十', 'zhengshi', '1996-08-22', 'female', 'https://api.dicebear.com/7.x/avataaars/svg?seed=zhengshi', 1, NOW(), NOW(), false),

-- 新用户 - 低活跃度
(10, '13800000010', 'chenyi@example.com', 'e10adc3949ba59abbe56e057f20f883e', '陈一', 'chenyi', '1989-02-14', 'male', 'https://api.dicebear.com/7.x/avataaars/svg?seed=chenyi', 1, NOW(), NOW(), false),
(11, '13800000011', 'liner@example.com', 'e10adc3949ba59abbe56e057f20f883e', '李二', 'liner', '1993-10-30', 'female', 'https://api.dicebear.com/7.x/avataaars/svg?seed=liner', 1, NOW(), NOW(), false),
(12, '13800000012', 'wangsan@example.com', 'e10adc3949ba59abbe56e057f20f883e', '王三', 'wangsan', '1990-06-18', 'male', 'https://api.dicebear.com/7.x/avataaars/svg?seed=wangsan', 1, NOW(), NOW(), false),

-- VIP用户 - 高价值用户
(13, '13800000013', 'liuvip@example.com', 'e10adc3949ba59abbe56e057f20f883e', 'VIP刘总', 'liuvip', '1985-03-08', 'male', 'https://api.dicebear.com/7.x/avataaars/svg?seed=liuvip', 1, NOW(), NOW(), false),
(14, '13800000014', 'zhangvip@example.com', 'e10adc3949ba59abbe56e057f20f883e', 'VIP张总', 'zhangvip', '1982-11-15', 'female', 'https://api.dicebear.com/7.x/avataaars/svg?seed=zhangvip', 1, NOW(), NOW(), false),
(15, '13800000015', 'testuser@example.com', 'e10adc3949ba59abbe56e057f20f883e', '测试用户', 'testuser', '1990-01-01', 'unknown', 'https://api.dicebear.com/7.x/avataaars/svg?seed=testuser', 1, NOW(), NOW(), false);

-- 说明：
-- 密码统一为 123456 的MD5值：e10adc3949ba59abbe56e057f20f883e
-- 角色：1=普通用户，2=管理员
-- 性别：male=男，female=女，unknown=未知
-- 头像使用 dicebear API 生成的随机头像
-- 用户分为不同活跃度等级，便于测试不同场景