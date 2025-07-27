/*
 * @Author: Will nanan_zhao@163.com
 * @Date: 2025-05-09 09:55:55
 * @LastEditors: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @LastEditTime: 2025-07-26 13:44:58
 * @FilePath: /joyshop_srvs/user_srv/tests/user_test.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package test

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"
	"user_srv/global"
	"user_srv/initialize"
	"user_srv/model"
	"user_srv/proto"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	userClient proto.UserClient
	conn       *grpc.ClientConn
)

// 测试专用的配置初始化
func initTestConfig() error {
	// 获取当前目录
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("获取工作目录失败: %v", err)
	}

	// 设置配置文件路径为父级目录下的config目录
	configFile := filepath.Join(dir, "..", "config", "config-develop.yaml")
	zap.S().Infof("正在加载配置文件: %s", configFile)

	// 检查配置文件是否存在
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return fmt.Errorf("配置文件不存在: %s", configFile)
	}

	// 初始化viper
	v := viper.New()
	v.SetConfigFile(configFile)

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 解析配置到结构体
	if err := v.Unmarshal(&global.ServerConfig); err != nil {
		return fmt.Errorf("解析配置文件失败: %v", err)
	}

	return nil
}

func TestMain(m *testing.M) {
	// 初始化日志
	initialize.InitLogger()

	// 使用测试专用的配置初始化
	if err := initTestConfig(); err != nil {
		zap.S().Fatalf("初始化测试配置失败: %v", err)
	}

	// 初始化数据库
	if err := initialize.InitDB(); err != nil {
		zap.S().Fatalf("初始化数据库失败: %v", err)
	}

	// 检查表是否存在
	var count int64
	if err := global.DB.Table("user").Count(&count).Error; err != nil {
		// 如果表不存在，则创建表
		if err := global.DB.AutoMigrate(&model.User{}); err != nil {
			zap.S().Fatalf("创建表结构失败: %v", err)
		}
		zap.S().Info("表结构创建成功")
	} else {
		zap.S().Info("表已存在，跳过创建")
	}

	// 初始化gRPC连接
	var err error
	conn, err = grpc.NewClient(fmt.Sprintf("%s:%d", global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ServerInfo.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.S().Errorw("[GetUserSrvClient] 连接 【用户服务失败】", "msg", err.Error())
		return
	}
	userClient = proto.NewUserClient(conn)

	// 运行测试
	code := m.Run()

	// 清理资源
	if err := conn.Close(); err != nil {
		zap.S().Errorf("关闭gRPC连接失败: %v", err)
	}

	// 退出
	if code != 0 {
		zap.S().Fatalf("测试失败，退出码: %d", code)
	}
}

// 测试获取用户列表
func TestGetUserList(t *testing.T) {
	rsp, err := userClient.GetUserList(context.Background(), &proto.PageInfo{
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		t.Errorf("获取用户列表失败: %v", err)
		return
	}
	for _, user := range rsp.Data {
		t.Logf("用户信息 - 手机号: %s, 昵称: %s", user.Mobile, user.Nickname)

		// 检查密码
		checkRsp, err := userClient.CheckPassword(context.Background(), &proto.PasswordCheckInof{
			Password:        "test123",
			EncryptPassword: user.Password,
		})
		if err != nil {
			t.Errorf("检查密码失败: %v", err)
			continue
		}
		t.Logf("密码验证结果: %v", checkRsp.Success)
	}
}

// 测试创建单个用户
func TestCreateUser(t *testing.T) {
	for i := 0; i < 10; i++ {
		rsp, err := userClient.CreateUser(context.Background(), &proto.CreateUserInfo{
			Nickname: fmt.Sprintf("testuser%d", i),
			Mobile:   generateRandomMobile(),
			Password: "admin123",
			Email:    generateRandomEmail(),
			Username: generateRandomUsername(),
		})
		if err != nil {
			t.Errorf("创建用户失败: %v", err)
			continue
		}
		t.Logf("创建用户成功，ID: %d", rsp.Id)
	}
}

// 测试根据手机号获取用户
func TestGetUserByMobile(t *testing.T) {
	rsp, err := userClient.GetUserByMobile(context.Background(), &proto.MobileRequest{
		Mobile: "18782222220",
	})
	if err != nil {
		t.Errorf("根据手机号获取用户失败: %v", err)
		return
	}
	t.Logf("获取到的用户信息: %+v", rsp)
}

// 测试根据ID获取用户
func TestGetUserById(t *testing.T) {
	// 先创建一个用户，然后根据ID获取
	createRsp, err := userClient.CreateUser(context.Background(), &proto.CreateUserInfo{
		Password: "test123",
		Mobile:   generateRandomMobile(),
		Email:    generateRandomEmail(),
		Nickname: generateRandomNickname(),
		Username: generateRandomUsername(),
		Birthday: uint64(time.Now().Unix()),
	})
	if err != nil {
		t.Errorf("创建用户失败: %v", err)
		return
	}

	// 根据ID获取用户
	rsp, err := userClient.GetUserById(context.Background(), &proto.IdRequest{
		Id: createRsp.Id,
	})
	if err != nil {
		t.Errorf("根据ID获取用户失败: %v", err)
		return
	}
	t.Logf("获取到的用户信息: %+v", rsp)
}

// 测试更新用户信息
func TestUpdateUser(t *testing.T) {
	// 先创建一个用户
	createRsp, err := userClient.CreateUser(context.Background(), &proto.CreateUserInfo{
		Password: "test123",
		Mobile:   generateRandomMobile(),
		Email:    generateRandomEmail(),
		Nickname: generateRandomNickname(),
		Username: generateRandomUsername(),
		Birthday: uint64(time.Now().Unix()),
	})
	if err != nil {
		t.Errorf("创建用户失败: %v", err)
		return
	}
	t.Logf("创建用户成功，ID: %d", createRsp.Id)

	// 更新用户信息
	_, err = userClient.UpdateUser(context.Background(), &proto.UpdateUserInfo{
		Id:       createRsp.Id,
		Nickname: generateRandomNickname(),
		Email:    generateRandomEmail(),
		Birthday: uint64(time.Now().AddDate(0, 0, 1).Unix()),
	})
	if err != nil {
		t.Errorf("更新用户失败: %v", err)
		return
	}
	t.Log("更新用户成功")

	// 验证更新结果
	getRsp, err := userClient.GetUserById(context.Background(), &proto.IdRequest{Id: createRsp.Id})
	if err != nil {
		t.Errorf("获取更新后的用户信息失败: %v", err)
		return
	}
	t.Logf("更新后的用户信息: 昵称=%s, 邮箱=%s", getRsp.Nickname, getRsp.Email)
}

// 测试删除用户
func TestDeleteUser(t *testing.T) {
	// 先创建一个用户
	createRsp, err := userClient.CreateUser(context.Background(), &proto.CreateUserInfo{
		Password: "test123",
		Mobile:   generateRandomMobile(),
		Email:    generateRandomEmail(),
		Nickname: generateRandomNickname(),
		Username: generateRandomUsername(),
		Birthday: uint64(time.Now().Unix()),
	})
	if err != nil {
		t.Errorf("创建用户失败: %v", err)
		return
	}
	t.Logf("创建用户成功，ID: %d", createRsp.Id)

	// 删除用户
	_, err = userClient.DeleteUser(context.Background(), &proto.IdRequest{Id: createRsp.Id})
	if err != nil {
		t.Errorf("删除用户失败: %v", err)
		return
	}
	t.Logf("删除用户成功，ID: %d", createRsp.Id)

	// 验证用户是否已删除
	_, err = userClient.GetUserById(context.Background(), &proto.IdRequest{Id: createRsp.Id})
	if err == nil {
		t.Error("用户应该已被删除")
		return
	}
	t.Log("验证用户已删除成功")
}

// 生成随机手机号
func generateRandomMobile() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	mobile := "1"
	for i := 0; i < 10; i++ {
		mobile += strconv.Itoa(r.Intn(10))
	}
	return mobile
}

// 生成随机邮箱
func generateRandomEmail() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
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
