-- 商品测试数据
-- 创建各个分类的真实商品数据，包含商品信息、价格、图片等

-- 清空现有数据
DELETE FROM goods_category;
DELETE FROM goods;

-- 重置自增ID
ALTER TABLE goods AUTO_INCREMENT = 1;

-- 插入商品数据
INSERT INTO goods (id, brand_id, on_sale, ship_free, is_new, is_hot, name, goods_sn, click_num, fav_num, market_price, shop_price, goods_brief, images, desc_images, goods_front_image, status, created_at, updated_at) VALUES
-- iPhone系列 (Brand: Apple - ID: 1)
(1, 1, true, true, true, true, 'iPhone 15 Pro Max 256GB 深空黑色', 'IP15PM256DSB', 5280, 892, 9999.00, 9199.00, '6.7英寸超视网膜XDR显示屏，A17 Pro芯片，专业级摄像头系统', '["https://img.joyshop.com/iphone15promax/black_1.jpg","https://img.joyshop.com/iphone15promax/black_2.jpg","https://img.joyshop.com/iphone15promax/black_3.jpg"]', '["https://img.joyshop.com/iphone15promax/desc_1.jpg","https://img.joyshop.com/iphone15promax/desc_2.jpg","https://img.joyshop.com/iphone15promax/desc_3.jpg"]', 'https://img.joyshop.com/iphone15promax/black_front.jpg', 1, NOW(), NOW()),

(2, 1, true, true, true, true, 'iPhone 14 128GB 午夜色', 'IP14128MN', 3256, 654, 6999.00, 5999.00, '6.1英寸超视网膜XDR显示屏，A15仿生芯片，先进的双摄像头系统', '["https://img.joyshop.com/iphone14/midnight_1.jpg","https://img.joyshop.com/iphone14/midnight_2.jpg"]', '["https://img.joyshop.com/iphone14/desc_1.jpg","https://img.joyshop.com/iphone14/desc_2.jpg"]', 'https://img.joyshop.com/iphone14/midnight_front.jpg', 1, NOW(), NOW()),

-- 华为系列 (Brand: HUAWEI - ID: 2)
(3, 2, true, true, false, true, 'HUAWEI Mate 60 Pro 512GB 雅川青', 'HWM60P512YCQ', 2876, 523, 7999.00, 7299.00, '6.82英寸OLED曲面屏，麒麟9000S芯片，卫星通话功能', '["https://img.joyshop.com/mate60pro/green_1.jpg","https://img.joyshop.com/mate60pro/green_2.jpg"]', '["https://img.joyshop.com/mate60pro/desc_1.jpg","https://img.joyshop.com/mate60pro/desc_2.jpg"]', 'https://img.joyshop.com/mate60pro/green_front.jpg', 1, NOW(), NOW()),

(4, 2, true, true, true, false, 'HUAWEI P60 Pro 256GB 洛可可白', 'HWP60P256LKK', 1956, 387, 6988.00, 5988.00, '6.67英寸LTPO OLED曲面屏，骁龙8+ Gen1芯片，XMAGE影像系统', '["https://img.joyshop.com/p60pro/white_1.jpg","https://img.joyshop.com/p60pro/white_2.jpg"]', '["https://img.joyshop.com/p60pro/desc_1.jpg","https://img.joyshop.com/p60pro/desc_2.jpg"]', 'https://img.joyshop.com/p60pro/white_front.jpg', 1, NOW(), NOW()),

-- 小米系列 (Brand: Xiaomi - ID: 3)
(5, 3, true, true, true, true, 'Xiaomi 14 Ultra 16GB+512GB 钛金属', 'XM14U16512TJS', 4521, 789, 6999.00, 6299.00, '6.73英寸2K LTPO OLED曲面屏，骁龙8 Gen3芯片，徕卡光学镜头', '["https://img.joyshop.com/xiaomi14ultra/titanium_1.jpg","https://img.joyshop.com/xiaomi14ultra/titanium_2.jpg"]', '["https://img.joyshop.com/xiaomi14ultra/desc_1.jpg","https://img.joyshop.com/xiaomi14ultra/desc_2.jpg"]', 'https://img.joyshop.com/xiaomi14ultra/titanium_front.jpg', 1, NOW(), NOW()),

