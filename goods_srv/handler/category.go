package handler

import (
	"context"
	"goods_srv/model"
	"goods_srv/proto"

	"google.golang.org/protobuf/types/known/emptypb"
)

// GetAllCategoriesList 获取所有分类
func (s *GoodsServer) GetAllCategoriesList(ctx context.Context, _ *emptypb.Empty) (*proto.CategoryListResponse, error) {
	categories, err := model.GetAllCategories()
	if err != nil {
		return nil, err
	}

	var categoryList []*proto.CategoryInfoResponse
	for _, c := range categories {
		categoryList = append(categoryList, ModelToProtoCategory(&c))
	}

	return &proto.CategoryListResponse{
		Total:    int32(len(categoryList)),
		Data:     categoryList,
		JsonData: "", // TODO: 如果需要JSON格式数据，需要实现转换
	}, nil
}

// GetSubCategory 获取子分类
func (s *GoodsServer) GetSubCategory(ctx context.Context, req *proto.CategoryListRequest) (*proto.SubCategoryListResponse, error) {
	category, subCategories, err := model.GetSubCategory(req.Id, req.Level)
	if err != nil {
		return nil, err
	}

	var subCategoryList []*proto.CategoryInfoResponse
	for _, c := range subCategories {
		subCategoryList = append(subCategoryList, ModelToProtoCategory(&c))
	}

	return &proto.SubCategoryListResponse{
		Total:         int32(len(subCategoryList)),
		Info:          ModelToProtoCategory(category),
		SubCategories: subCategoryList,
	}, nil
}

// CreateCategory 创建分类
func (s *GoodsServer) CreateCategory(ctx context.Context, req *proto.CategoryInfoRequest) (*proto.CategoryInfoResponse, error) {
	category := ProtoToModelCategory(req)
	err := model.CreateCategory(category)
	if err != nil {
		return nil, err
	}

	return ModelToProtoCategory(category), nil
}

// DeleteCategory 删除分类
func (s *GoodsServer) DeleteCategory(ctx context.Context, req *proto.DeleteCategoryRequest) (*emptypb.Empty, error) {
	err := model.DeleteCategory(uint(req.Id))
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// UpdateCategory 更新分类
func (s *GoodsServer) UpdateCategory(ctx context.Context, req *proto.CategoryInfoRequest) (*emptypb.Empty, error) {
	category := ProtoToModelCategory(req)
	err := model.UpdateCategory(category)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
