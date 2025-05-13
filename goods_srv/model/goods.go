package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"goods_srv/global"

	"gorm.io/gorm"
)

// GormList 自定义类型，用于存储字符串列表
type GormList []string

// Value 实现 driver.Valuer 接口
func (g GormList) Value() (driver.Value, error) {
	if len(g) == 0 {
		return "[]", nil
	}
	return json.Marshal(g)
}

// Scan 实现 sql.Scanner 接口
func (g *GormList) Scan(value interface{}) error {
	if value == nil {
		*g = GormList{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, g)
}

// Category 商品分类
type Category struct {
	gorm.Model
	Name     string    `gorm:"type:varchar(50);not null;comment:分类名称"`
	ParentId uint      `gorm:"type:int;default:0;comment:父分类ID;"`
	Category *Category `gorm:"foreignKey:ParentId;references:ID;constraint:OnDelete:SET NULL"`
	Level    int       `gorm:"type:int;not null;default:1;comment:分类层级"`
	Sort     int       `gorm:"type:int;not null;default:0;comment:排序"`
	IsTab    bool      `gorm:"type:boolean;not null;default:false;comment:是否显示在导航栏"`
}

// Brand 品牌
type Brand struct {
	gorm.Model
	Name string `gorm:"type:varchar(50);not null;comment:品牌名称"`
	Logo string `gorm:"type:varchar(200);not null;default:'';comment:品牌logo"`
	Desc string `gorm:"type:varchar(200);not null;default:'';comment:品牌描述"`
}

// GoodsCategoryBrand 商品分类和品牌关联表
type CategoryBrand struct {
	gorm.Model
	CategoryId uint     `gorm:"type:int;not null;comment:分类ID"`
	Category   Category `gorm:"foreignKey:CategoryId;references:ID"`
	BrandId    uint     `gorm:"type:int;not null;comment:品牌ID"`
	Brand      Brand    `gorm:"foreignKey:BrandId;references:ID"`
}

// Goods 商品
type Goods struct {
	gorm.Model
	BrandId         uint       `gorm:"type:int;not null;comment:品牌ID"`
	Brand           Brand      `gorm:"foreignKey:BrandId;references:ID"`
	OnSale          bool       `gorm:"type:boolean;not null;default:false;comment:是否上架"`
	ShipFree        bool       `gorm:"type:boolean;not null;default:false;comment:是否包邮"`
	IsNew           bool       `gorm:"type:boolean;not null;default:false;comment:是否新品"`
	IsHot           bool       `gorm:"type:boolean;not null;default:false;comment:是否热销"`
	Name            string     `gorm:"type:varchar(100);not null;comment:商品名称"`
	GoodsSn         string     `gorm:"type:varchar(50);not null;comment:商品编号"`
	MarketPrice     float64    `gorm:"type:decimal(10,2);not null;comment:市场价格"`
	ShopPrice       float64    `gorm:"type:decimal(10,2);not null;comment:本店价格"`
	GoodsBrief      string     `gorm:"type:varchar(200);not null;comment:商品简介"`
	Images          GormList   `gorm:"type:json;not null;comment:商品图片"`
	DescImages      GormList   `gorm:"type:json;not null;comment:商品详情图片"`
	GoodsFrontImage string     `gorm:"type:varchar(200);not null;comment:商品主图"`
	Status          int        `gorm:"type:tinyint;not null;default:1;comment:商品状态"`
	Categories      []Category `gorm:"many2many:goods_category;"`
}

// TableName 设置表名
func (Goods) TableName() string {
	return "goods"
}

// TableName 设置表名
func (Category) TableName() string {
	return "category"
}

// TableName 设置表名
func (Brand) TableName() string {
	return "brand"
}

// TableName 设置表名
func (CategoryBrand) TableName() string {
	return "category_brand"
}

// CreateGoods 创建商品
func CreateGoods(goods *Goods) error {
	return global.DB.Create(goods).Error
}

// GetGoodsById 根据ID获取商品
func GetGoodsById(id uint) (*Goods, error) {
	var goods Goods
	err := global.DB.Preload("Categories").First(&goods, id).Error
	if err != nil {
		return nil, err
	}
	return &goods, nil
}

// UpdateGoods 更新商品
func UpdateGoods(goods *Goods) error {
	return global.DB.Save(goods).Error
}

// DeleteGoods 删除商品
func DeleteGoods(id uint) error {
	return global.DB.Delete(&Goods{}, id).Error
}

// GetGoodsList 获取商品列表
func GetGoodsList(page, pageSize int) ([]Goods, int64, error) {
	var goods []Goods
	var total int64

	// 获取总数
	if err := global.DB.Model(&Goods{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	err := global.DB.Preload("Categories").Offset(offset).Limit(pageSize).Find(&goods).Error
	if err != nil {
		return nil, 0, err
	}

	return goods, total, nil
}

// GetGoodsByCategory 按分类获取商品
func GetGoodsByCategory(categoryId int32, page, pageSize int32) ([]Goods, int64, error) {
	var goods []Goods
	var total int64

	result := global.DB.Model(&Goods{}).Where("category_id = ?", categoryId).Count(&total)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	result = global.DB.Where("category_id = ?", categoryId).Offset(int((page - 1) * pageSize)).Limit(int(pageSize)).Find(&goods)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return goods, total, nil
}

// SearchGoods 搜索商品
func SearchGoods(keyword string, page, pageSize int32) ([]Goods, int64, error) {
	var goods []Goods
	var total int64

	result := global.DB.Model(&Goods{}).Where("name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%").Count(&total)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	result = global.DB.Where("name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%").Offset(int((page - 1) * pageSize)).Limit(int(pageSize)).Find(&goods)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return goods, total, nil
}
