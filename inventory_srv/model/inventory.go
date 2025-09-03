/*
 * @Author: Will nanan_zhao@163.com
 * @Date: 2025-05-12 17:13:18
 * @LastEditors: Please set LastEditors
 * @LastEditTime: 2025-06-04 10:19:17
 * @FilePath: /joyshop_srvs/inventory_srv/model/inventory.go
 * @Description: 库存服务数据模型定义，包含库存基本信息、商品关联和版本控制
 */
package model

//	type StockInfo struct {
//		BaseModel
//		Name 	string `json:"name" gorm:"type:varchar(100);not null;comment:仓库名称"`
//		Address string `json:"address" gorm:"type:varchar(255);not null;comment:仓库地址"`
//	}
type Inventory struct {
	BaseModel
	GoodsID int32 `json:"goods_id" gorm:"type:int;not null;index:idx_goods_id;comment:商品ID"` // 商品ID
	Stock   int32 `json:"stock" gorm:"type:int;not null;default:0;comment:库存数量"`
	Version int32 `json:"version" gorm:"type:int;not null;default:0;comment:版本号"` // 分布式锁使用的版本号（乐观锁）
}

// 哪个用户，哪个商品，哪个订单
type InventoryHistory struct {
	user  int32
	goods int32
	nums  int32
	order int32
	state int32 // 1 预扣减 2 表示已经支付成功
}
