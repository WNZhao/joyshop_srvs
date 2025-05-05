package main

import (
	"crypto/sha512"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"joyshop_srvs/user_srv/model"
	"log"
	"os"
	"path/filepath"
	"time"
)

// 同步表结构到数据库
type Config struct {
	Database struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Name     string `yaml:"name"`
	} `yaml:"database"`
}

func loadConfig(path string) Config {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("failed to read config file:", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatal("failed to parse config file:", err)
	}

	return cfg
}

// 自动获取配置文件绝对路径（从项目根开始）
func getConfigPath(relative string) (string, error) {
	// 获取当前运行时的“工作目录”（通常是根目录）
	rootPath, err := os.Getwd()
	fmt.Println(rootPath)
	if err != nil {
		return "", fmt.Errorf("get working directory failed: %v", err)
	}

	// 拼接相对路径
	configPath := filepath.Join(rootPath, relative)
	return configPath, nil
}

func main() {
	// 自动拼出配置路径：从项目根读取 config/db/db.yaml
	configPath, err := getConfigPath("user_srv/config/db/db.yaml")
	if err != nil {
		log.Fatal("配置路径解析失败:", err)
	}
	// 读取配置
	cfg := loadConfig(configPath)
	// 拼接DSN
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)
	// 设置打印日志
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second * 10,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)
	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
		Logger: newLogger,
	})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}
	// 设置全局的logger,这个logger在我们执行每个sql语句时都会打印出来

	fmt.Println("数据库连接成功！", db)
	// 定义表结构，将表结构直接生成对应的表 自动建表
	//db.AutoMigrate(&model.User{})
	// 创建用户
	options := &password.Options{10, 100, 32, sha512.New}
	salt, encodedPwd := password.Encode("admin123", options)
	//把盐值和加密后的密码存储到数据库中 盐值整合到加密后的密码中
	newPassword := fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)
	//fmt.Println(newPassword)
	//fmt.Println(salt) // 8c7a9f3b
	//passwordInfo := strings.Split(newPassword, "$")

	for i := 0; i < 10; i++ {
		user := &model.User{
			Mobile:   fmt.Sprintf("1380013800%d", i),
			NickName: fmt.Sprintf("小明%d", i),
			UserName: fmt.Sprintf("xiaoming%d", i),
			Password: newPassword,
			Email:    fmt.Sprintf("xiaoming%d@123.com", i),
		}
		db.Save(&user)
	}
}
