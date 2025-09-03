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
	t.Log(TestScenarios.UserLogin)
	
	rsp, err := userClient.GetUserList(context.Background(), &proto.PageInfo{
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		t.Errorf("获取用户列表失败: %v", err)
		return
	}
	
	// 验证返回的用户数量应该大于0
	if len(rsp.Data) == 0 {
		t.Error("用户列表为空")
		return
	}
	
	t.Logf("获取到用户列表，总数: %d", rsp.Total)
	for _, user := range rsp.Data {
		t.Logf("用户信息 - ID: %d, 用户名: %s, 手机号: %s, 昵称: %s, 角色: %d", 
			user.Id, user.Username, user.Mobile, user.Nickname, user.Role)

		// 使用真实密码检查（统一为123456）
		checkRsp, err := userClient.CheckPassword(context.Background(), &proto.PasswordCheckInof{
			Password:        "123456",  // 使用真实密码
			EncryptPassword: user.Password,
		})
		if err != nil {
			t.Errorf("检查用户ID %d密码失败: %v", user.Id, err)
			continue
		}
		t.Logf("用户ID %d 密码验证结果: %v", user.Id, checkRsp.Success)
	}
}

// 测试创建单个用户
func TestCreateUser(t *testing.T) {
	t.Log(TestScenarios.UserRegistration)
	
	// 使用预定义的测试用户数据创建用户
	testUser := NewTestUser
	rsp, err := userClient.CreateUser(context.Background(), &proto.CreateUserInfo{
		Nickname: testUser.Nickname,
		Mobile:   testUser.Mobile,
		Password: testUser.Password,
		Email:    testUser.Email,
		Username: testUser.Username,
	})
	if err != nil {
		t.Errorf("创建用户失败: %v", err)
		return
	}
	t.Logf("创建用户成功 - ID: %d, 用户名: %s, 昵称: %s", rsp.Id, testUser.Username, testUser.Nickname)
	
	// 验证创建的用户信息
	getRsp, err := userClient.GetUserById(context.Background(), &proto.IdRequest{Id: rsp.Id})
	if err != nil {
		t.Errorf("获取新创建的用户信息失败: %v", err)
		return
	}
	
	if getRsp.Username != testUser.Username {
		t.Errorf("用户名不匹配: 期望 %s, 实际 %s", testUser.Username, getRsp.Username)
	}
	if getRsp.Mobile != testUser.Mobile {
		t.Errorf("手机号不匹配: 期望 %s, 实际 %s", testUser.Mobile, getRsp.Mobile)
	}
}

// 测试根据手机号获取用户
func TestGetUserByMobile(t *testing.T) {
	t.Log("测试根据手机号获取用户")
	
	// 使用预定义的测试用户手机号
	testUser := RegularUsers[0] // 张三
	rsp, err := userClient.GetUserByMobile(context.Background(), &proto.MobileRequest{
		Mobile: testUser.Mobile,
	})
	if err != nil {
		t.Errorf("根据手机号获取用户失败: %v", err)
		return
	}
	
	// 验证返回的用户信息
	if rsp.Mobile != testUser.Mobile {
		t.Errorf("手机号不匹配: 期望 %s, 实际 %s", testUser.Mobile, rsp.Mobile)
	}
	if rsp.Username != testUser.Username {
		t.Errorf("用户名不匹配: 期望 %s, 实际 %s", testUser.Username, rsp.Username)
	}
	
	t.Logf("获取到的用户信息 - ID: %d, 用户名: %s, 昵称: %s, 手机号: %s", 
		rsp.Id, rsp.Username, rsp.Nickname, rsp.Mobile)
}

