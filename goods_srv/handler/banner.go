/*
 * @Author: Will nanan_zhao@163.com
 * @Date: 2025-05-17 14:15:04
 * @LastEditors: Will nanan_zhao@163.com
 * @LastEditTime: 2025-05-18 16:58:20
 * @FilePath: /joyshop_srvs/goods_srv/handler/banner.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package handler

import (
	"context"
	"goods_srv/global"
	"goods_srv/model"
	"goods_srv/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
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
	// 检查轮播图是否存在
	var banner model.Banner
	if err := global.DB.First(&banner, req.Id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "轮播图不存在")
		}
		return nil, status.Error(codes.Internal, "查询轮播图失败")
	}

	// 更新轮播图信息
	banner.Image = req.Image
	banner.Url = req.Url
	banner.Index = int(req.Index)

	err := model.UpdateBanner(&banner)
	if err != nil {
		return nil, status.Error(codes.Internal, "更新轮播图失败")
	}

	return &emptypb.Empty{}, nil
}
