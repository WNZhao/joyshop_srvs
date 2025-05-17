package handler

import (
	"context"
	"goods_srv/model"
	"goods_srv/proto"

	"google.golang.org/protobuf/types/known/emptypb"
)

// BannerList 获取轮播图列表
func (s *GoodsServer) BannerList(ctx context.Context, _ *emptypb.Empty) (*proto.BannerListResponse, error) {
	banners, err := model.GetBannerList()
	if err != nil {
		return nil, err
	}

	var bannerList []*proto.BannerResponse
	for _, b := range banners {
		bannerList = append(bannerList, ModelToProtoBanner(&b))
	}

	return &proto.BannerListResponse{
		Total: int32(len(bannerList)),
		Data:  bannerList,
	}, nil
}

// CreateBanner 创建轮播图
func (s *GoodsServer) CreateBanner(ctx context.Context, req *proto.BannerRequest) (*proto.BannerResponse, error) {
	banner := ProtoToModelBanner(req)
	err := model.CreateBanner(banner)
	if err != nil {
		return nil, err
	}

	return ModelToProtoBanner(banner), nil
}

// DeleteBanner 删除轮播图
func (s *GoodsServer) DeleteBanner(ctx context.Context, req *proto.BannerRequest) (*emptypb.Empty, error) {
	err := model.DeleteBanner(uint(req.Id))
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// UpdateBanner 更新轮播图
func (s *GoodsServer) UpdateBanner(ctx context.Context, req *proto.BannerRequest) (*emptypb.Empty, error) {
	banner := ProtoToModelBanner(req)
	err := model.UpdateBanner(banner)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
