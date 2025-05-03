package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"joyshop_srvs/user_srv/proto"
)

var userClient proto.UserClient
var conn *grpc.ClientConn

func Init() {
	var err error
	conn, err = grpc.NewClient("localhost:50051",
		//grpc.WithInsecure()
		grpc.WithTransportCredentials(insecure.NewCredentials()), // ✅ 新方式
	)
	if err != nil {
		panic(err)
	}
	userClient = proto.NewUserClient(conn)
}

func TestGetUserList() {
	rsp, err := userClient.GetUserList(context.Background(), &proto.PageInfo{Page: 1, PageSize: 2})
	if err != nil {
		panic(err)
	}
	for _, user := range rsp.Data {
		println(user.Id, user.Mobile, user.Nickname, user.Password)
		rsp, err := userClient.CheckPassword(context.Background(), &proto.PasswordCheckInof{
			Password:        "admin123",
			EncryptPassword: user.Password,
		})
		if err != nil {
			panic(err)
		}
		println("密码是否正确：", rsp.Success)
	}
}

func main() {
	Init()
	defer conn.Close()

	TestGetUserList()
}
