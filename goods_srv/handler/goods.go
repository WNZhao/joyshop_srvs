package handler

import (
	"context"
	"goods_srv/model"
	"goods_srv/proto"

	"google.golang.org/protobuf/types/known/emptypb"
)

type GoodsServer struct {
	proto.UnimplementedGoodsServer
}

// GetGoodsList 获取商品列表
func (s *GoodsServer) GetGoodsList(ctx context.Context, req *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error) {
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
func (s *GoodsServer) GetGoodsById(ctx context.Context, req *proto.GoodInfoRequest) (*proto.GoodsInfoResponse, error) {
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
func (s *GoodsServer) UpdateGoods(ctx context.Context, req *proto.CreateGoodsInfo) (*emptypb.Empty, error) {
	goods := ProtoToModelGoodsUpdate(req)
	err := model.UpdateGoods(goods)
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
