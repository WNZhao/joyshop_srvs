-- 品牌测试数据
-- 创建各个分类的知名品牌数据

-- 清空现有数据
DELETE FROM category_brand;
DELETE FROM brand;

-- 重置自增ID
ALTER TABLE brand AUTO_INCREMENT = 1;

-- 插入品牌数据
INSERT INTO brand (id, name, logo, `desc`, created_at, updated_at) VALUES
-- 电子数码品牌
(1, 'Apple', 'https://cdn.jsdelivr.net/gh/joyshop/brand-logos/apple.png', '苹果公司，创新科技的引领者', NOW(), NOW()),
(2, 'HUAWEI', 'https://cdn.jsdelivr.net/gh/joyshop/brand-logos/huawei.png', '华为技术有限公司，全球领先的ICT基础设施和智能终端提供商', NOW(), NOW()),
(3, 'Xiaomi', 'https://cdn.jsdelivr.net/gh/joyshop/brand-logos/xiaomi.png', '小米集团，以手机、智能硬件和IoT平台为核心的互联网公司', NOW(), NOW()),
(4, 'Samsung', 'https://cdn.jsdelivr.net/gh/joyshop/brand-logos/samsung.png', '三星电子，全球科技企业领导者', NOW(), NOW()),
(5, 'Lenovo', 'https://cdn.jsdelivr.net/gh/joyshop/brand-logos/lenovo.png', '联想集团，全球化的科技公司', NOW(), NOW()),
(6, 'Dell', 'https://cdn.jsdelivr.net/gh/joyshop/brand-logos/dell.png', '戴尔科技集团，全球领先的技术解决方案提供商', NOW(), NOW()),
(7, 'Sony', 'https://cdn.jsdelivr.net/gh/joyshop/brand-logos/sony.png', '索尼公司，娱乐科技的创新者', NOW(), NOW()),
(8, 'Canon', 'https://cdn.jsdelivr.net/gh/joyshop/brand-logos/canon.png', '佳能公司，影像技术的领导者', NOW(), NOW()),

-- 服装品牌
(11, 'Nike', 'https://cdn.jsdelivr.net/gh/joyshop/brand-logos/nike.png', '耐克公司，全球知名运动品牌', NOW(), NOW()),
(12, 'Adidas', 'https://cdn.jsdelivr.net/gh/joyshop/brand-logos/adidas.png', '阿迪达斯，德国运动用品制造商', NOW(), NOW()),
(13, 'UNIQLO', 'https://cdn.jsdelivr.net/gh/joyshop/brand-logos/uniqlo.png', '优衣库，日本服装品牌', NOW(), NOW()),
(14, 'ZARA', 'https://cdn.jsdelivr.net/gh/joyshop/brand-logos/zara.png', 'ZARA，西班牙时装连锁品牌', NOW(), NOW()),
(15, 'H&M', 'https://cdn.jsdelivr.net/gh/joyshop/brand-logos/hm.png', 'H&M，瑞典快时尚品牌', NOW(), NOW()),
(16, 'Levi\'s', 'https://cdn.jsdelivr.net/gh/joyshop/brand-logos/levis.png', 'Levi\'s，美国牛仔裤品牌', NOW(), NOW()),

-- 家居生活品牌
(21, 'IKEA', 'https://cdn.jsdelivr.net/gh/joyshop/brand-logos/ikea.png', '宜家家居，瑞典家具零售企业', NOW(), NOW()),
(22, '全友家居', 'https://cdn.jsdelivr.net/gh/joyshop/brand-logos/quanyou.png', '全友家私，中国家具行业领军企业', NOW(), NOW()),
(23, '美的', 'https://cdn.jsdelivr.net/gh/joyshop/brand-logos/midea.png', '美的集团，中国家电制造企业', NOW(), NOW()),
(24, '海尔', 'https://cdn.jsdelivr.net/gh/joyshop/brand-logos/haier.png', '海尔集团，中国家电企业', NOW(), NOW()),

-- 图书出版品牌
(31, '人民邮电出版社', 'https://cdn.jsdelivr.net/gh/joyshop/brand-logos/ptpress.png', '人民邮电出版社，专业IT图书出版社', NOW(), NOW()),
(32, '机械工业出版社', 'https://cdn.jsdelivr.net/gh/joyshop/brand-logos/cmpedu.png', '机械工业出版社，综合性出版社', NOW(), NOW()),
(33, '晨光文具', 'https://cdn.jsdelivr.net/gh/joyshop/brand-logos/chenguang.png', '上海晨光文具股份有限公司', NOW(), NOW()),

