package model

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"user_srv/config"
	"user_srv/global"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// 测试数据
var (
	testUsers = []struct {
		Mobile   string
		Email    string
		Password string
		NickName string
		UserName string
		Gender   string
		Role     int
	}{
		{
			"13800138001",
			"user1@example.com",
			"$pbkdf2-sha512$test_salt$encoded_password",
			"测试用户1",
			"testuser1",
			"male",
			1,
		},
		{
			"13800138002",
			"user2@example.com",
			"$pbkdf2-sha512$test_salt$encoded_password",
			"测试用户2",
			"testuser2",
			"female",
			1,
		},
		{
			"13800138003",
			"admin@example.com",
			"$pbkdf2-sha512$test_salt$encoded_password",
			"管理员",
			"admin",
			"male",
			2,
		},
	}
)

// 数据库模型测试 用户信息
func setupTestDB(t *testing.T) {
	// 初始化日志
	zap.S().Info("开始初始化测试数据库...")

	// 获取当前工作目录
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("获取工作目录失败: %v", err)
	}

	// 设置配置文件路径
	configFile := filepath.Join(dir, "..", "config", "config-develop.yaml")
	zap.S().Infof("正在加载配置文件: %s", configFile)

	// 检查配置文件是否存在
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		t.Fatalf("配置文件不存在: %s", configFile)
	}

	// 初始化viper
	v := viper.New()
	v.SetConfigFile(configFile)

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		t.Fatalf("读取配置文件失败: %v", err)
	}

	// 解析配置到结构体
	var cfg config.ServerConfig
	if err := v.Unmarshal(&cfg); err != nil {
		t.Fatalf("解析配置文件失败: %v", err)
	}

	// 构建DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBConfig.User,
		cfg.DBConfig.Password,
		cfg.DBConfig.Host,
		cfg.DBConfig.Port,
		cfg.DBConfig.DBName,
	)

	// 配置GORM
	gormConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
		Logger: logger.Default.LogMode(logger.Info),
	}

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		t.Fatalf("连接数据库失败: %v", err)
	}

	// 设置全局DB实例
	global.DB = db

	// 在测试开始前检查表是否存在
	if !db.Migrator().HasTable(&User{}) {
		if err := db.AutoMigrate(
			&User{},
		); err != nil {
			t.Fatalf("自动迁移表结构失败: %v", err)
		}
	}
}

func TestUserCRUD(t *testing.T) {
	// 设置测试数据库
	setupTestDB(t)

	// 创建测试用户
	user := &User{
		Mobile:   testUsers[0].Mobile,
		Email:    testUsers[0].Email,
		Password: testUsers[0].Password,
		NickName: testUsers[0].NickName,
		UserName: testUsers[0].UserName,
		Gender:   testUsers[0].Gender,
		Role:     testUsers[0].Role,
	}

	// 测试创建用户
	t.Run("Create", func(t *testing.T) {
		result := global.DB.Create(user)
		if result.Error != nil {
			t.Fatalf("创建用户失败: %v", result.Error)
		}
		if user.ID == 0 {
			t.Fatal("创建用户后ID为0")
		}
	})

	// 测试查询用户
	t.Run("Read", func(t *testing.T) {
		var found User
		result := global.DB.First(&found, user.ID)
		if result.Error != nil {
			t.Fatalf("查询用户失败: %v", result.Error)
		}
		if found.Mobile != user.Mobile {
			t.Errorf("查询到的手机号不匹配，期望: %s, 实际: %s", user.Mobile, found.Mobile)
		}
		if found.Email != user.Email {
			t.Errorf("查询到的邮箱不匹配，期望: %s, 实际: %s", user.Email, found.Email)
		}
		if found.UserName != user.UserName {
			t.Errorf("查询到的用户名不匹配，期望: %s, 实际: %s", user.UserName, found.UserName)
		}
	})

	// 测试更新用户
	t.Run("Update", func(t *testing.T) {
		newNickName := "更新后的昵称"
		result := global.DB.Model(user).Update("nick_name", newNickName)
		if result.Error != nil {
			t.Fatalf("更新用户失败: %v", result.Error)
		}

		// 验证更新
		var updated User
		global.DB.First(&updated, user.ID)
		if updated.NickName != newNickName {
			t.Errorf("更新后的昵称不匹配，期望: %s, 实际: %s", newNickName, updated.NickName)
		}
	})

	// 测试删除用户
	t.Run("Delete", func(t *testing.T) {
		result := global.DB.Delete(user)
		if result.Error != nil {
			t.Fatalf("删除用户失败: %v", result.Error)
		}

		// 验证删除
		var deleted User
		result = global.DB.First(&deleted, user.ID)
		if result.Error == nil {
			t.Error("删除后仍能查询到记录")
		}
	})
}

