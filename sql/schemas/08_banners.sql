-- 轮播图测试数据
-- 创建首页轮播图数据

-- 清空现有数据
DELETE FROM banner;

-- 重置自增ID
ALTER TABLE banner AUTO_INCREMENT = 1;

-- 插入轮播图数据
INSERT INTO banner (id, image, url, `index`, created_at, updated_at) VALUES
(1, 'https://img.joyshop.com/banners/iphone15-banner.jpg', '/category/111', 1, NOW(), NOW()),  -- iPhone系列促销
(2, 'https://img.joyshop.com/banners/nike-banner.jpg', '/category/231', 2, NOW(), NOW()),     -- Nike运动鞋促销  
(3, 'https://img.joyshop.com/banners/macbook-banner.jpg', '/category/121', 3, NOW(), NOW()),  -- MacBook促销
(4, 'https://img.joyshop.com/banners/xiaomi-banner.jpg', '/category/13', 4, NOW(), NOW()),    -- 小米家电促销
(5, 'https://img.joyshop.com/banners/fashion-banner.jpg', '/category/2', 5, NOW(), NOW());    -- 服装时尚促销

-- 说明：
-- image: 轮播图片URL
-- url: 点击跳转链接，指向具体分类页面
-- index: 显示顺序，数字越小越靠前
-- 轮播图内容覆盖主要商品分类，引导用户浏览