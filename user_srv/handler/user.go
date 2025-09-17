package handler

import (
	"context"
	"crypto/sha512"
	"fmt"
	"strings"
	"time"
	"user_srv/global"
	"user_srv/model"
	"user_srv/proto"

	"github.com/anaskhan96/go-password-encoder"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

// 实现接口
//GetUserList(context.Context, *PageInfo) (*UserListResponse, error)
//GetUserByMobile(context.Context, *MobileRequest) (*UserInfoResponse, error)
//GetUserById(context.Context, *IdRequest) (*UserInfoResponse, error)
//CreateUser(context.Context, *CreateUserInfo) (*UserInfoResponse, error)
//UpdateUser(context.Context, *UpdateUserInfo) (*emptypb.Empty, error)
//DeleteUser(context.Context, *IdRequest) (*emptypb.Empty, error)
//CheckPassword(context.Context, *PasswordCheckInof) (*CheckReponse, error)
//mustEmbedUnimplementedUserServer()

type UserServer struct {
	proto.UnimplementedUserServer
}

// Paginate 生成分页函数，供 GORM 使用
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// 获取页码

		if page <= 0 {
			page = 1
		}
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		// 计算偏移量并返回分页处理后的 DB
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func ModelToRsponse(user model.User) proto.UserInfoResponse {
	// 将 model.User 转换为 proto.UserInfoResponse
	// 在grpc的message中，字段有默认值，你不能随便赋值nil进去，容易出错
	rsp := proto.UserInfoResponse{
		Id:       user.ID,
		Password: user.Password,
		Mobile:   user.Mobile,
		Email:    user.Email,
		Nickname: user.NickName,
		Username: user.UserName,
		Gender:   user.Gender,
		Avatar:   user.Avatar,
		Role:     int32(user.Role),
	}
	if user.Birthday != nil {
		rsp.Birthday = uint64(user.Birthday.Unix())
	}
	return rsp
}

func (s *UserServer) GetUserList(ctx context.Context, info *proto.PageInfo) (*proto.UserListResponse, error) {
	// 获取用户列表
	var users []model.User
	result := global.DB.Table("user").Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	// 获取总数
	rsp := &proto.UserListResponse{}
	rsp.Total = uint32(result.RowsAffected)
	global.DB.Table("user").Scopes(Paginate(int(info.Page), int(info.PageSize))).Find(&users)

	for _, user := range users {
		userInfoRsp := ModelToRsponse(user)
		rsp.Data = append(rsp.Data, &userInfoRsp)
	}
	return rsp, nil
}

func (s *UserServer) GetUserByMobile(ctx context.Context, req *proto.MobileRequest) (*proto.UserInfoResponse, error) {
	// 通过手机号获取用户信息
	var user model.User
	result := global.DB.Where("mobile = ?", req.Mobile).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	// 将 model.User 转换为 proto.UserInfoResponse
	userInfoRsp := ModelToRsponse(user)
	return &userInfoRsp, nil
}

func (s *UserServer) GetUserById(ctx context.Context, req *proto.IdRequest) (*proto.UserInfoResponse, error) {
	// 通过手机号获取用户信息
	var user model.User
	result := global.DB.First(&user, req.Id)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	// 将 model.User 转换为 proto.UserInfoResponse
	userInfoRsp := ModelToRsponse(user)
	return &userInfoRsp, nil
}

func (s *UserServer) CreateUser(ctx context.Context, req *proto.CreateUserInfo) (*proto.UserInfoResponse, error) {
	// 创建用户
	var user model.User
	result := global.DB.Where(&model.User{Mobile: req.Mobile}).First(&user)
	if result.RowsAffected == 1 {
		return nil, status.Errorf(codes.AlreadyExists, "用户已存在")
	}
	// 创建用户
	user.Mobile = req.Mobile
	user.Email = req.Email
	user.NickName = req.Nickname
	user.UserName = req.Username

	options := &password.Options{10, 100, 32, sha512.New}
	salt, encodedPwd := password.Encode(req.Password, options)
	//把盐值和加密后的密码存储到数据库中 盐值整合到加密后的密码中
	newPassword := fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)
	//fmt.Println(newPassword)
	//fmt.Println(salt) // 8c7a9f3b
	user.Password = newPassword
	result = global.DB.Create(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	userInfoRsp := ModelToRsponse(user)
	return &userInfoRsp, nil
}

func (s *UserServer) UpdateUser(ctx context.Context, req *proto.UpdateUserInfo) (*emptypb.Empty, error) {
	// 更新用户
	var user model.User
	result := global.DB.First(&user, req.Id)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}

	// 只更新非空字段，支持部分字段更新
	updates := make(map[string]interface{})

	// 检查用户名和邮箱是否重复（如果要更新的话）
	if req.Username != "" {
		var existUser model.User
		result := global.DB.Where("user_name = ? AND id != ?", req.Username, req.Id).First(&existUser)
		if result.RowsAffected > 0 {
			return nil, status.Errorf(codes.AlreadyExists, "用户名已被使用")
		}
		updates["user_name"] = req.Username
	}

	if req.Email != "" {
		var existUser model.User
		result := global.DB.Where("email = ? AND id != ?", req.Email, req.Id).First(&existUser)
		if result.RowsAffected > 0 {
			return nil, status.Errorf(codes.AlreadyExists, "邮箱已被使用")
		}
		updates["email"] = req.Email
	}

	if req.Mobile != "" {
		var existUser model.User
		result := global.DB.Where("mobile = ? AND id != ?", req.Mobile, req.Id).First(&existUser)
		if result.RowsAffected > 0 {
			return nil, status.Errorf(codes.AlreadyExists, "手机号已被使用")
		}
		updates["mobile"] = req.Mobile
	}

	if req.Nickname != "" {
		updates["nick_name"] = req.Nickname
	}
	if req.Gender != "" {
		updates["gender"] = req.Gender
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}
	if req.Birthday != 0 {
		birthday := time.Unix(int64(req.Birthday), 0)
		updates["birthday"] = &birthday
	}
	if req.Password != "" {
		// 如果更新密码，需要重新加密
		options := &password.Options{10, 100, 32, sha512.New}
		salt, encodedPwd := password.Encode(req.Password, options)
		newPassword := fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)
		updates["password"] = newPassword
	}

	// 如果没有要更新的字段，直接返回
	if len(updates) == 0 {
		return &emptypb.Empty{}, nil
	}

	// 使用 Updates 方法进行部分更新
	result = global.DB.Model(&user).Updates(updates)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	zap.S().Infof("用户 ID=%d 更新成功，更新字段: %+v", req.Id, updates)
	return &emptypb.Empty{}, nil
}

func (s *UserServer) CheckPassword(ctx context.Context, req *proto.PasswordCheckInof) (*proto.CheckReponse, error) {
	// 检查密码
	passwordInfo := strings.Split(req.EncryptPassword, "$")
	
	// 检查数组长度，防止越界
	if len(passwordInfo) < 4 {
		zap.S().Errorf("密码格式错误: %s, 期望格式: $pbkdf2-sha512$salt$encodedPwd", req.EncryptPassword)
		return &proto.CheckReponse{
			Success: false,
		}, nil
	}
	
	options := &password.Options{10, 100, 32, sha512.New}
	check := password.Verify(req.Password, passwordInfo[2], passwordInfo[3], options)
	return &proto.CheckReponse{
		Success: check,
	}, nil
}

func (s *UserServer) DeleteUser(ctx context.Context, req *proto.IdRequest) (*emptypb.Empty, error) {
	// 删除用户
	var user model.User
	result := global.DB.First(&user, req.Id)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	result = global.DB.Delete(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	return &emptypb.Empty{}, nil
}
