package handler

import (
	"context"
	"encoding/json"
	"goods_srv/global"
	"goods_srv/model"
	"goods_srv/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GoodsServer struct {
	proto.UnimplementedGoodsServer
}

// GetGoodsList 获取商品列表
func (s *GoodsServer) GoodsList(ctx context.Context, req *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error) {
	// 转换查询参数为过滤器
	filter := ProtoToModelFilter(req)

	// 使用过滤器查询商品
	goods, total, err := model.GetGoodsList(filter)
	if err != nil {
		return nil, err
	}

	var goodsList []*proto.GoodsInfoResponse
	for _, g := range goods {
		goodsList = append(goodsList, ModelToProtoGoods(&g))
	}

	return &proto.GoodsListResponse{
		Total: int32(total),
		Data:  goodsList,
	}, nil
}

// GetGoodsById 获取商品详情
func (s *GoodsServer) GetGoodsDetail(ctx context.Context, req *proto.GoodInfoRequest) (*proto.GoodsInfoResponse, error) {
	goods, err := model.GetGoodsById(uint(req.Id))
	if err != nil {
		return nil, err
	}

	return ModelToProtoGoods(goods), nil
}

// CreateGoods 创建商品
func (s *GoodsServer) CreateGoods(ctx context.Context, req *proto.CreateGoodsInfo) (*proto.GoodsInfoResponse, error) {
	goods := ProtoToModelGoods(req)
	err := model.CreateGoods(goods)
	if err != nil {
		return nil, err
	}

	return ModelToProtoGoods(goods), nil
}

// UpdateGoods 更新商品
// func (s *GoodsServer) UpdateGoods(ctx context.Context, req *proto.CreateGoodsInfo) (*emptypb.Empty, error) {
// 	goods := ProtoToModelGoodsUpdate(req)
// 	err := model.UpdateGoods(goods)
// 	if err != nil {
// 		return nil, err
// 	}

//		return &emptypb.Empty{}, nil
//	}
func (s *GoodsServer) UpdateGoods(ctx context.Context, req *proto.CreateGoodsInfo) (*emptypb.Empty, error) {
	imagesJson, _ := json.Marshal(req.Images)
	descImagesJson, _ := json.Marshal(req.DescImages)
	updateMap := map[string]interface{}{
		"brand_id":          req.BrandId,
		"on_sale":           req.OnSale,
		"ship_free":         req.ShipFree,
		"is_new":            req.IsNew,
		"is_hot":            req.IsHot,
		"name":              req.Name,
		"goods_sn":          req.GoodsSn,
		"market_price":      req.MarketPrice,
		"shop_price":        req.ShopPrice,
		"goods_brief":       req.GoodsBrief,
		"images":            string(imagesJson),
		"desc_images":       string(descImagesJson),
		"goods_front_image": req.GoodsFrontImage,
		"status":            req.Status,
	}
	err := model.UpdateGoodsByMap(uint(req.Id), updateMap)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// DeleteGoods 删除商品
func (s *GoodsServer) DeleteGoods(ctx context.Context, req *proto.DeleteGoodsInfo) (*emptypb.Empty, error) {
	err := model.DeleteGoods(uint(req.Id))
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// BatchGetGoods 批量获取商品信息
func (s *GoodsServer) BatchGetGoods(ctx context.Context, req *proto.BatchGoodsIdInfo) (*proto.GoodsListResponse, error) {
	global.Logger.Infof("批量获取商品信息，商品ID列表: %v", req.Id)

	if len(req.Id) == 0 {
		return &proto.GoodsListResponse{
			Total: 0,
			Data:  []*proto.GoodsInfoResponse{},
		}, nil
	}

	// 从数据库批量查询商品
	var goods []model.Goods
	if err := global.DB.Where("id IN ?", req.Id).Find(&goods).Error; err != nil {
		global.Logger.Errorf("批量查询商品失败: %v", err)
		return nil, status.Errorf(codes.Internal, "批量查询商品失败")
	}

	// 转换为proto格式并设置库存信息
	var goodsList []*proto.GoodsInfoResponse
	for _, g := range goods {
		goodsProto := ModelToProtoGoods(&g)
		// 从库存表获取库存信息
		goodsProto.Stocks = s.getGoodsStockFromDB(int32(g.ID))
		goodsList = append(goodsList, goodsProto)
	}

	global.Logger.Infof("成功批量获取商品信息，返回%d个商品", len(goodsList))
	return &proto.GoodsListResponse{
		Total: int32(len(goodsList)),
		Data:  goodsList,
	}, nil
}

// getGoodsStockFromDB 从库存数据库获取商品库存
func (s *GoodsServer) getGoodsStockFromDB(goodsId int32) int32 {
	// 查询库存表
	type Inventory struct {
		Stock int32 `gorm:"column:stock"`
	}
	var inv Inventory
	
	// 直接查询库存数据库（这里简化处理，实际应该通过库存服务）
	err := global.DB.Raw("SELECT stock FROM joyshop_inventory.inventory WHERE goods_id = ? AND deleted_at IS NULL", goodsId).Scan(&inv).Error
	if err != nil {
		return 0 // 查询失败返回0库存
	}
	return inv.Stock
}

// GetGoodsByCategory 按分类获取商品
func (s *GoodsServer) GetGoodsByCategory(ctx context.Context, req *proto.CategoryListRequest) (*proto.GoodsListResponse, error) {
	goods, total, err := model.GetGoodsByCategory(req.Id, req.Level, 1) // 使用默认分页大小
	if err != nil {
		return nil, err
	}

	var goodsList []*proto.GoodsInfoResponse
	for _, g := range goods {
		goodsList = append(goodsList, ModelToProtoGoods(&g))
	}

	return &proto.GoodsListResponse{
		Total: int32(total),
		Data:  goodsList,
	}, nil
}

// SearchGoods 搜索商品
func (s *GoodsServer) SearchGoods(ctx context.Context, req *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error) {
	goods, total, err := model.SearchGoods(req.Keywords, req.Pages, req.PagePerNums)
	if err != nil {
		return nil, err
	}

	var goodsList []*proto.GoodsInfoResponse
	for _, g := range goods {
		goodsList = append(goodsList, ModelToProtoGoods(&g))
	}

	return &proto.GoodsListResponse{
		Total: int32(total),
		Data:  goodsList,
	}, nil
}
