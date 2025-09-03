-- 商品分类测试数据
-- 创建3级分类结构，模拟真实电商分类体系

-- 清空现有数据
DELETE FROM goods_category;
DELETE FROM category;

-- 重置自增ID
ALTER TABLE category AUTO_INCREMENT = 1;

-- 一级分类（Level 1）
INSERT INTO category (id, name, parent_id, level, sort, is_tab, created_at, updated_at) VALUES
(1, '电子数码', NULL, 1, 1, true, NOW(), NOW()),
(2, '服装鞋包', NULL, 1, 2, true, NOW(), NOW()),
(3, '家居生活', NULL, 1, 3, true, NOW(), NOW()),
(4, '图书文教', NULL, 1, 4, true, NOW(), NOW()),
(5, '运动户外', NULL, 1, 5, true, NOW(), NOW()),
(6, '美妆个护', NULL, 1, 6, true, NOW(), NOW()),
(7, '食品饮料', NULL, 1, 7, false, NOW(), NOW()),
(8, '母婴玩具', NULL, 1, 8, false, NOW(), NOW());

-- 二级分类（Level 2）- 电子数码
INSERT INTO category (id, name, parent_id, level, sort, is_tab, created_at, updated_at) VALUES
-- 电子数码子分类
(11, '手机通讯', 1, 2, 1, false, NOW(), NOW()),
(12, '电脑办公', 1, 2, 2, false, NOW(), NOW()),
(13, '家用电器', 1, 2, 3, false, NOW(), NOW()),
(14, '摄影摄像', 1, 2, 4, false, NOW(), NOW()),
(15, '智能设备', 1, 2, 5, false, NOW(), NOW()),

-- 服装鞋包子分类
(21, '男装', 2, 2, 1, false, NOW(), NOW()),
(22, '女装', 2, 2, 2, false, NOW(), NOW()),
(23, '鞋类', 2, 2, 3, false, NOW(), NOW()),
(24, '箱包配饰', 2, 2, 4, false, NOW(), NOW()),

-- 家居生活子分类  
(31, '家具', 3, 2, 1, false, NOW(), NOW()),
(32, '家装建材', 3, 2, 2, false, NOW(), NOW()),
(33, '家纺', 3, 2, 3, false, NOW(), NOW()),
(34, '厨具', 3, 2, 4, false, NOW(), NOW()),

-- 图书文教子分类
(41, '图书', 4, 2, 1, false, NOW(), NOW()),
(42, '文具用品', 4, 2, 2, false, NOW(), NOW()),
(43, '教育培训', 4, 2, 3, false, NOW(), NOW()),

-- 运动户外子分类
(51, '运动服装', 5, 2, 1, false, NOW(), NOW()),
(52, '运动鞋', 5, 2, 2, false, NOW(), NOW()),
(53, '运动器材', 5, 2, 3, false, NOW(), NOW()),
(54, '户外用品', 5, 2, 4, false, NOW(), NOW()),

-- 美妆个护子分类
(61, '护肤', 6, 2, 1, false, NOW(), NOW()),
(62, '彩妆', 6, 2, 2, false, NOW(), NOW()),
(63, '个人护理', 6, 2, 3, false, NOW(), NOW());

-- 三级分类（Level 3）- 细分类别
INSERT INTO category (id, name, parent_id, level, sort, is_tab, created_at, updated_at) VALUES
-- 手机通讯细分
(111, 'iPhone', 11, 3, 1, false, NOW(), NOW()),
(112, 'Android手机', 11, 3, 2, false, NOW(), NOW()),
(113, '手机配件', 11, 3, 3, false, NOW(), NOW()),

-- 电脑办公细分
(121, '笔记本电脑', 12, 3, 1, false, NOW(), NOW()),
(122, '台式机', 12, 3, 2, false, NOW(), NOW()),
(123, '平板电脑', 12, 3, 3, false, NOW(), NOW()),
(124, '电脑配件', 12, 3, 4, false, NOW(), NOW()),

-- 家用电器细分
(131, '大家电', 13, 3, 1, false, NOW(), NOW()),
(132, '小家电', 13, 3, 2, false, NOW(), NOW()),

-- 摄影摄像细分  
(141, '单反相机', 14, 3, 1, false, NOW(), NOW()),
(142, '微单相机', 14, 3, 2, false, NOW(), NOW()),
(143, '运动相机', 14, 3, 3, false, NOW(), NOW()),

-- 男装细分
(211, 'T恤', 21, 3, 1, false, NOW(), NOW()),
(212, '衬衫', 21, 3, 2, false, NOW(), NOW()),
(213, '牛仔裤', 21, 3, 3, false, NOW(), NOW()),
(214, '休闲裤', 21, 3, 4, false, NOW(), NOW()),

-- 女装细分
(221, '连衣裙', 22, 3, 1, false, NOW(), NOW()),
(222, '针织衫', 22, 3, 2, false, NOW(), NOW()),
(223, '半身裙', 22, 3, 3, false, NOW(), NOW()),

-- 鞋类细分
(231, '运动鞋', 23, 3, 1, false, NOW(), NOW()),
(232, '休闲鞋', 23, 3, 2, false, NOW(), NOW()),
(233, '正装鞋', 23, 3, 3, false, NOW(), NOW()),

-- 运动服装细分
(511, '运动T恤', 51, 3, 1, false, NOW(), NOW()),
(512, '运动裤', 51, 3, 2, false, NOW(), NOW()),
(513, '运动套装', 51, 3, 3, false, NOW(), NOW()),

-- 运动鞋细分  
(521, '跑步鞋', 52, 3, 1, false, NOW(), NOW()),
(522, '篮球鞋', 52, 3, 2, false, NOW(), NOW()),
(523, '足球鞋', 52, 3, 3, false, NOW(), NOW());

-- 说明：
-- Level 1: 一级分类（顶级分类）
-- Level 2: 二级分类（主要分类）  
-- Level 3: 三级分类（详细分类）
-- is_tab: true 表示在导航栏显示
-- sort: 排序字段，数字越小越靠前