-- 笔记本电脑
(6, 1, true, true, false, true, 'MacBook Pro 14英寸 M3芯片 512GB 深空灰色', 'MBP14M3512SG', 2145, 456, 16999.00, 15999.00, '14英寸Liquid视网膜XDR显示屏，M3芯片，16GB统一内存', '["https://img.joyshop.com/mbp14/spacegray_1.jpg","https://img.joyshop.com/mbp14/spacegray_2.jpg"]', '["https://img.joyshop.com/mbp14/desc_1.jpg","https://img.joyshop.com/mbp14/desc_2.jpg"]', 'https://img.joyshop.com/mbp14/spacegray_front.jpg', 1, NOW(), NOW()),

(7, 5, true, true, true, false, 'ThinkPad X1 Carbon Gen11 14英寸', 'TPX1C11', 1876, 234, 13999.00, 11999.00, '14英寸2.8K IPS防眩光显示屏，第13代Intel Core i7处理器', '["https://img.joyshop.com/thinkpadx1/black_1.jpg","https://img.joyshop.com/thinkpadx1/black_2.jpg"]', '["https://img.joyshop.com/thinkpadx1/desc_1.jpg","https://img.joyshop.com/thinkpadx1/desc_2.jpg"]', 'https://img.joyshop.com/thinkpadx1/black_front.jpg', 1, NOW(), NOW()),

-- 运动鞋类
(8, 11, true, false, true, true, 'Nike Air Jordan 1 Retro High OG 黑红配色', 'NAJR1RHBR', 8965, 1245, 1399.00, 1199.00, '经典篮球鞋，优质皮革材质，气垫缓震技术', '["https://img.joyshop.com/aj1/bred_1.jpg","https://img.joyshop.com/aj1/bred_2.jpg","https://img.joyshop.com/aj1/bred_3.jpg"]', '["https://img.joyshop.com/aj1/desc_1.jpg","https://img.joyshop.com/aj1/desc_2.jpg"]', 'https://img.joyshop.com/aj1/bred_front.jpg', 1, NOW(), NOW()),

(9, 12, true, false, false, true, 'Adidas Ultraboost 22 跑步鞋 黑白配色', 'ADUB22BW', 3245, 687, 1599.00, 1299.00, 'BOOST中底缓震科技，Primeknit鞋面，专业跑步鞋', '["https://img.joyshop.com/ultraboost22/blackwhite_1.jpg","https://img.joyshop.com/ultraboost22/blackwhite_2.jpg"]', '["https://img.joyshop.com/ultraboost22/desc_1.jpg","https://img.joyshop.com/ultraboost22/desc_2.jpg"]', 'https://img.joyshop.com/ultraboost22/blackwhite_front.jpg', 1, NOW(), NOW()),

-- 服装类
(10, 13, true, true, true, false, 'UNIQLO 男士纯棉T恤 白色 L码', 'UQ001WL', 1564, 289, 99.00, 79.00, '100%纯棉材质，舒适透气，多色可选', '["https://img.joyshop.com/uniqlo-tshirt/white_1.jpg","https://img.joyshop.com/uniqlo-tshirt/white_2.jpg"]', '["https://img.joyshop.com/uniqlo-tshirt/desc_1.jpg","https://img.joyshop.com/uniqlo-tshirt/desc_2.jpg"]', 'https://img.joyshop.com/uniqlo-tshirt/white_front.jpg', 1, NOW(), NOW()),

(11, 16, true, false, false, true, 'Levi\'s 501经典直筒牛仔裤 深蓝色 32码', 'LV501DB32', 2356, 445, 799.00, 699.00, '经典501版型，100%纯棉牛仔布，美国原装进口', '["https://img.joyshop.com/levis501/darkblue_1.jpg","https://img.joyshop.com/levis501/darkblue_2.jpg"]', '["https://img.joyshop.com/levis501/desc_1.jpg","https://img.joyshop.com/levis501/desc_2.jpg"]', 'https://img.joyshop.com/levis501/darkblue_front.jpg', 1, NOW(), NOW()),