// 测试根据ID获取用户
func TestGetUserById(t *testing.T) {
	t.Log("测试根据ID获取用户")
	
	// 使用预定义的测试用户ID
	testUser := RegularUsers[1] // 李四
	rsp, err := userClient.GetUserById(context.Background(), &proto.IdRequest{
		Id: testUser.ID,
	})
	if err != nil {
		t.Errorf("根据ID获取用户失败: %v", err)
		return
	}
	
	// 验证返回的用户信息
	if rsp.Id != testUser.ID {
		t.Errorf("用户ID不匹配: 期望 %d, 实际 %d", testUser.ID, rsp.Id)
	}
	if rsp.Username != testUser.Username {
		t.Errorf("用户名不匹配: 期望 %s, 实际 %s", testUser.Username, rsp.Username)
	}
	if rsp.Mobile != testUser.Mobile {
		t.Errorf("手机号不匹配: 期望 %s, 实际 %s", testUser.Mobile, rsp.Mobile)
	}
	
	t.Logf("获取到的用户信息 - ID: %d, 用户名: %s, 昵称: %s, 手机号: %s, 角色: %d", 
		rsp.Id, rsp.Username, rsp.Nickname, rsp.Mobile, rsp.Role)
}

// 测试更新用户信息
func TestUpdateUser(t *testing.T) {
	t.Log(TestScenarios.UserUpdate)
	
	// 使用预定义的测试用户
	testUser := RegularUsers[2] // 王五
	updatedNickname := "更新后的昵称"
	updatedEmail := "updated@example.com"
	
	// 更新用户信息
	_, err := userClient.UpdateUser(context.Background(), &proto.UpdateUserInfo{
		Id:       testUser.ID,
		Nickname: updatedNickname,
		Email:    updatedEmail,
		Birthday: uint64(time.Now().Unix()),
	})
	if err != nil {
		t.Errorf("更新用户失败: %v", err)
		return
	}
	t.Logf("更新用户成功 - ID: %d", testUser.ID)

	// 验证更新结果
	getRsp, err := userClient.GetUserById(context.Background(), &proto.IdRequest{Id: testUser.ID})
	if err != nil {
		t.Errorf("获取更新后的用户信息失败: %v", err)
		return
	}
	
	// 验证更新的字段
	if getRsp.Nickname != updatedNickname {
		t.Errorf("昵称更新失败: 期望 %s, 实际 %s", updatedNickname, getRsp.Nickname)
	}
	if getRsp.Email != updatedEmail {
		t.Errorf("邮箱更新失败: 期望 %s, 实际 %s", updatedEmail, getRsp.Email)
	}
	
	t.Logf("更新后的用户信息: ID=%d, 昵称=%s, 邮箱=%s", getRsp.Id, getRsp.Nickname, getRsp.Email)
}