-- 运动户外品牌
(41, 'Under Armour', 'https://cdn.jsdelivr.net/gh/joyshop/brand-logos/underarmour.png', '安德玛，美国运动品牌', NOW(), NOW()),
(42, 'New Balance', 'https://cdn.jsdelivr.net/gh/joyshop/brand-logos/newbalance.png', '新百伦，美国运动品牌', NOW(), NOW()),
(43, 'The North Face', 'https://cdn.jsdelivr.net/gh/joyshop/brand-logos/tnf.png', '北面，美国户外品牌', NOW(), NOW()),

-- 美妆个护品牌
(51, 'L\'Oréal', 'https://cdn.jsdelivr.net/gh/joyshop/brand-logos/loreal.png', '欧莱雅集团，全球化妆品领导者', NOW(), NOW()),
(52, 'Estée Lauder', 'https://cdn.jsdelivr.net/gh/joyshop/brand-logos/esteelauder.png', '雅诗兰黛，美国化妆品品牌', NOW(), NOW()),
(53, 'SK-II', 'https://cdn.jsdelivr.net/gh/joyshop/brand-logos/skii.png', 'SK-II，日本护肤品牌', NOW(), NOW()),

-- 自有品牌
(61, 'JoyShop', 'https://cdn.jsdelivr.net/gh/joyshop/brand-logos/joyshop.png', 'JoyShop自有品牌，精选优质商品', NOW(), NOW());

-- 插入分类品牌关联数据
INSERT INTO category_brand (category_id, brand_id, created_at, updated_at) VALUES
-- 电子数码品牌关联
-- 手机通讯品牌
(11, 1, NOW(), NOW()), -- Apple - 手机通讯
(11, 2, NOW(), NOW()), -- HUAWEI - 手机通讯  
(11, 3, NOW(), NOW()), -- Xiaomi - 手机通讯
(11, 4, NOW(), NOW()), -- Samsung - 手机通讯

-- 电脑办公品牌
(12, 1, NOW(), NOW()), -- Apple - 电脑办公
(12, 5, NOW(), NOW()), -- Lenovo - 电脑办公
(12, 6, NOW(), NOW()), -- Dell - 电脑办公
(12, 2, NOW(), NOW()), -- HUAWEI - 电脑办公

-- 家用电器品牌
(13, 23, NOW(), NOW()), -- 美的 - 家用电器
(13, 24, NOW(), NOW()), -- 海尔 - 家用电器
(13, 3, NOW(), NOW()),  -- 小米 - 家用电器

-- 摄影摄像品牌
(14, 8, NOW(), NOW()), -- Canon - 摄影摄像
(14, 7, NOW(), NOW()), -- Sony - 摄影摄像

-- 服装品牌关联
-- 男装品牌
(21, 13, NOW(), NOW()), -- UNIQLO - 男装
(21, 14, NOW(), NOW()), -- ZARA - 男装  
(21, 15, NOW(), NOW()), -- H&M - 男装
(21, 16, NOW(), NOW()), -- Levi's - 男装

-- 女装品牌
(22, 13, NOW(), NOW()), -- UNIQLO - 女装
(22, 14, NOW(), NOW()), -- ZARA - 女装
(22, 15, NOW(), NOW()), -- H&M - 女装

-- 鞋类品牌
(23, 11, NOW(), NOW()), -- Nike - 鞋类
(23, 12, NOW(), NOW()), -- Adidas - 鞋类
(23, 42, NOW(), NOW()), -- New Balance - 鞋类

-- 家居生活品牌关联
(31, 21, NOW(), NOW()), -- IKEA - 家具
(31, 22, NOW(), NOW()), -- 全友家居 - 家具

-- 图书文教品牌关联
(41, 31, NOW(), NOW()), -- 人民邮电出版社 - 图书
(41, 32, NOW(), NOW()), -- 机械工业出版社 - 图书
(42, 33, NOW(), NOW()), -- 晨光文具 - 文具用品

-- 运动户外品牌关联
(51, 11, NOW(), NOW()), -- Nike - 运动服装
(51, 12, NOW(), NOW()), -- Adidas - 运动服装
(51, 41, NOW(), NOW()), -- Under Armour - 运动服装

(52, 11, NOW(), NOW()), -- Nike - 运动鞋
(52, 12, NOW(), NOW()), -- Adidas - 运动鞋
(52, 42, NOW(), NOW()), -- New Balance - 运动鞋

(54, 43, NOW(), NOW()), -- The North Face - 户外用品

-- 美妆个护品牌关联  
(61, 51, NOW(), NOW()), -- L'Oréal - 护肤
(61, 52, NOW(), NOW()), -- Estée Lauder - 护肤
(61, 53, NOW(), NOW()), -- SK-II - 护肤

(62, 51, NOW(), NOW()), -- L'Oréal - 彩妆
(62, 52, NOW(), NOW()); -- Estée Lauder - 彩妆

-- 说明：
-- 品牌涵盖各主要分类，包含国际知名品牌和国内品牌
-- category_brand 表建立了分类和品牌的多对多关系
-- 一个品牌可以关联多个分类，一个分类可以有多个品牌