-- 家电类
(12, 23, true, true, true, true, '美的变频空调 1.5匹 一级能效 静音型', 'MD15BPJY', 1234, 298, 3999.00, 3299.00, '变频压缩机，一级能效，静音运行，智能控制', '["https://img.joyshop.com/midea-ac/white_1.jpg","https://img.joyshop.com/midea-ac/white_2.jpg"]', '["https://img.joyshop.com/midea-ac/desc_1.jpg","https://img.joyshop.com/midea-ac/desc_2.jpg"]', 'https://img.joyshop.com/midea-ac/white_front.jpg', 1, NOW(), NOW()),

(13, 24, true, true, false, false, '海尔双开门冰箱 408L 变频静音', 'HR408SK', 987, 156, 4999.00, 4199.00, '408L大容量，变频压缩机，风冷无霜，节能静音', '["https://img.joyshop.com/haier-fridge/silver_1.jpg","https://img.joyshop.com/haier-fridge/silver_2.jpg"]', '["https://img.joyshop.com/haier-fridge/desc_1.jpg","https://img.joyshop.com/haier-fridge/desc_2.jpg"]', 'https://img.joyshop.com/haier-fridge/silver_front.jpg', 1, NOW(), NOW()),

-- 相机类
(14, 8, true, false, true, true, 'Canon EOS R6 Mark II 微单相机 单机身', 'CANR6M2', 756, 134, 16999.00, 14999.00, '2420万像素全画幅CMOS，双像素CMOS AF II，8级防抖', '["https://img.joyshop.com/canonr6m2/black_1.jpg","https://img.joyshop.com/canonr6m2/black_2.jpg"]', '["https://img.joyshop.com/canonr6m2/desc_1.jpg","https://img.joyshop.com/canonr6m2/desc_2.jpg"]', 'https://img.joyshop.com/canonr6m2/black_front.jpg', 1, NOW(), NOW()),

-- 图书类
(15, 31, true, true, false, false, '深入理解计算机系统（原书第3版）', 'CSAPP3', 2145, 567, 139.00, 99.00, '计算机科学经典教材，深入讲解计算机系统原理', '["https://img.joyshop.com/csapp/cover_1.jpg"]', '["https://img.joyshop.com/csapp/desc_1.jpg","https://img.joyshop.com/csapp/desc_2.jpg"]', 'https://img.joyshop.com/csapp/cover_front.jpg', 1, NOW(), NOW()),

(16, 32, true, true, true, true, 'Python编程：从入门到实践（第2版）', 'PYTHON2', 3456, 789, 89.00, 69.00, 'Python编程入门经典教材，适合初学者', '["https://img.joyshop.com/python2/cover_1.jpg"]', '["https://img.joyshop.com/python2/desc_1.jpg","https://img.joyshop.com/python2/desc_2.jpg"]', 'https://img.joyshop.com/python2/cover_front.jpg', 1, NOW(), NOW()),

-- 美妆护肤类
(17, 53, true, true, true, true, 'SK-II 神仙水 230ml 正装', 'SKII230', 4521, 892, 1690.00, 1390.00, 'Pitera™精华，改善肌理，提升肌肤光泽度', '["https://img.joyshop.com/skii-fts/bottle_1.jpg","https://img.joyshop.com/skii-fts/bottle_2.jpg"]', '["https://img.joyshop.com/skii-fts/desc_1.jpg","https://img.joyshop.com/skii-fts/desc_2.jpg"]', 'https://img.joyshop.com/skii-fts/bottle_front.jpg', 1, NOW(), NOW()),

-- 家具类
(18, 21, true, false, false, true, 'IKEA 马尔姆 双人床架 白色 1.8m', 'IK-MALM180W', 567, 89, 899.00, 699.00, '实木材质，简约北欧风格，坚固耐用', '["https://img.joyshop.com/ikea-malm/white_1.jpg","https://img.joyshop.com/ikea-malm/white_2.jpg"]', '["https://img.joyshop.com/ikea-malm/desc_1.jpg","https://img.joyshop.com/ikea-malm/desc_2.jpg"]', 'https://img.joyshop.com/ikea-malm/white_front.jpg', 1, NOW(), NOW()),

-- 更多商品...
(19, 3, true, true, true, false, 'Xiaomi 小米电视65英寸4K智能电视', 'XMTV65', 1234, 234, 3999.00, 3499.00, '65英寸4K超高清显示，MIUI TV系统，语音控制', '["https://img.joyshop.com/xiaomi-tv65/black_1.jpg","https://img.joyshop.com/xiaomi-tv65/black_2.jpg"]', '["https://img.joyshop.com/xiaomi-tv65/desc_1.jpg","https://img.joyshop.com/xiaomi-tv65/desc_2.jpg"]', 'https://img.joyshop.com/xiaomi-tv65/black_front.jpg', 1, NOW(), NOW()),