// 测试删除用户
func TestDeleteUser(t *testing.T) {
	t.Log(TestScenarios.UserDelete)
	
	// 先创建一个测试用户用于删除
	testUser := TestUserConstants{
		Username: "delete_test_user",
		Password: "123456",
		Mobile:   "19900000000",
		Email:    "delete_test@example.com",
		Nickname: "待删除用户",
		Role:     1,
	}
	
	createRsp, err := userClient.CreateUser(context.Background(), &proto.CreateUserInfo{
		Password: testUser.Password,
		Mobile:   testUser.Mobile,
		Email:    testUser.Email,
		Nickname: testUser.Nickname,
		Username: testUser.Username,
		Birthday: uint64(time.Now().Unix()),
	})
	if err != nil {
		t.Errorf("创建测试用户失败: %v", err)
		return
	}
	t.Logf("创建测试用户成功，ID: %d", createRsp.Id)

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

// 测试管理员登录
func TestAdminLogin(t *testing.T) {
	t.Log(TestScenarios.AdminLogin)
	
	// 测试管理员用户登录
	adminUser := AdminUser
	rsp, err := userClient.GetUserByMobile(context.Background(), &proto.MobileRequest{
		Mobile: adminUser.Mobile,
	})
	if err != nil {
		t.Errorf("获取管理员用户失败: %v", err)
		return
	}
	
	// 验证管理员角色
	if rsp.Role != adminUser.Role {
		t.Errorf("管理员角色不正确: 期望 %d, 实际 %d", adminUser.Role, rsp.Role)
	}
	
	// 验证管理员密码
	checkRsp, err := userClient.CheckPassword(context.Background(), &proto.PasswordCheckInof{
		Password:        adminUser.Password,
		EncryptPassword: rsp.Password,
	})
	if err != nil {
		t.Errorf("检查管理员密码失败: %v", err)
		return
	}
	
	if !checkRsp.Success {
		t.Error("管理员密码验证失败")
	}
	
	t.Logf("管理员登录验证成功 - ID: %d, 用户名: %s, 角色: %d", rsp.Id, rsp.Username, rsp.Role)
}

// 测试VIP用户功能
func TestVIPUserFeatures(t *testing.T) {
	t.Log(TestScenarios.VIPUserLogin)
	
	// 测试VIP用户
	for i, vipUser := range VIPUsers {
		t.Run(fmt.Sprintf("VIP用户%d", i+1), func(t *testing.T) {
			rsp, err := userClient.GetUserById(context.Background(), &proto.IdRequest{
				Id: vipUser.ID,
			})
			if err != nil {
				t.Errorf("获取VIP用户失败: %v", err)
				return
			}
			
			// 验证VIP用户信息
			if rsp.Id != vipUser.ID {
				t.Errorf("VIP用户ID不匹配: 期望 %d, 实际 %d", vipUser.ID, rsp.Id)
			}
			if rsp.Username != vipUser.Username {
				t.Errorf("VIP用户名不匹配: 期望 %s, 实际 %s", vipUser.Username, rsp.Username)
			}
			
			// 验证密码
			checkRsp, err := userClient.CheckPassword(context.Background(), &proto.PasswordCheckInof{
				Password:        vipUser.Password,
				EncryptPassword: rsp.Password,
			})
			if err != nil {
				t.Errorf("检查VIP用户密码失败: %v", err)
				return
			}
			
			if !checkRsp.Success {
				t.Errorf("VIP用户密码验证失败")
			}
			
			t.Logf("VIP用户验证成功 - ID: %d, 用户名: %s, 昵称: %s", rsp.Id, rsp.Username, rsp.Nickname)
		})
	}
}

// 测试批量用户信息验证
func TestBatchUserValidation(t *testing.T) {
	t.Log("批量用户信息验证测试")
	
	// 验证所有预定义用户
	allUsers := AllTestUsers()
	for _, testUser := range allUsers {
		t.Run(fmt.Sprintf("用户ID_%d", testUser.ID), func(t *testing.T) {
			rsp, err := userClient.GetUserById(context.Background(), &proto.IdRequest{
				Id: testUser.ID,
			})
			if err != nil {
				t.Errorf("获取用户ID %d 失败: %v", testUser.ID, err)
				return
			}
			
			// 基础信息验证
			if rsp.Username != testUser.Username {
				t.Errorf("用户名不匹配: 期望 %s, 实际 %s", testUser.Username, rsp.Username)
			}
			if rsp.Mobile != testUser.Mobile {
				t.Errorf("手机号不匹配: 期望 %s, 实际 %s", testUser.Mobile, rsp.Mobile)
			}
			if rsp.Role != testUser.Role {
				t.Errorf("用户角色不匹配: 期望 %d, 实际 %d", testUser.Role, rsp.Role)
			}
			
			// 密码验证
			checkRsp, err := userClient.CheckPassword(context.Background(), &proto.PasswordCheckInof{
				Password:        testUser.Password,
				EncryptPassword: rsp.Password,
			})
			if err != nil {
				t.Errorf("检查用户ID %d 密码失败: %v", testUser.ID, err)
				return
			}
			
			if !checkRsp.Success {
				t.Errorf("用户ID %d 密码验证失败", testUser.ID)
			}
		})
	}
	
	t.Logf("批量验证完成，共验证 %d 个用户", len(allUsers))
}
