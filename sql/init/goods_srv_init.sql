-- 商品服务初始化脚本
-- 初始化商品服务相关的所有数据：分类、品牌、商品、轮播图

-- 设置字符集
SET NAMES utf8mb4;

-- 禁用外键检查
SET FOREIGN_KEY_CHECKS = 0;

-- 按依赖顺序执行
-- 1. 分类数据（商品依赖分类）
SOURCE 02_categories.sql;

-- 2. 品牌数据（商品依赖品牌）
SOURCE 03_brands.sql;

-- 3. 商品数据（核心数据）
SOURCE 04_goods.sql;

-- 4. 轮播图数据（首页展示）
SOURCE 08_banners.sql;

-- 启用外键检查
SET FOREIGN_KEY_CHECKS = 1;

-- 显示结果
SELECT '===========================================' AS '';
SELECT '✅ 商品服务测试数据初始化完成！' AS 'Status';
SELECT '===========================================' AS '';

-- 统计信息
SELECT 'category' AS table_name, COUNT(*) AS record_count FROM category
UNION ALL
SELECT 'brand', COUNT(*) FROM brand  
UNION ALL
SELECT 'goods', COUNT(*) FROM goods
UNION ALL
SELECT 'category_brand', COUNT(*) FROM category_brand
UNION ALL
SELECT 'goods_category', COUNT(*) FROM goods_category
UNION ALL
SELECT 'banner', COUNT(*) FROM banner;

-- 分类层级统计
SELECT '📁 分类结构：' AS '';
SELECT CONCAT('Level ', level, ': ', COUNT(*), ' 个分类') AS 'Category Structure'
FROM category 
GROUP BY level 
ORDER BY level;

-- 品牌统计
SELECT '🏷️ 品牌信息：' AS '';
SELECT CONCAT(name, ' (ID: ', id, ')') AS 'Brand List' 
FROM brand 
ORDER BY id 
LIMIT 10;

-- 商品概况
SELECT '🛍️ 商品概况：' AS '';
SELECT CONCAT('总计 ', COUNT(*), ' 个商品') AS 'Goods Summary' FROM goods;
SELECT CONCAT('上架商品：', COUNT(*), ' 个') AS 'Goods Summary' FROM goods WHERE on_sale = true;
SELECT CONCAT('热销商品：', COUNT(*), ' 个') AS 'Goods Summary' FROM goods WHERE is_hot = true;
SELECT CONCAT('新品商品：', COUNT(*), ' 个') AS 'Goods Summary' FROM goods WHERE is_new = true;

SELECT '===========================================' AS '';
SELECT '🧪 测试建议：' AS '';
SELECT '1. 测试分类查询：查看不同层级的分类' AS 'Test Suggestions';
SELECT '2. 测试商品搜索：按分类、品牌、关键词搜索' AS 'Test Suggestions';
SELECT '3. 测试商品详情：查看具体商品信息' AS 'Test Suggestions';
SELECT '4. 测试轮播图：获取首页轮播图列表' AS 'Test Suggestions';
SELECT '===========================================' AS '';