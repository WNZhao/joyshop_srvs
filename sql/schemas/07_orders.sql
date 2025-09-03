-- 订单服务测试数据
-- 创建订单信息和订单商品数据，模拟不同状态的订单

-- 清空现有数据
DELETE FROM order_goods;
DELETE FROM order_info;

-- 重置自增ID
ALTER TABLE order_info AUTO_INCREMENT = 1;
ALTER TABLE order_goods AUTO_INCREMENT = 1;

-- 插入订单信息数据
INSERT INTO order_info (id, user, order_sn, pay_type, status, trade_no, order_mount, pay_time, address, signer_name, singer_mobile, post, created_at, updated_at, is_deleted) VALUES
-- 已完成的订单（用于测试完整流程）
(1, 2, '20241201120000000000020001', 'alipay', 'TRADE_FINISHED', 'ALI_20241201_001', 10588.00, '2024-12-01 12:30:00', '北京市朝阳区建国路1号', '张三', '13800000002', '请尽快发货，谢谢', '2024-12-01 12:00:00', '2024-12-01 15:00:00', false),

(2, 3, '20241201130000000000030001', 'wechat', 'TRADE_FINISHED', 'WX_20241201_002', 2688.00, '2024-12-01 13:15:00', '上海市浦东新区陆家嘴环路1000号', '李四', '13800000003', '工作日配送', '2024-12-01 13:00:00', '2024-12-01 16:00:00', false),

(3, 4, '20241202090000000000040001', 'alipay', 'TRADE_SUCCESS', 'ALI_20241202_003', 31067.00, '2024-12-02 09:30:00', '广州市天河区珠江新城华夏路1号', '王五', '13800000004', '请包装仔细，贵重物品', '2024-12-02 09:00:00', '2024-12-02 12:00:00', false),

-- 待支付订单
(4, 5, '20241203100000000000050001', 'alipay', 'WAIT_BUYER_PAY', '', 12697.00, NULL, '深圳市南山区科技园南区深南大道999号', '赵六', '13800000005', '家用电器，请确保包装完整', '2024-12-03 10:00:00', '2024-12-03 10:00:00', false),

-- 支付中订单
(5, 6, '20241203140000000000060001', 'wechat', 'PAYING', '', 2897.00, NULL, '杭州市西湖区文三路1号', '孙七', '13800000006', '运动用品订单', '2024-12-03 14:00:00', '2024-12-03 14:05:00', false),

-- 已关闭订单（超时未支付）
(6, 7, '20241201080000000000070001', 'alipay', 'TRADE_CLOSED', '', 6178.00, NULL, '成都市锦江区红星路2段1号', '周八', '13800000007', '学生用品', '2024-12-01 08:00:00', '2024-12-01 20:00:00', false),

-- VIP大额订单
(7, 13, '20241203160000000000130001', 'alipay', 'TRADE_SUCCESS', 'ALI_20241203_007', 45876.00, '2024-12-03 16:30:00', '北京市海淀区中关村大街1号', 'VIP刘总', '13800000013', '企业采购，请开发票', '2024-12-03 16:00:00', '2024-12-03 18:00:00', false),

(8, 14, '20241203170000000000140001', 'wechat', 'TRADE_SUCCESS', 'WX_20241203_008', 23675.00, '2024-12-03 17:15:00', '上海市黄浦区南京东路100号', 'VIP张总', '13800000014', '高端用户，优先配送', '2024-12-03 17:00:00', '2024-12-03 19:00:00', false),

-- 近期订单（不同状态用于测试）
(9, 10, '20241204090000000000100001', 'alipay', 'WAIT_BUYER_PAY', '', 99.80, NULL, '武汉市武昌区中南路1号', '陈一', '13800000010', '', '2024-12-04 09:00:00', '2024-12-04 09:00:00', false),

(10, 11, '20241204100000000000110001', 'wechat', 'TRADE_SUCCESS', 'WX_20241204_010', 857.00, '2024-12-04 10:30:00', '西安市雁塔区高新路1号', '李二', '13800000011', '服装订单', '2024-12-04 10:00:00', '2024-12-04 13:00:00', false);

-- 插入订单商品数据
INSERT INTO order_goods (id, `order`, goods, goods_name, goods_image, goods_price, nums, created_at, updated_at, is_deleted) VALUES
-- 订单1的商品（张三的订单：iPhone + AJ1 + T恤）
(1, 1, 1, 'iPhone 15 Pro Max 256GB 深空黑色', 'https://img.joyshop.com/iphone15promax/black_front.jpg', 9199.00, 1, '2024-12-01 12:00:00', '2024-12-01 12:00:00', false),
(2, 1, 8, 'Nike Air Jordan 1 Retro High OG 黑红配色', 'https://img.joyshop.com/aj1/bred_front.jpg', 1199.00, 1, '2024-12-01 12:00:00', '2024-12-01 12:00:00', false),
(3, 1, 10, 'UNIQLO 男士纯棉T恤 白色 L码', 'https://img.joyshop.com/uniqlo-tshirt/white_front.jpg', 79.00, 2, '2024-12-01 12:00:00', '2024-12-01 12:00:00', false),

