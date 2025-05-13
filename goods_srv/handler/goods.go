package handler

import (
	"context"
	"goods_srv/model"
	"goods_srv/proto"
	"time"
)

type GoodsServer struct {
	proto.UnimplementedGoodsServer
}

// GetGoodsList 获取商品列表
func (s *GoodsServer) GetGoodsList(ctx context.Context, req *proto.GoodsListRequest) (*proto.GoodsListResponse, error) {
	goods, total, err := model.GetGoodsList(req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	var goodsList []*proto.GoodsInfo
	for _, g := range goods {
		goodsList = append(goodsList, &proto.GoodsInfo{
			Id:          g.ID,
			Name:        g.Name,
			Description: g.Description,
			Price:       g.Price,
			Stock:       g.Stock,
			Category:    g.Category,
			Image:       g.Image,
			Status:      g.Status,
			CreatedAt:   g.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   g.UpdatedAt.Format(time.RFC3339),
		})
	}

	return &proto.GoodsListResponse{
		Total: int32(total),
		Data:  goodsList,
	}, nil
}

// GetGoodsById 获取商品详情
func (s *GoodsServer) GetGoodsById(ctx context.Context, req *proto.GoodsByIdRequest) (*proto.GoodsInfoResponse, error) {
	goods, err := model.GetGoodsById(req.Id)
	if err != nil {
		return nil, err
	}

	return &proto.GoodsInfoResponse{
		Data: &proto.GoodsInfo{
			Id:          goods.ID,
			Name:        goods.Name,
			Description: goods.Description,
			Price:       goods.Price,
			Stock:       goods.Stock,
			Category:    goods.Category,
			Image:       goods.Image,
			Status:      goods.Status,
			CreatedAt:   goods.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   goods.UpdatedAt.Format(time.RFC3339),
		},
	}, nil
}

// CreateGoods 创建商品
func (s *GoodsServer) CreateGoods(ctx context.Context, req *proto.CreateGoodsRequest) (*proto.GoodsInfoResponse, error) {
	goods := &model.Goods{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Category:    req.Category,
		Image:       req.Image,
		Status:      1,
	}

	err := model.CreateGoods(goods)
	if err != nil {
		return nil, err
	}

	return &proto.GoodsInfoResponse{
		Data: &proto.GoodsInfo{
			Id:          goods.ID,
			Name:        goods.Name,
			Description: goods.Description,
			Price:       goods.Price,
			Stock:       goods.Stock,
			Category:    goods.Category,
			Image:       goods.Image,
			Status:      goods.Status,
			CreatedAt:   goods.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   goods.UpdatedAt.Format(time.RFC3339),
		},
	}, nil
}

// UpdateGoods 更新商品
func (s *GoodsServer) UpdateGoods(ctx context.Context, req *proto.UpdateGoodsRequest) (*proto.GoodsInfoResponse, error) {
	goods := &model.Goods{
		ID:          req.Id,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Category:    req.Category,
		Image:       req.Image,
		Status:      req.Status,
	}

	err := model.UpdateGoods(goods)
	if err != nil {
		return nil, err
	}

	return &proto.GoodsInfoResponse{
		Data: &proto.GoodsInfo{
			Id:          goods.ID,
			Name:        goods.Name,
			Description: goods.Description,
			Price:       goods.Price,
			Stock:       goods.Stock,
			Category:    goods.Category,
			Image:       goods.Image,
			Status:      goods.Status,
			CreatedAt:   goods.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   goods.UpdatedAt.Format(time.RFC3339),
		},
	}, nil
}

// DeleteGoods 删除商品
func (s *GoodsServer) DeleteGoods(ctx context.Context, req *proto.DeleteGoodsRequest) (*proto.Empty, error) {
	err := model.DeleteGoods(req.Id)
	if err != nil {
		return nil, err
	}
	return &proto.Empty{}, nil
}

// GetGoodsByCategory 按分类获取商品
func (s *GoodsServer) GetGoodsByCategory(ctx context.Context, req *proto.GoodsByCategoryRequest) (*proto.GoodsListResponse, error) {
	goods, total, err := model.GetGoodsByCategory(req.Category, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	var goodsList []*proto.GoodsInfo
	for _, g := range goods {
		goodsList = append(goodsList, &proto.GoodsInfo{
			Id:          g.ID,
			Name:        g.Name,
			Description: g.Description,
			Price:       g.Price,
			Stock:       g.Stock,
			Category:    g.Category,
			Image:       g.Image,
			Status:      g.Status,
			CreatedAt:   g.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   g.UpdatedAt.Format(time.RFC3339),
		})
	}

	return &proto.GoodsListResponse{
		Total: int32(total),
		Data:  goodsList,
	}, nil
}

// SearchGoods 搜索商品
func (s *GoodsServer) SearchGoods(ctx context.Context, req *proto.SearchGoodsRequest) (*proto.GoodsListResponse, error) {
	goods, total, err := model.SearchGoods(req.Keyword, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	var goodsList []*proto.GoodsInfo
	for _, g := range goods {
		goodsList = append(goodsList, &proto.GoodsInfo{
			Id:          g.ID,
			Name:        g.Name,
			Description: g.Description,
			Price:       g.Price,
			Stock:       g.Stock,
			Category:    g.Category,
			Image:       g.Image,
			Status:      g.Status,
			CreatedAt:   g.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   g.UpdatedAt.Format(time.RFC3339),
		})
	}

	return &proto.GoodsListResponse{
		Total: int32(total),
		Data:  goodsList,
	}, nil
}