func TestUserList(t *testing.T) {
	// 设置测试数据库
	setupTestDB(t)

	// 创建测试用户
	for _, u := range testUsers {
		user := &User{
			Mobile:   u.Mobile,
			Email:    u.Email,
			Password: u.Password,
			NickName: u.NickName,
			UserName: u.UserName,
			Gender:   u.Gender,
			Role:     u.Role,
		}
		if err := global.DB.Create(user).Error; err != nil {
			t.Fatalf("创建测试用户失败: %v", err)
		}
	}

	// 测试获取用户列表
	t.Run("List", func(t *testing.T) {
		var users []User
		result := global.DB.Find(&users)
		if result.Error != nil {
			t.Fatalf("获取用户列表失败: %v", result.Error)
		}
		if len(users) < len(testUsers) {
			t.Fatalf("获取用户列表数量错误，期望至少 %d，实际 %d", len(testUsers), len(users))
		}
	})

	// 测试按条件查询
	t.Run("QueryByCondition", func(t *testing.T) {
		var users []User
		result := global.DB.Where("role = ?", 1).Find(&users)
		if result.Error != nil {
			t.Fatalf("按角色查询用户失败: %v", result.Error)
		}

		// 验证查询结果
		for _, user := range users {
			if user.Role != 1 {
				t.Errorf("查询到的用户角色不正确，期望: 1, 实际: %d", user.Role)
			}
		}
	})
}

func TestUserUniqueConstraints(t *testing.T) {
	// 设置测试数据库
	setupTestDB(t)

	// 创建第一个用户
	user1 := &User{
		Mobile:   "13800138001",
		Email:    "test1@example.com",
		Password: "password",
		NickName: "用户1",
		UserName: "user1",
		Gender:   "male",
		Role:     1,
	}

	if err := global.DB.Create(user1).Error; err != nil {
		t.Fatalf("创建第一个用户失败: %v", err)
	}

	// 测试手机号唯一性约束
	t.Run("MobileUnique", func(t *testing.T) {
		user2 := &User{
			Mobile:   "13800138001", // 重复的手机号
			Email:    "test2@example.com",
			Password: "password",
			NickName: "用户2",
			UserName: "user2",
			Gender:   "female",
			Role:     1,
		}

		if err := global.DB.Create(user2).Error; err == nil {
			t.Error("创建重复手机号用户应该失败")
		}
	})

	// 测试邮箱唯一性约束
	t.Run("EmailUnique", func(t *testing.T) {
		user3 := &User{
			Mobile:   "13800138002",
			Email:    "test1@example.com", // 重复的邮箱
			Password: "password",
			NickName: "用户3",
			UserName: "user3",
			Gender:   "male",
			Role:     1,
		}

		if err := global.DB.Create(user3).Error; err == nil {
			t.Error("创建重复邮箱用户应该失败")
		}
	})

	// 测试用户名唯一性约束
	t.Run("UserNameUnique", func(t *testing.T) {
		user4 := &User{
			Mobile:   "13800138003",
			Email:    "test4@example.com",
			Password: "password",
			NickName: "用户4",
			UserName: "user1", // 重复的用户名
			Gender:   "female",
			Role:     1,
		}

		if err := global.DB.Create(user4).Error; err == nil {
			t.Error("创建重复用户名用户应该失败")
		}
	})
}