-- 订单2的商品（李四的订单：SK-II + Adidas跑步鞋 + 中性笔）
(4, 2, 17, 'SK-II 神仙水 230ml 正装', 'https://img.joyshop.com/skii-fts/bottle_front.jpg', 1390.00, 1, '2024-12-01 13:00:00', '2024-12-01 13:00:00', false),
(5, 2, 9, 'Adidas Ultraboost 22 跑步鞋 黑白配色', 'https://img.joyshop.com/ultraboost22/blackwhite_front.jpg', 1299.00, 1, '2024-12-01 13:00:00', '2024-12-01 13:00:00', false),
(6, 2, 20, '晨光文具 中性笔 0.5mm 黑色 10支装', 'https://img.joyshop.com/chenguang-pen/black_front.jpg', 19.80, 5, '2024-12-01 13:00:00', '2024-12-01 13:00:00', false),

-- 订单3的商品（王五的订单：MacBook + Canon相机 + Python书）
(7, 3, 6, 'MacBook Pro 14英寸 M3芯片 512GB 深空灰色', 'https://img.joyshop.com/mbp14/spacegray_front.jpg', 15999.00, 1, '2024-12-02 09:00:00', '2024-12-02 09:00:00', false),
(8, 3, 14, 'Canon EOS R6 Mark II 微单相机 单机身', 'https://img.joyshop.com/canonr6m2/black_front.jpg', 14999.00, 1, '2024-12-02 09:00:00', '2024-12-02 09:00:00', false),
(9, 3, 16, 'Python编程：从入门到实践（第2版）', 'https://img.joyshop.com/python2/cover_front.jpg', 69.00, 1, '2024-12-02 09:00:00', '2024-12-02 09:00:00', false),

-- 订单4的商品（赵六的订单：家电套装）
(10, 4, 12, '美的变频空调 1.5匹 一级能效 静音型', 'https://img.joyshop.com/midea-ac/white_front.jpg', 3299.00, 1, '2024-12-03 10:00:00', '2024-12-03 10:00:00', false),
(11, 4, 13, '海尔双开门冰箱 408L 变频静音', 'https://img.joyshop.com/haier-fridge/silver_front.jpg', 4199.00, 1, '2024-12-03 10:00:00', '2024-12-03 10:00:00', false),
(12, 4, 18, 'IKEA 马尔姆 双人床架 白色 1.8m', 'https://img.joyshop.com/ikea-malm/white_front.jpg', 699.00, 1, '2024-12-03 10:00:00', '2024-12-03 10:00:00', false),
(13, 4, 19, 'Xiaomi 小米电视65英寸4K智能电视', 'https://img.joyshop.com/xiaomi-tv65/black_front.jpg', 3499.00, 1, '2024-12-03 10:00:00', '2024-12-03 10:00:00', false),

-- 订单5的商品（孙七的订单：运动套装）
(14, 5, 8, 'Nike Air Jordan 1 Retro High OG 黑红配色', 'https://img.joyshop.com/aj1/bred_front.jpg', 1199.00, 1, '2024-12-03 14:00:00', '2024-12-03 14:00:00', false),
(15, 5, 9, 'Adidas Ultraboost 22 跑步鞋 黑白配色', 'https://img.joyshop.com/ultraboost22/blackwhite_front.jpg', 1299.00, 2, '2024-12-03 14:00:00', '2024-12-03 14:00:00', false),
(16, 5, 10, 'UNIQLO 男士纯棉T恤 白色 L码', 'https://img.joyshop.com/uniqlo-tshirt/white_front.jpg', 79.00, 3, '2024-12-03 14:00:00', '2024-12-03 14:00:00', false),

-- 订单6的商品（周八的关闭订单：学生用品）
(17, 6, 2, 'iPhone 14 128GB 午夜色', 'https://img.joyshop.com/iphone14/midnight_front.jpg', 5999.00, 1, '2024-12-01 08:00:00', '2024-12-01 08:00:00', false),
(18, 6, 15, '深入理解计算机系统（原书第3版）', 'https://img.joyshop.com/csapp/cover_front.jpg', 99.00, 1, '2024-12-01 08:00:00', '2024-12-01 08:00:00', false),
(19, 6, 16, 'Python编程：从入门到实践（第2版）', 'https://img.joyshop.com/python2/cover_front.jpg', 69.00, 1, '2024-12-01 08:00:00', '2024-12-01 08:00:00', false),
(20, 6, 20, '晨光文具 中性笔 0.5mm 黑色 10支装', 'https://img.joyshop.com/chenguang-pen/black_front.jpg', 19.80, 10, '2024-12-01 08:00:00', '2024-12-01 08:00:00', false),

