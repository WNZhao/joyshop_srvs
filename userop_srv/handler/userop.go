package handler

import (
	"context"
	"userop_srv/global"
	"userop_srv/model"
	"userop_srv/proto"

	"gorm.io/gorm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserOpServer struct {
	proto.UnimplementedUserOpServer
}

// GetFavList 获取收藏列表
func (s *UserOpServer) GetFavList(ctx context.Context, req *proto.UserFavRequest) (*proto.FavListResponse, error) {
	var favs []model.UserFav
	var favList []*proto.GoodsInfo
	var count int64

	// 查询收藏记录
	result := global.DB.Where(&model.UserFav{User: req.UserId}).Count(&count)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "查询收藏失败")
	}

	global.DB.Where(&model.UserFav{User: req.UserId}).Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Find(&favs)

	for _, fav := range favs {
		favList = append(favList, &proto.GoodsInfo{
			Id: fav.Goods,
		})
	}

	return &proto.FavListResponse{
		Total: int32(count),
		Data:  favList,
	}, nil
}

// AddUserFav 添加收藏
func (s *UserOpServer) AddUserFav(ctx context.Context, req *proto.UserFavRequest) (*emptypb.Empty, error) {
	var userFav model.UserFav

	// 检查是否已经收藏
	if result := global.DB.Where(&model.UserFav{User: req.UserId, Goods: req.GoodsId}).First(&userFav); result.RowsAffected == 1 {
		return nil, status.Errorf(codes.AlreadyExists, "已经收藏过了")
	}

	userFav.User = req.UserId
	userFav.Goods = req.GoodsId

	if result := global.DB.Create(&userFav); result.Error != nil {
		return nil, status.Errorf(codes.Internal, "添加收藏失败")
	}

	return &emptypb.Empty{}, nil
}

// DeleteUserFav 删除收藏
func (s *UserOpServer) DeleteUserFav(ctx context.Context, req *proto.UserFavRequest) (*emptypb.Empty, error) {
	if result := global.DB.Where(&model.UserFav{User: req.UserId, Goods: req.GoodsId}).Delete(&model.UserFav{}); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "收藏记录不存在")
	}
	return &emptypb.Empty{}, nil
}

// GetUserFavDetail 获取收藏详情
func (s *UserOpServer) GetUserFavDetail(ctx context.Context, req *proto.UserFavRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

// IsFav 是否收藏
func (s *UserOpServer) IsFav(ctx context.Context, req *proto.UserFavRequest) (*proto.IsFavResponse, error) {
	var count int64
	global.DB.Model(&model.UserFav{}).Where(&model.UserFav{User: req.UserId, Goods: req.GoodsId}).Count(&count)
	return &proto.IsFavResponse{
		Success: count > 0,
	}, nil
}

// GetAddressList 获取地址列表
func (s *UserOpServer) GetAddressList(ctx context.Context, req *proto.AddressRequest) (*proto.AddressListResponse, error) {
	var addresses []model.Address
	var addressList []*proto.Address
	var count int64

	result := global.DB.Where(&model.Address{User: req.UserId}).Count(&count)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "查询地址失败")
	}

	global.DB.Where(&model.Address{User: req.UserId}).Find(&addresses)

	for _, address := range addresses {
		addressList = append(addressList, &proto.Address{
			Id:           int32(address.ID),
			UserId:       address.User,
			Province:     address.Province,
			City:         address.City,
			District:     address.District,
			Address:      address.Address,
			SignerName:   address.SignerName,
			SignerMobile: address.SignerMobile,
		})
	}

	return &proto.AddressListResponse{
		Total: int32(count),
		Data:  addressList,
	}, nil
}

// CreateAddress 创建地址
func (s *UserOpServer) CreateAddress(ctx context.Context, req *proto.Address) (*proto.AddressResponse, error) {
	var address model.Address

	address.User = req.UserId
	address.Province = req.Province
	address.City = req.City
	address.District = req.District
	address.Address = req.Address
	address.SignerName = req.SignerName
	address.SignerMobile = req.SignerMobile

	if result := global.DB.Create(&address); result.Error != nil {
		return nil, status.Errorf(codes.Internal, "创建地址失败")
	}

	return &proto.AddressResponse{
		Id: int32(address.ID),
	}, nil
}

// UpdateAddress 更新地址
func (s *UserOpServer) UpdateAddress(ctx context.Context, req *proto.Address) (*emptypb.Empty, error) {
	var address model.Address

	if result := global.DB.Where(&model.Address{User: req.UserId}).First(&address, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "地址不存在")
	}

	if req.Province != "" {
		address.Province = req.Province
	}
	if req.City != "" {
		address.City = req.City
	}
	if req.District != "" {
		address.District = req.District
	}
	if req.Address != "" {
		address.Address = req.Address
	}
	if req.SignerName != "" {
		address.SignerName = req.SignerName
	}
	if req.SignerMobile != "" {
		address.SignerMobile = req.SignerMobile
	}

	if result := global.DB.Save(&address); result.Error != nil {
		return nil, status.Errorf(codes.Internal, "更新地址失败")
	}

	return &emptypb.Empty{}, nil
}

// DeleteAddress 删除地址
func (s *UserOpServer) DeleteAddress(ctx context.Context, req *proto.AddressRequest) (*emptypb.Empty, error) {
	if result := global.DB.Where(&model.Address{User: req.UserId}).Delete(&model.Address{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "地址不存在")
	}
	return &emptypb.Empty{}, nil
}

// CreateMessage 创建留言
func (s *UserOpServer) CreateMessage(ctx context.Context, req *proto.MessageRequest) (*proto.MessageResponse, error) {
	var message model.LeavingMessage

	message.User = req.UserId
	message.MessageType = req.MessageType
	message.Subject = req.Subject
	message.Message = req.Message
	message.File = req.File

	if result := global.DB.Create(&message); result.Error != nil {
		return nil, status.Errorf(codes.Internal, "创建留言失败")
	}

	return &proto.MessageResponse{
		Id: int32(message.ID),
	}, nil
}

// MessageList 留言列表
func (s *UserOpServer) MessageList(ctx context.Context, req *proto.MessageRequest) (*proto.MessageListResponse, error) {
	var messages []model.LeavingMessage
	var messageList []*proto.MessageInfo
	var count int64

	result := global.DB.Where(&model.LeavingMessage{User: req.UserId}).Count(&count)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "查询留言失败")
	}

	global.DB.Where(&model.LeavingMessage{User: req.UserId}).Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Find(&messages)

	for _, message := range messages {
		messageList = append(messageList, &proto.MessageInfo{
			Id:          int32(message.ID),
			UserId:      message.User,
			MessageType: message.MessageType,
			Subject:     message.Subject,
			Message:     message.Message,
			File:        message.File,
		})
	}

	return &proto.MessageListResponse{
		Total: int32(count),
		Data:  messageList,
	}, nil
}

// Paginate 分页
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}