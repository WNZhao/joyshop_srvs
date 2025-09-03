-- 库存服务初始化脚本
-- 初始化库存服务相关数据

-- 设置字符集
SET NAMES utf8mb4;

-- 禁用外键检查
SET FOREIGN_KEY_CHECKS = 0;

-- 执行库存数据初始化
SOURCE 05_inventory.sql;

-- 启用外键检查
SET FOREIGN_KEY_CHECKS = 1;

-- 显示结果
SELECT '===========================================' AS '';
SELECT '✅ 库存服务测试数据初始化完成！' AS 'Status';
SELECT '===========================================' AS '';

-- 统计信息
SELECT 'inventory' AS table_name, COUNT(*) AS record_count FROM inventory;

-- 库存分析
SELECT '📦 库存分析：' AS '';
SELECT CONCAT('总库存记录：', COUNT(*), ' 条') AS 'Inventory Analysis' FROM inventory;
SELECT CONCAT('总库存数量：', SUM(stock), ' 件') AS 'Inventory Analysis' FROM inventory;
SELECT CONCAT('平均库存：', ROUND(AVG(stock), 2), ' 件/SKU') AS 'Inventory Analysis' FROM inventory;

-- 库存分布
SELECT '📊 库存分布：' AS '';
SELECT 
    CASE 
        WHEN stock = 0 THEN '零库存'
        WHEN stock BETWEEN 1 AND 10 THEN '低库存(1-10)'
        WHEN stock BETWEEN 11 AND 50 THEN '正常库存(11-50)'
        WHEN stock BETWEEN 51 AND 200 THEN '充足库存(51-200)'
        ELSE '大量库存(200+)'
    END AS 'Stock Level',
    COUNT(*) AS 'Count'
FROM inventory 
GROUP BY 
    CASE 
        WHEN stock = 0 THEN '零库存'
        WHEN stock BETWEEN 1 AND 10 THEN '低库存(1-10)'
        WHEN stock BETWEEN 11 AND 50 THEN '正常库存(11-50)'
        WHEN stock BETWEEN 51 AND 200 THEN '充足库存(51-200)'
        ELSE '大量库存(200+)'
    END;

-- 库存告警商品
SELECT '⚠️ 库存告警商品（库存 ≤ 10）：' AS '';
SELECT CONCAT('商品ID ', goods_id, ': 库存 ', stock, ' 件') AS 'Low Stock Alert'
FROM inventory 
WHERE stock <= 10 
ORDER BY stock ASC, goods_id ASC;

-- 高库存商品  
SELECT '📈 高库存商品（库存 > 500）：' AS '';
SELECT CONCAT('商品ID ', goods_id, ': 库存 ', stock, ' 件') AS 'High Stock'
FROM inventory 
WHERE stock > 500 
ORDER BY stock DESC;

SELECT '===========================================' AS '';
SELECT '🧪 测试建议：' AS '';
SELECT '1. 测试库存查询：查看特定商品库存' AS 'Test Suggestions';
SELECT '2. 测试库存扣减：模拟商品购买扣减库存' AS 'Test Suggestions';
SELECT '3. 测试库存回滚：模拟订单取消回滚库存' AS 'Test Suggestions';
SELECT '4. 测试库存告警：检查低库存商品提醒' AS 'Test Suggestions';
SELECT '5. 测试并发扣减：验证乐观锁版本控制' AS 'Test Suggestions';
SELECT '===========================================' AS '';