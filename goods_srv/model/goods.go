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
	Name          string     `gorm:"type:varchar(50);not null;comment:分类名称"`
	ParentId      *int       `gorm:"type:int unsigned;default:null;comment:父分类ID;"`
	Category      *Category  `gorm:"foreignKey:ParentId;references:ID;constraint:OnDelete:SET NULL"`
	SubCategories []Category `gorm:"foreignKey:ParentId;references:ID;constraint:OnDelete:CASCADE"`
	Level         int        `gorm:"type:int;not null;default:1;comment:分类层级"`
	Sort          int        `gorm:"type:int;not null;default:0;comment:排序"`
	IsTab         bool       `gorm:"type:boolean;not null;default:false;comment:是否显示在导航栏"`
}

// Brand 品牌
type Brand struct {
	gorm.Model
	Name string `gorm:"type:varchar(50);not null;comment:品牌名称"`
	Logo string `gorm:"type:varchar(200);not null;default:'';comment:品牌logo"`
	Desc string `gorm:"type:varchar(200);not null;default:'';comment:品牌描述"`
}

// GoodsCategoryBrand 商品分类和品牌关联表,这两个建立一个同名索引
type CategoryBrand struct {
	gorm.Model
	CategoryId uint     `gorm:"type:int;not null;comment:分类ID;index:idx_category_id,unique"`
	Category   Category `gorm:"foreignKey:CategoryId;references:ID"`
	BrandId    uint     `gorm:"type:int;not null;comment:品牌ID;index:idx_category_id,unique"`
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
	ClickNum        int        `gorm:"type:int;not null;default:0;comment:点击数"`
	FavNum          int        `gorm:"type:int;not null;default:0;comment:收藏数"`
	MarketPrice     float64    `gorm:"type:decimal(10,2);not null;comment:市场价格"`
	ShopPrice       float64    `gorm:"type:decimal(10,2);not null;comment:本店价格"`
	GoodsBrief      string     `gorm:"type:varchar(200);not null;comment:商品简介"`
	Images          GormList   `gorm:"type:json;not null;comment:商品图片"`
	DescImages      GormList   `gorm:"type:json;not null;comment:商品详情图片"`
	GoodsFrontImage string     `gorm:"type:varchar(200);not null;comment:商品主图"`
	Status          int        `gorm:"type:tinyint;not null;default:1;comment:商品状态"`
	Categories      []Category `gorm:"many2many:goods_category;"`
}

// Banner 轮播图
type Banner struct {
	gorm.Model
	Image string `gorm:"type:varchar(200);not null;comment:轮播图图片"`
	Url   string `gorm:"type:varchar(200);not null;comment:轮播图链接"`
	Index int    `gorm:"type:int;not null;default:0;comment:轮播图索引"`
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
	err := global.DB.Preload("SubCategories").First(&goods, id).Error
	if err != nil {
		return nil, err
	}
	return &goods, nil
}

// UpdateGoods 更新商品
func UpdateGoods(goods *Goods) error {
	return global.DB.Model(&Goods{}).
		Where("id = ?", goods.ID).
		Updates(map[string]interface{}{
			"brand_id":          goods.BrandId,
			"on_sale":           goods.OnSale,
			"ship_free":         goods.ShipFree,
			"is_new":            goods.IsNew,
			"is_hot":            goods.IsHot,
			"name":              goods.Name,
			"goods_sn":          goods.GoodsSn,
			"market_price":      goods.MarketPrice,
			"shop_price":        goods.ShopPrice,
			"goods_brief":       goods.GoodsBrief,
			"images":            goods.Images,
			"desc_images":       goods.DescImages,
			"goods_front_image": goods.GoodsFrontImage,
			"status":            goods.Status,
		}).Error
}

// DeleteGoods 删除商品
func DeleteGoods(id uint) error {
	return global.DB.Delete(&Goods{}, id).Error
}

// UpdateGoodsByMap 只更新指定字段
func UpdateGoodsByMap(id uint, updateMap map[string]interface{}) error {
	return global.DB.Model(&Goods{}).Where("id = ?", id).Updates(updateMap).Error
}

// GoodsFilter 商品查询过滤器
type GoodsFilter struct {
	Page       int
	PageSize   int
	MinPrice   float64
	MaxPrice   float64
	IsHot      *bool
	IsNew      *bool
	IsTab      *bool
	BrandId    uint
	CategoryId uint
	Keywords   string
	OnSale     *bool
}