-- 订单7的商品（VIP刘总的大额订单）
(21, 7, 1, 'iPhone 15 Pro Max 256GB 深空黑色', 'https://img.joyshop.com/iphone15promax/black_front.jpg', 9199.00, 2, '2024-12-03 16:00:00', '2024-12-03 16:00:00', false),
(22, 7, 6, 'MacBook Pro 14英寸 M3芯片 512GB 深空灰色', 'https://img.joyshop.com/mbp14/spacegray_front.jpg', 15999.00, 1, '2024-12-03 16:00:00', '2024-12-03 16:00:00', false),
(23, 7, 14, 'Canon EOS R6 Mark II 微单相机 单机身', 'https://img.joyshop.com/canonr6m2/black_front.jpg', 14999.00, 1, '2024-12-03 16:00:00', '2024-12-03 16:00:00', false),
(24, 7, 12, '美的变频空调 1.5匹 一级能效 静音型', 'https://img.joyshop.com/midea-ac/white_front.jpg', 3299.00, 2, '2024-12-03 16:00:00', '2024-12-03 16:00:00', false),

-- 订单8的商品（VIP张总的订单）
(25, 8, 17, 'SK-II 神仙水 230ml 正装', 'https://img.joyshop.com/skii-fts/bottle_front.jpg', 1390.00, 5, '2024-12-03 17:00:00', '2024-12-03 17:00:00', false),
(26, 8, 6, 'MacBook Pro 14英寸 M3芯片 512GB 深空灰色', 'https://img.joyshop.com/mbp14/spacegray_front.jpg', 15999.00, 1, '2024-12-03 17:00:00', '2024-12-03 17:00:00', false),
(27, 8, 9, 'Adidas Ultraboost 22 跑步鞋 黑白配色', 'https://img.joyshop.com/ultraboost22/blackwhite_front.jpg', 1299.00, 2, '2024-12-03 17:00:00', '2024-12-03 17:00:00', false),
(28, 8, 13, '海尔双开门冰箱 408L 变频静音', 'https://img.joyshop.com/haier-fridge/silver_front.jpg', 4199.00, 1, '2024-12-03 17:00:00', '2024-12-03 17:00:00', false),

-- 订单9的商品（陈一的小额订单）
(29, 9, 10, 'UNIQLO 男士纯棉T恤 白色 L码', 'https://img.joyshop.com/uniqlo-tshirt/white_front.jpg', 79.00, 1, '2024-12-04 09:00:00', '2024-12-04 09:00:00', false),
(30, 9, 20, '晨光文具 中性笔 0.5mm 黑色 10支装', 'https://img.joyshop.com/chenguang-pen/black_front.jpg', 19.80, 2, '2024-12-04 09:00:00', '2024-12-04 09:00:00', false),

-- 订单10的商品（李二的服装订单）
(31, 10, 11, 'Levi\'s 501经典直筒牛仔裤 深蓝色 32码', 'https://img.joyshop.com/levis501/darkblue_front.jpg', 699.00, 1, '2024-12-04 10:00:00', '2024-12-04 10:00:00', false),
(32, 10, 10, 'UNIQLO 男士纯棉T恤 白色 L码', 'https://img.joyshop.com/uniqlo-tshirt/white_front.jpg', 79.00, 2, '2024-12-04 10:00:00', '2024-12-04 10:00:00', false);

-- 说明：
-- 订单状态说明：
--   WAIT_BUYER_PAY: 待支付（刚创建的订单）
--   PAYING: 支付中（用户正在支付）
--   TRADE_SUCCESS: 支付成功（待发货/已发货）
--   TRADE_FINISHED: 交易完成（已收货确认）
--   TRADE_CLOSED: 交易关闭（超时未支付或取消）
-- 
-- 支付方式：
--   alipay: 支付宝
--   wechat: 微信支付
-- 
-- 订单编号格式：YYYYMMDDHHMMSS + 8位用户ID + 4位序号
-- 交易号：第三方支付平台的交易流水号
-- 
-- 数据特点：
-- 1. 包含各种状态的订单，便于测试不同场景
-- 2. 订单金额从小额到大额，覆盖不同用户群体
-- 3. 商品快照价格，防止价格变动影响历史订单
-- 4. VIP用户订单金额较高，普通用户相对较低
-- 5. 时间分布合理，有历史订单和最新订单