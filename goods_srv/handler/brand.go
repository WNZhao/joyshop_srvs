/*
 * @Author: Will nanan_zhao@163.com
 * @Date: 2025-05-17 14:14:53
 * @LastEditors: Will nanan_zhao@163.com
 * @LastEditTime: 2025-05-17 14:19:48
 * @FilePath: /joyshop_srvs/goods_srv/handler/brand.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package handler

import (
	"context"
	"goods_srv/model"
	"goods_srv/proto"

	"google.golang.org/protobuf/types/known/emptypb"
)

// BrandList 获取品牌列表
func (s *GoodsServer) BrandList(ctx context.Context, req *proto.BrandFilterRequest) (*proto.BrandListResponse, error) {
	brands, total, err := model.GetBrandList(int(req.Pages), int(req.PagePerNums))
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
	brand := ProtoToModelBrand(req)
	err := model.CreateBrand(brand)
	if err != nil {
		return nil, err
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
	brand := ProtoToModelBrand(req)
	err := model.UpdateBrand(brand)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
