#!/bin/bash

# JoyShop 数据导入脚本
# 用于清理旧数据并导入新的测试数据

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 数据库连接参数
DB_HOST="127.0.0.1"
DB_PORT="3306"
DB_USER="root"

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}JoyShop 数据导入工具${NC}"
echo -e "${GREEN}========================================${NC}"

# 提示输入密码
echo -n "请输入MySQL密码: "
read -s DB_PASS
echo ""

# 切换到SQL目录
cd /Users/walker/gitsource/github/joyshop_srvs/sql

# 步骤1: 清理现有数据
echo -e "\n${YELLOW}步骤 1: 清理现有测试数据...${NC}"
mysql -h $DB_HOST -P $DB_PORT -u $DB_USER -p$DB_PASS < scripts/clean_data.sql
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ 数据清理完成${NC}"
else
    echo -e "${RED}❌ 数据清理失败${NC}"
    exit 1
fi

# 步骤2: 导入商品服务数据
echo -e "\n${YELLOW}步骤 2: 导入商品服务数据...${NC}"

echo "  导入分类数据..."
mysql -h $DB_HOST -P $DB_PORT -u $DB_USER -p$DB_PASS joyshop_goods < schemas/02_categories.sql
if [ $? -eq 0 ]; then
    echo -e "  ${GREEN}✅ 分类数据导入成功${NC}"
else
    echo -e "  ${RED}❌ 分类数据导入失败${NC}"
    exit 1
fi

echo "  导入品牌数据..."
mysql -h $DB_HOST -P $DB_PORT -u $DB_USER -p$DB_PASS joyshop_goods < schemas/03_brands.sql
if [ $? -eq 0 ]; then
    echo -e "  ${GREEN}✅ 品牌数据导入成功${NC}"
else
    echo -e "  ${RED}❌ 品牌数据导入失败${NC}"
    exit 1
fi

echo "  导入商品数据(3500条)..."
mysql -h $DB_HOST -P $DB_PORT -u $DB_USER -p$DB_PASS joyshop_goods < schemas/04_goods.sql
if [ $? -eq 0 ]; then
    echo -e "  ${GREEN}✅ 商品数据导入成功${NC}"
else
    echo -e "  ${RED}❌ 商品数据导入失败${NC}"
    exit 1
fi

echo "  导入轮播图数据..."
mysql -h $DB_HOST -P $DB_PORT -u $DB_USER -p$DB_PASS joyshop_goods < schemas/08_banners.sql
if [ $? -eq 0 ]; then
    echo -e "  ${GREEN}✅ 轮播图数据导入成功${NC}"
else
    echo -e "  ${RED}❌ 轮播图数据导入失败${NC}"
    exit 1
fi

# 步骤3: 导入库存数据
echo -e "\n${YELLOW}步骤 3: 导入库存数据...${NC}"
mysql -h $DB_HOST -P $DB_PORT -u $DB_USER -p$DB_PASS joyshop_inventory < schemas/05_inventory.sql
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ 库存数据导入成功${NC}"
else
    echo -e "${RED}❌ 库存数据导入失败${NC}"
    exit 1
fi

# 步骤4: 导入用户数据
echo -e "\n${YELLOW}步骤 4: 导入用户数据...${NC}"
mysql -h $DB_HOST -P $DB_PORT -u $DB_USER -p$DB_PASS joyshop_user < schemas/01_users.sql
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ 用户数据导入成功${NC}"
else
    echo -e "${RED}❌ 用户数据导入失败${NC}"
    exit 1
fi

# 步骤5: 导入订单相关数据
echo -e "\n${YELLOW}步骤 5: 导入订单相关数据...${NC}"

echo "  导入购物车数据..."
mysql -h $DB_HOST -P $DB_PORT -u $DB_USER -p$DB_PASS joyshop_order < schemas/06_shopping_cart.sql
if [ $? -eq 0 ]; then
    echo -e "  ${GREEN}✅ 购物车数据导入成功${NC}"
else
    echo -e "  ${RED}❌ 购物车数据导入失败${NC}"
    exit 1
fi

echo "  导入订单数据..."
mysql -h $DB_HOST -P $DB_PORT -u $DB_USER -p$DB_PASS joyshop_order < schemas/07_orders.sql
if [ $? -eq 0 ]; then
    echo -e "  ${GREEN}✅ 订单数据导入成功${NC}"
else
    echo -e "  ${RED}❌ 订单数据导入失败${NC}"
    exit 1
fi

# 步骤6: 验证数据
echo -e "\n${YELLOW}步骤 6: 验证数据完整性...${NC}"
mysql -h $DB_HOST -P $DB_PORT -u $DB_USER -p$DB_PASS < scripts/validate_data_integrity.sql
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ 数据验证完成${NC}"
else
    echo -e "${RED}❌ 数据验证失败${NC}"
    exit 1
fi

echo -e "\n${GREEN}========================================${NC}"
echo -e "${GREEN}🎉 所有数据导入成功！${NC}"
echo -e "${GREEN}========================================${NC}"

# 显示统计信息
echo -e "\n${YELLOW}数据统计：${NC}"
mysql -h $DB_HOST -P $DB_PORT -u $DB_USER -p$DB_PASS -e "
USE joyshop_goods;
SELECT '商品' as type, COUNT(*) as count FROM goods
UNION ALL SELECT '分类', COUNT(*) FROM category
UNION ALL SELECT '品牌', COUNT(*) FROM brand;

USE joyshop_inventory;
SELECT '库存' as type, COUNT(*) as count FROM inventory;

USE joyshop_user;
SELECT '用户' as type, COUNT(*) as count FROM users;

USE joyshop_order;
SELECT '订单' as type, COUNT(*) as count FROM order_info
UNION ALL SELECT '购物车', COUNT(*) FROM shopping_cart;
"

echo -e "\n${GREEN}数据导入完成！可以开始测试了。${NC}"