// GetGoodsList 获取商品列表
func GetGoodsList(filter *GoodsFilter) ([]Goods, int64, error) {
	var goods []Goods
	var total int64

	// 构建查询
	query := global.DB.Model(&Goods{})

	// 应用过滤条件
	if filter.MinPrice > 0 {
		query = query.Where("shop_price >= ?", filter.MinPrice)
	}
	if filter.MaxPrice > 0 {
		query = query.Where("shop_price <= ?", filter.MaxPrice)
	}
	if filter.IsHot != nil {
		query = query.Where("is_hot = ?", *filter.IsHot)
	}
	if filter.IsNew != nil {
		query = query.Where("is_new = ?", *filter.IsNew)
	}
	if filter.IsTab != nil {
		query = query.Where("is_tab = ?", *filter.IsTab)
	}
	if filter.BrandId > 0 {
		query = query.Where("brand_id = ?", filter.BrandId)
	}
	if filter.CategoryId > 0 {
		query = query.Joins("JOIN goods_category ON goods.id = goods_category.goods_id").
			Where("goods_category.category_id = ?", filter.CategoryId)
	}
	if filter.Keywords != "" {
		query = query.Where("name LIKE ? OR goods_brief LIKE ?", "%"+filter.Keywords+"%", "%"+filter.Keywords+"%")
	}
	if filter.OnSale != nil {
		query = query.Where("on_sale = ?", *filter.OnSale)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (filter.Page - 1) * filter.PageSize
	err := query.Preload("Categories").
		Preload("Brand").
		Offset(offset).
		Limit(filter.PageSize).
		Find(&goods).Error

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

// GetBrandList 获取品牌列表
func GetBrandList(page, pageSize int, name, desc string) ([]Brand, int64, error) {
	var brands []Brand
	var total int64

	// 构建查询
	query := global.DB.Model(&Brand{})

	// 添加名称过滤条件
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	// 添加描述过滤条件
	if desc != "" {
		query = query.Where("desc LIKE ?", "%"+desc+"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&brands).Error; err != nil {
		return nil, 0, err
	}

	return brands, total, nil
}

// CreateBrand 创建品牌
func CreateBrand(brand *Brand) error {
	return global.DB.Create(brand).Error
}

// DeleteBrand 删除品牌
func DeleteBrand(id uint) error {
	return global.DB.Delete(&Brand{}, id).Error
}

// UpdateBrand 更新品牌
func UpdateBrand(brand *Brand) error {
	return global.DB.Save(brand).Error
}

// GetAllCategories 获取所有分类
func GetAllCategories() ([]Category, error) {
	var categories []Category
	if err := global.DB.Preload("SubCategories.SubCategories").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// GetSubCategory 获取子分类
func GetSubCategory(id int32, level int32) (*Category, []Category, error) {
	var category Category
	var subCategories []Category

	// 获取当前分类
	if err := global.DB.First(&category, id).Error; err != nil {
		return nil, nil, err
	}

	// 获取子分类
	if err := global.DB.Where("parent_id = ? AND level = ?", id, level).Find(&subCategories).Error; err != nil {
		return nil, nil, err
	}

	return &category, subCategories, nil
}

// CreateCategory 创建分类
func CreateCategory(category *Category) error {
	return global.DB.Create(category).Error
}

// DeleteCategory 删除分类
func DeleteCategory(id uint) error {
	return global.DB.Delete(&Category{}, id).Error
}

// UpdateCategory 更新分类
func UpdateCategory(category *Category) error {
	return global.DB.Save(category).Error
}

// GetBannerList 获取轮播图列表
func GetBannerList() ([]Banner, error) {
	var banners []Banner
	if err := global.DB.Order("`index` asc").Find(&banners).Error; err != nil {
		return nil, err
	}
	return banners, nil
}

// CreateBanner 创建轮播图
func CreateBanner(banner *Banner) error {
	return global.DB.Create(banner).Error
}

// DeleteBanner 删除轮播图
func DeleteBanner(id uint) error {
	return global.DB.Unscoped().Delete(&Banner{}, id).Error
}

// UpdateBanner 更新轮播图
func UpdateBanner(banner *Banner) error {
	return global.DB.Save(banner).Error
}

// CheckBrandNameExists 检查品牌名称是否存在
func CheckBrandNameExists(name string, excludeId ...uint) (bool, error) {
	var count int64
	query := global.DB.Model(&Brand{}).Where("name = ?", name)

	// 如果提供了排除ID，则排除该ID的品牌
	if len(excludeId) > 0 {
		query = query.Where("id != ?", excludeId[0])
	}

	err := query.Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
