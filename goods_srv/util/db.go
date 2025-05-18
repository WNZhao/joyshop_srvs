/*
 * @Author: Will nanan_zhao@163.com
 * @Date: 2025-05-18 19:13:14
 * @LastEditors: Will nanan_zhao@163.com
 * @LastEditTime: 2025-05-18 19:24:56
 * @FilePath: /joyshop_srvs/goods_srv/util/db.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package util

import (
	"goods_srv/global"
	"goods_srv/model"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// 初始化测试数据库表结构（如不存在则自动迁移）
func SetupTestTables() error {
	zap.S().Info("检查并自动迁移测试数据库表结构...")
	return global.DB.AutoMigrate(
		&model.Category{},
		&model.Brand{},
		&model.Goods{},
		&model.CategoryBrand{},
		&model.Banner{},
	)
}

// 清空所有测试表数据
func CleanTestTables() error {
	zap.S().Info("清空所有测试表数据...")
	tables := []interface{}{
		&model.Goods{},
		&model.CategoryBrand{},
		&model.Category{},
		&model.Brand{},
		&model.Banner{},
	}
	for _, table := range tables {
		if err := global.DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(table).Error; err != nil {
			return err
		}
	}
	return nil
}

// 删除所有测试表结构
func DropAllTables() error {
	zap.S().Info("删除所有测试表结构...")
	return global.DB.Migrator().DropTable(
		&model.Goods{},
		&model.CategoryBrand{},
		&model.Category{},
		&model.Brand{},
		&model.Banner{},
	)
}
