package main

import (
	"context"
	"joyshop_srvs/user_srv/proto"
	"math/rand"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

func TestDeleteUser() {
	// 先创建一个用户用于测试删除
	createReq := &proto.CreateUserInfo{
		Password: "test123",
		Mobile:   generateRandomMobile(),
		Email:    generateRandomEmail(),
		Nickname: generateRandomNickname(),
		Username: generateRandomUsername(),
		Birthday: uint64(time.Now().Unix()),
	}
	createRsp, err := userClient.CreateUser(context.Background(), createReq)
	if err != nil {
		panic(err)
	}
	println("创建用户成功，ID:", createRsp.Id)

	// 删除用户
	_, err = userClient.DeleteUser(context.Background(), &proto.IdRequest{Id: createRsp.Id})
	if err != nil {
		panic(err)
	}
	println("删除用户成功，ID:", createRsp.Id)

	// 验证用户是否真的被删除
	_, err = userClient.GetUserById(context.Background(), &proto.IdRequest{Id: createRsp.Id})
	if err == nil {
		panic("用户应该已经被删除")
	}
	println("验证用户已删除成功")
}

// 生成随机手机号
func generateRandomMobile() string {
	// 生成11位随机数字
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	// 确保手机号以1开头
	mobile := "1"
	for i := 0; i < 10; i++ {
		mobile += strconv.Itoa(r.Intn(10))
	}
	return mobile
}

// 生成随机邮箱
func generateRandomEmail() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	// 生成随机字符串作为邮箱前缀
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	prefix := make([]byte, 8)
	for i := range prefix {
		prefix[i] = letters[r.Intn(len(letters))]
	}
	return string(prefix) + "@example.com"
}

// 生成随机用户名
func generateRandomUsername() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	username := make([]byte, 8)
	for i := range username {
		username[i] = letters[r.Intn(len(letters))]
	}
	return "user_" + string(username)
}

// 生成随机昵称
func generateRandomNickname() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	adjectives := []string{"快乐的", "聪明的", "勇敢的", "善良的", "可爱的", "活泼的", "温柔的", "优雅的"}
	nouns := []string{"小猫", "小狗", "小兔", "小熊", "小鹿", "小象", "小马", "小羊"}
	return adjectives[r.Intn(len(adjectives))] + nouns[r.Intn(len(nouns))]
}

func TestCreateUser() {
	// 创建用户
	createReq := &proto.CreateUserInfo{
		Password: "test123",
		Mobile:   generateRandomMobile(),
		Email:    generateRandomEmail(),
		Nickname: generateRandomNickname(),
		Username: generateRandomUsername(),
		Birthday: uint64(time.Now().Unix()),
	}
	createRsp, err := userClient.CreateUser(context.Background(), createReq)
	if err != nil {
		panic(err)
	}
	println("创建用户成功，ID:", createRsp.Id)
	println("用户信息：", createRsp.Mobile, createRsp.Email, createRsp.Nickname, createRsp.Username)

	// 验证密码
	checkRsp, err := userClient.CheckPassword(context.Background(), &proto.PasswordCheckInof{
		Password:        "test123",
		EncryptPassword: createRsp.Password,
	})
	if err != nil {
		panic(err)
	}
	println("密码验证结果：", checkRsp.Success)
}

func TestBatchCreateUsers() {
	// 批量创建5个用户
	for i := 0; i < 5; i++ {
		createReq := &proto.CreateUserInfo{
			Password: "test123",
			Mobile:   generateRandomMobile(),
			Email:    generateRandomEmail(),
			Nickname: generateRandomNickname(),
			Username: generateRandomUsername(),
			Birthday: uint64(time.Now().Unix()),
		}
		createRsp, err := userClient.CreateUser(context.Background(), createReq)
		if err != nil {
			panic(err)
		}
		println("创建用户成功，ID:", createRsp.Id)
		println("用户信息：", createRsp.Mobile, createRsp.Email, createRsp.Nickname, createRsp.Username)
	}
}

func TestUpdateUser() {
	// 先创建一个用户用于测试更新
	createReq := &proto.CreateUserInfo{
		Password: "test123",
		Mobile:   generateRandomMobile(),
		Email:    generateRandomEmail(),
		Nickname: generateRandomNickname(),
		Username: generateRandomUsername(),
		Birthday: uint64(time.Now().Unix()),
	}
	createRsp, err := userClient.CreateUser(context.Background(), createReq)
	if err != nil {
		panic(err)
	}
	println("创建用户成功，ID:", createRsp.Id)

	// 更新用户信息
	updateReq := &proto.UpdateUserInfo{
		Id:       createRsp.Id,
		Nickname: generateRandomNickname(),
		Email:    generateRandomEmail(),
		Birthday: uint64(time.Now().AddDate(0, 0, 1).Unix()),
	}
	_, err = userClient.UpdateUser(context.Background(), updateReq)
	if err != nil {
		panic(err)
	}
	println("更新用户成功")

	// 验证更新结果
	getRsp, err := userClient.GetUserById(context.Background(), &proto.IdRequest{Id: createRsp.Id})
	if err != nil {
		panic(err)
	}
	println("更新后的用户信息：")
	println("昵称:", getRsp.Nickname)
	println("邮箱:", getRsp.Email)
	birthday := time.Unix(int64(getRsp.Birthday), 0)
	println("生日:", birthday.Format("2006-01-02"))
}

func main() {
	Init()
	defer conn.Close()

	TestGetUserList()
	TestCreateUser()
	TestBatchCreateUsers()
	TestUpdateUser()
	TestDeleteUser()
}
