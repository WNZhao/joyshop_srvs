package handler

import (
	"goods_srv/model"
	"goods_srv/proto"
)

// ModelToProtoGoods 将model.Goods转换为proto.GoodsInfoResponse
func ModelToProtoGoods(g *model.Goods) *proto.GoodsInfoResponse {
	return &proto.GoodsInfoResponse{
		Id:              int32(g.ID),
		Name:            g.Name,
		GoodsSn:         g.GoodsSn,
		Stocks:          0, // 需要从库存服务获取
		MarketPrice:     float32(g.MarketPrice),
		ShopPrice:       float32(g.ShopPrice),
		GoodsBrief:      g.GoodsBrief,
		GoodsDesc:       "", // 需要从商品详情服务获取
		Images:          g.Images,
		DescImages:      g.DescImages,
		GoodsFrontImage: g.GoodsFrontImage,
		Status:          int32(g.Status),
		IsHot:           g.IsHot,
		IsNew:           g.IsNew,
		OnSale:          g.OnSale,
		ShipFree:        g.ShipFree,
		BrandId:         int32(g.BrandId),
		CategoryIds:     nil, // 需要从关联表获取
	}
}

// ProtoToModelGoods 将proto.CreateGoodsInfo转换为model.Goods
func ProtoToModelGoods(req *proto.CreateGoodsInfo) *model.Goods {
	return &model.Goods{
		BrandId:         uint(req.BrandId),
		OnSale:          req.OnSale,
		ShipFree:        req.ShipFree,
		IsNew:           req.IsNew,
		IsHot:           req.IsHot,
		Name:            req.Name,
		GoodsSn:         req.GoodsSn,
		MarketPrice:     float64(req.MarketPrice),
		ShopPrice:       float64(req.ShopPrice),
		GoodsBrief:      req.GoodsBrief,
		Images:          req.Images,
		DescImages:      req.DescImages,
		GoodsFrontImage: req.GoodsFrontImage,
		Status:          int(req.Status),
	}
}

// ProtoToModelGoodsUpdate 将proto.CreateGoodsInfo转换为model.Goods
func ProtoToModelGoodsUpdate(req *proto.CreateGoodsInfo) *model.Goods {
	goods := &model.Goods{
		BrandId:         uint(req.BrandId),
		OnSale:          req.OnSale,
		ShipFree:        req.ShipFree,
		IsNew:           req.IsNew,
		IsHot:           req.IsHot,
		Name:            req.Name,
		GoodsSn:         req.GoodsSn,
		MarketPrice:     float64(req.MarketPrice),
		ShopPrice:       float64(req.ShopPrice),
		GoodsBrief:      req.GoodsBrief,
		Images:          req.Images,
		DescImages:      req.DescImages,
		GoodsFrontImage: req.GoodsFrontImage,
		Status:          int(req.Status),
	}
	goods.ID = uint(req.Id)
	return goods
}

// ProtoToModelFilter 将proto.GoodsFilterRequest转换为model.GoodsFilter
func ProtoToModelFilter(req *proto.GoodsFilterRequest) *model.GoodsFilter {
	filter := &model.GoodsFilter{
		Page:     int(req.Pages),
		PageSize: int(req.PagePerNums),
	}

	// 商品属性
	if req.IsHot {
		filter.IsHot = &req.IsHot
	}
	if req.IsNew {
		filter.IsNew = &req.IsNew
	}
	if req.IsTab {
		filter.IsTab = &req.IsTab
	}
	if req.OnSale {
		filter.OnSale = &req.OnSale
	}

	// 品牌和分类
	if req.BrandId > 0 {
		filter.BrandId = uint(req.BrandId)
	}
	if req.CategoryId > 0 {
		filter.CategoryId = uint(req.CategoryId)
	}

	// 关键词搜索
	if req.Keywords != "" {
		filter.Keywords = req.Keywords
	}

	return filter
}

// ModelToProtoCategory 将model.Category转换为proto.CategoryInfoResponse
func ModelToProtoCategory(c *model.Category) *proto.CategoryInfoResponse {
	var parentId int32
	if c.ParentId != nil {
		parentId = int32(*c.ParentId)
	}
	return &proto.CategoryInfoResponse{
		Id:       int32(c.ID),
		Name:     c.Name,
		ParentId: parentId,
		Level:    int32(c.Level),
		IsTab:    c.IsTab,
	}
}

// ProtoToModelCategory 将proto.CategoryInfoRequest转换为model.Category
func ProtoToModelCategory(req *proto.CategoryInfoRequest) *model.Category {
	var parentId *int
	if req.ParentId != 0 {
		tmp := int(req.ParentId)
		parentId = &tmp
	}
	return &model.Category{
		Name:     req.Name,
		ParentId: parentId,
		Level:    int(req.Level),
		IsTab:    req.IsTab,
	}
}

// ModelToProtoBrand 将model.Brand转换为proto.BrandInfoResponse
func ModelToProtoBrand(b *model.Brand) *proto.BrandInfoResponse {
	return &proto.BrandInfoResponse{
		Id:   int32(b.ID),
		Name: b.Name,
		Logo: b.Logo,
		Desc: b.Desc,
	}
}

// ProtoToModelBrand 将proto.BrandRequest转换为model.Brand
func ProtoToModelBrand(req *proto.BrandRequest) *model.Brand {
	return &model.Brand{
		Name: req.Name,
		Logo: req.Logo,
		Desc: req.Desc,
	}
}

// ModelToProtoBanner 将model.Banner转换为proto.BannerResponse
func ModelToProtoBanner(b *model.Banner) *proto.BannerResponse {
	return &proto.BannerResponse{
		Id:    int32(b.ID),
		Image: b.Image,
		Url:   b.Url,
		Index: int32(b.Index),
	}
}

// ProtoToModelBanner 将proto.BannerRequest转换为model.Banner
func ProtoToModelBanner(req *proto.BannerRequest) *model.Banner {
	return &model.Banner{
		Image: req.Image,
		Url:   req.Url,
		Index: int(req.Index),
	}
}

// ModelToProtoCategoryBrand 将model.CategoryBrand转换为proto.CategoryBrandResponse
func ModelToProtoCategoryBrand(cb *model.CategoryBrand) *proto.CategoryBrandResponse {
	return &proto.CategoryBrandResponse{
		Id:         int32(cb.ID),
		CategoryId: int32(cb.CategoryId),
		BrandId:    int32(cb.BrandId),
	}
}

// ProtoToModelCategoryBrand 将proto.CategoryBrandRequest转换为model.CategoryBrand
func ProtoToModelCategoryBrand(req *proto.CategoryBrandRequest) *model.CategoryBrand {
	return &model.CategoryBrand{
		CategoryId: uint(req.CategoryId),
		BrandId:    uint(req.BrandId),
	}
}
