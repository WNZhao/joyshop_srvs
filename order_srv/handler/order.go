package handler

import (
	context "context"
	"order_srv/proto"

	"google.golang.org/protobuf/types/known/emptypb"
)

type OrderServiceServer struct {
	proto.UnimplementedOrderServiceServer
}

// 购物车相关
func (s *OrderServiceServer) CartItemList(ctx context.Context, req *proto.UserInfo) (*proto.CartItemListResponse, error) {
	// TODO: 获取用户购物车信息
	return &proto.CartItemListResponse{}, nil
}

func (s *OrderServiceServer) CartItemAdd(ctx context.Context, req *proto.CartItemRequest) (*proto.ShopCartInfoResponse, error) {
	// TODO: 添加购物车
	return &proto.ShopCartInfoResponse{}, nil
}

func (s *OrderServiceServer) CartItemUpdate(ctx context.Context, req *proto.CartItemRequest) (*emptypb.Empty, error) {
	// TODO: 更新购物车
	return &emptypb.Empty{}, nil
}

func (s *OrderServiceServer) CartItemDelete(ctx context.Context, req *proto.CartItemRequest) (*emptypb.Empty, error) {
	// TODO: 删除购物车
	return &emptypb.Empty{}, nil
}

// 订单相关
func (s *OrderServiceServer) OrderCreate(ctx context.Context, req *proto.OrderRequest) (*proto.OrderInfoResponse, error) {
	// TODO: 创建订单
	return &proto.OrderInfoResponse{}, nil
}

func (s *OrderServiceServer) OrderList(ctx context.Context, req *proto.OrderFilterRequest) (*proto.OrderListResponse, error) {
	// TODO: 获取订单列表
	return &proto.OrderListResponse{}, nil
}

func (s *OrderServiceServer) OrderDetail(ctx context.Context, req *proto.OrderRequest) (*proto.OrderInfoDetailResponse, error) {
	// TODO: 获取订单详情
	return &proto.OrderInfoDetailResponse{}, nil
}

func (s *OrderServiceServer) OrderUpdate(ctx context.Context, req *proto.OrderStatus) (*emptypb.Empty, error) {
	// TODO: 更新订单状态
	return &emptypb.Empty{}, nil
}

func (s *OrderServiceServer) OrderDelete(ctx context.Context, req *proto.OrderDelRequest) (*emptypb.Empty, error) {
	// TODO: 删除订单
	return &emptypb.Empty{}, nil
}