(20, 33, true, true, false, false, '晨光文具 中性笔 0.5mm 黑色 10支装', 'CG05B10', 856, 142, 25.00, 19.80, '书写流畅，笔墨均匀，办公学习必备', '["https://img.joyshop.com/chenguang-pen/black_1.jpg"]', '["https://img.joyshop.com/chenguang-pen/desc_1.jpg"]', 'https://img.joyshop.com/chenguang-pen/black_front.jpg', 1, NOW(), NOW());

-- 建立商品与分类的关联关系
INSERT INTO goods_category (goods_id, category_id) VALUES
-- iPhone 15 Pro Max 关联分类
(1, 1),   -- 电子数码
(1, 11),  -- 手机通讯
(1, 111), -- iPhone

-- iPhone 14 关联分类
(2, 1),   -- 电子数码  
(2, 11),  -- 手机通讯
(2, 111), -- iPhone

-- HUAWEI Mate 60 Pro 关联分类
(3, 1),   -- 电子数码
(3, 11),  -- 手机通讯
(3, 112), -- Android手机

-- HUAWEI P60 Pro 关联分类
(4, 1),   -- 电子数码
(4, 11),  -- 手机通讯
(4, 112), -- Android手机

-- Xiaomi 14 Ultra 关联分类
(5, 1),   -- 电子数码
(5, 11),  -- 手机通讯
(5, 112), -- Android手机

-- MacBook Pro 关联分类
(6, 1),   -- 电子数码
(6, 12),  -- 电脑办公
(6, 121), -- 笔记本电脑

-- ThinkPad X1 关联分类
(7, 1),   -- 电子数码
(7, 12),  -- 电脑办公
(7, 121), -- 笔记本电脑

-- Nike Air Jordan 1 关联分类
(8, 5),   -- 运动户外
(8, 52),  -- 运动鞋
(8, 522), -- 篮球鞋
(8, 23),  -- 鞋类
(8, 231), -- 运动鞋

-- Adidas Ultraboost 关联分类
(9, 5),   -- 运动户外
(9, 52),  -- 运动鞋
(9, 521), -- 跑步鞋
(9, 23),  -- 鞋类
(9, 231), -- 运动鞋

-- UNIQLO T恤 关联分类
(10, 2),   -- 服装鞋包
(10, 21),  -- 男装
(10, 211), -- T恤

-- Levi's 牛仔裤 关联分类
(11, 2),   -- 服装鞋包
(11, 21),  -- 男装
(11, 213), -- 牛仔裤

-- 美的空调 关联分类
(12, 1),  -- 电子数码
(12, 13), -- 家用电器
(12, 131), -- 大家电

-- 海尔冰箱 关联分类
(13, 1),  -- 电子数码
(13, 13), -- 家用电器
(13, 131), -- 大家电

-- Canon 相机 关联分类
(14, 1),  -- 电子数码
(14, 14), -- 摄影摄像
(14, 142), -- 微单相机

-- 计算机系统图书 关联分类
(15, 4),  -- 图书文教
(15, 41), -- 图书

-- Python编程图书 关联分类
(16, 4),  -- 图书文教
(16, 41), -- 图书

-- SK-II 神仙水 关联分类
(17, 6),  -- 美妆个护
(17, 61), -- 护肤

-- IKEA 床架 关联分类
(18, 3),  -- 家居生活
(18, 31), -- 家具

-- 小米电视 关联分类
(19, 1),  -- 电子数码
(19, 13), -- 家用电器
(19, 131), -- 大家电

-- 晨光中性笔 关联分类
(20, 4),  -- 图书文教
(20, 42); -- 文具用品

-- 说明：
-- 商品数据涵盖各个主要分类
-- 价格设置有市场价和店铺价，体现促销
-- 点击数和收藏数模拟用户行为数据
-- 商品编号采用品牌+型号的规则
-- 图片URL使用CDN地址，实际项目中需要替换为真实地址
-- 商品状态：1=正常，0=下架