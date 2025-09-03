/*
 * @Author: Will nanan_zhao@163.com
 * @Date: 2025-05-12 17:13:18
 * @LastEditors: Please set LastEditors
 * @LastEditTime: 2025-06-08 14:12:50
 * @FilePath: /joyshop_srvs/order_srv/model/inventory.go
 * @Description: 库存服务数据模型定义，包含库存基本信息、商品关联和版本控制
 */
package model

import (
	"time"
)

type ShoppingCart struct {
	BaseModel
	User    int32 `gorm:"type:int;index;not null;comment:用户ID" json:"user"` // 在购物车列表中我们要查询当前用户的购物车记录
	Goods   int32 `gorm:"type:int;index;not null;comment:<UNK>" json:"goods"`
	Nums    int32 `gorm:"type:int;index;not null;comment:<UNK>" json:"nums"`
	Checked bool  `gorm:"type:bool;default:true;comment:是否选中" json:"checked"` //是否选中
}

type OrderInfo struct {
	BaseModel
	User    int32  `gorm:"type:int;index"`
	OrderSn string `gorm:"type:varchar(30);index"` // 平台自己生成的订单号
	PayType string `gorm:"type:varchar(20);comment:'alipay(支付宝), wechat(微信)'"`
	// status大家可以考虑用iota来做
	Status       string `gorm:"type:varchar(20);comment:'PAYING(待支付), TRADE_SUCCESS(成功), TRADE_CLOSED(超时关闭), WAIT_BUYER_PAY(交易创建), TRADE_FINISHED(交易结束)'"`
	TradeNo      string `gorm:"type:varchar(100);comment:'交易号'"` // 交易号就是支付订单号
	OrderMount   float32
	PayTime      time.Time
	Address      string `gorm:"type:varchar(100)"`
	SignerName   string `gorm:"type:varchar(20)"`
	SingerMobile string `gorm:"type:varchar(11)"`
	Post         string `gorm:"type:varchar(20)"` //留言信息
}

// 订单商品信息
type OrderGoods struct {
	BaseModel
	Order int32 `gorm:"type:int;index"`
	Goods int32 `gorm:"type:int;index"`
	// 商品名称、商品图片、商品价格、商品数量，高并发场景下都不会遵守第三范式，所以这里不使用外键（字段冗余）
	GoodsName  string  `gorm:"type:varchar(100);index"`
	GoodsImage string  `gorm:"type:varchar(200)"`
	GoodsPrice float32 // 快照价格
	Nums       int32   `gorm:"type:int"`
}
