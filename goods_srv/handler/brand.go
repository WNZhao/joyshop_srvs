/*
 * @Author: Will nanan_zhao@163.com
 * @Date: 2025-05-17 14:14:53
 * @LastEditors: Will nanan_zhao@163.com
 * @LastEditTime: 2025-05-18 15:41:58
 * @FilePath: /joyshop_srvs/goods_srv/handler/brand.go
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

// BrandList 获取品牌列表
func (s *GoodsServer) BrandList(ctx context.Context, req *proto.BrandFilterRequest) (*proto.BrandListResponse, error) {
	brands, total, err := model.GetBrandList(int(req.Pages), int(req.PagePerNums), req.Name, req.Desc)
	if err != nil {
		return nil, err
	}

	var brandList []*proto.BrandInfoResponse
	for _, b := range brands {
		brandList = append(brandList, ModelToProtoBrand(&b))
	}

	return &proto.BrandListResponse{
		Total: int32(total),
		Data:  brandList,
	}, nil
}

// CreateBrand 创建品牌
func (s *GoodsServer) CreateBrand(ctx context.Context, req *proto.BrandRequest) (*proto.BrandInfoResponse, error) {
	// 检查品牌名称是否已存在
	exists, err := model.CheckBrandNameExists(req.Name)
	if err != nil {
		return nil, status.Error(codes.Internal, "检查品牌名称失败")
	}
	if exists {
		return nil, status.Error(codes.InvalidArgument, "品牌名称已存在")
	}

	brand := ProtoToModelBrand(req)
	err = model.CreateBrand(brand)
	if err != nil {
		return nil, status.Error(codes.Internal, "创建品牌失败")
	}

	return ModelToProtoBrand(brand), nil
}

// DeleteBrand 删除品牌
func (s *GoodsServer) DeleteBrand(ctx context.Context, req *proto.BrandRequest) (*emptypb.Empty, error) {
	err := model.DeleteBrand(uint(req.Id))
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// UpdateBrand 更新品牌
func (s *GoodsServer) UpdateBrand(ctx context.Context, req *proto.BrandRequest) (*emptypb.Empty, error) {
	// 检查品牌是否存在
	var brand model.Brand
	if err := global.DB.First(&brand, req.Id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "品牌不存在")
		}
		return nil, status.Error(codes.Internal, "查询品牌失败")
	}

	// 检查品牌名称是否已存在（排除当前品牌）
	exists, err := model.CheckBrandNameExists(req.Name, uint(req.Id))
	if err != nil {
		return nil, status.Error(codes.Internal, "检查品牌名称失败")
	}
	if exists {
		return nil, status.Error(codes.InvalidArgument, "品牌名称已存在")
	}

	// 更新品牌信息
	brand = *ProtoToModelBrand(req)
	err = model.UpdateBrand(&brand)
	if err != nil {
		return nil, status.Error(codes.Internal, "更新品牌失败")
	}

	return &emptypb.Empty{}, nil
}
