package model

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"userop_srv/config"
	"userop_srv/global"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// setupTestDB 初始化测试数据库
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
		cfg.MySQL.User,
		cfg.MySQL.Password,
		cfg.MySQL.Host,
		cfg.MySQL.Port,
		cfg.MySQL.DBName,
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

	// 自动迁移用户操作相关表结构
	if err := db.AutoMigrate(&UserFav{}, &Address{}, &LeavingMessage{}); err != nil {
		t.Fatalf("自动迁移表结构失败: %v", err)
	}
}

func cleanTestData() {
	global.DB.Exec("DELETE FROM userfav")
	global.DB.Exec("DELETE FROM address")
	global.DB.Exec("DELETE FROM leavingmessages")
}

// TestUserFavCRUD 测试用户收藏的CRUD操作
func TestUserFavCRUD(t *testing.T) {
	setupTestDB(t)
	defer cleanTestData()

	// 创建测试
	t.Run("Create", func(t *testing.T) {
		userFav := &UserFav{
			User:  1,
			Goods: 100,
		}
		err := global.DB.Create(userFav).Error
		assert.NoError(t, err)
		assert.NotZero(t, userFav.ID)
	})

	// 查询测试
	t.Run("Read", func(t *testing.T) {
		// 先创建数据
		userFav := &UserFav{
			User:  2,
			Goods: 200,
		}
		global.DB.Create(userFav)

		// 查询单条
		var found UserFav
		err := global.DB.Where("user = ? AND goods = ?", 2, 200).First(&found).Error
		assert.NoError(t, err)
		assert.Equal(t, int32(2), found.User)
		assert.Equal(t, int32(200), found.Goods)

		// 查询列表
		var favList []UserFav
		err = global.DB.Where("user = ?", 2).Find(&favList).Error
		assert.NoError(t, err)
		assert.Len(t, favList, 1)
	})

	// 更新测试
	t.Run("Update", func(t *testing.T) {
		// 先创建数据
		userFav := &UserFav{
			User:  3,
			Goods: 300,
		}
		global.DB.Create(userFav)

		// 更新商品ID
		err := global.DB.Model(&UserFav{}).Where("id = ?", userFav.ID).Update("goods", 301).Error
		assert.NoError(t, err)

		// 验证更新
		var updated UserFav
		global.DB.First(&updated, userFav.ID)
		assert.Equal(t, int32(301), updated.Goods)
	})

	// 删除测试
	t.Run("Delete", func(t *testing.T) {
		// 先创建数据
		userFav := &UserFav{
			User:  4,
			Goods: 400,
		}
		global.DB.Create(userFav)

		// 删除
		err := global.DB.Where("user = ? AND goods = ?", 4, 400).Delete(&UserFav{}).Error
		assert.NoError(t, err)

		// 验证删除
		var count int64
		global.DB.Model(&UserFav{}).Where("user = ? AND goods = ?", 4, 400).Count(&count)
		assert.Equal(t, int64(0), count)
	})

	// 唯一索引测试
	t.Run("UniqueIndex", func(t *testing.T) {
		// 创建第一条记录
		userFav1 := &UserFav{
			User:  5,
			Goods: 500,
		}
		err := global.DB.Create(userFav1).Error
		assert.NoError(t, err)

		// 尝试创建重复记录（相同的user和goods组合）
		userFav2 := &UserFav{
			User:  5,
			Goods: 500,
		}
		err = global.DB.Create(userFav2).Error
		assert.Error(t, err) // 应该报错，因为违反唯一索引
	})
}

// TestAddressCRUD 测试收货地址的CRUD操作
func TestAddressCRUD(t *testing.T) {
	setupTestDB(t)
	defer cleanTestData()

	// 创建测试
	t.Run("Create", func(t *testing.T) {
		address := &Address{
			User:         1,
			Province:     "浙江省",
			City:         "杭州市",
			District:     "西湖区",
			Address:      "文三路100号",
			SignerName:   "张三",
			SignerMobile: "13800138000",
		}
		err := global.DB.Create(address).Error
		assert.NoError(t, err)
		assert.NotZero(t, address.ID)
	})

	// 查询测试
	t.Run("Read", func(t *testing.T) {
		// 先创建数据
		address := &Address{
			User:         2,
			Province:     "北京市",
			City:         "北京市",
			District:     "朝阳区",
			Address:      "建国路88号",
			SignerName:   "李四",
			SignerMobile: "13900139000",
		}
		global.DB.Create(address)

		// 查询单条
		var found Address
		err := global.DB.Where("user = ?", 2).First(&found).Error
		assert.NoError(t, err)
		assert.Equal(t, "李四", found.SignerName)
		assert.Equal(t, "北京市", found.Province)

		// 查询用户的所有地址
		var addressList []Address
		err = global.DB.Where("user = ?", 2).Find(&addressList).Error
		assert.NoError(t, err)
		assert.Len(t, addressList, 1)
	})

	// 更新测试
	t.Run("Update", func(t *testing.T) {
		// 先创建数据
		address := &Address{
			User:         3,
			Province:     "广东省",
			City:         "深圳市",
			District:     "南山区",
			Address:      "科技园路1号",
			SignerName:   "王五",
			SignerMobile: "13700137000",
		}
		global.DB.Create(address)

		// 更新地址信息
		updates := map[string]interface{}{
			"address":       "科技园路2号",
			"signer_mobile": "13700137001",
		}
		err := global.DB.Model(&Address{}).Where("id = ?", address.ID).Updates(updates).Error
		assert.NoError(t, err)

		// 验证更新
		var updated Address
		global.DB.First(&updated, address.ID)
		assert.Equal(t, "科技园路2号", updated.Address)
		assert.Equal(t, "13700137001", updated.SignerMobile)
	})

	// 删除测试
	t.Run("Delete", func(t *testing.T) {
		// 先创建数据
		address := &Address{
			User:         4,
			Province:     "上海市",
			City:         "上海市",
			District:     "浦东新区",
			Address:      "张江高科技园区",
			SignerName:   "赵六",
			SignerMobile: "13600136000",
		}
		global.DB.Create(address)

		// 删除
		err := global.DB.Delete(&Address{}, address.ID).Error
		assert.NoError(t, err)

		// 验证删除（软删除）
		var count int64
		global.DB.Model(&Address{}).Where("id = ?", address.ID).Count(&count)
		assert.Equal(t, int64(0), count)

		// 包含软删除的查询
		var countWithDeleted int64
		global.DB.Unscoped().Model(&Address{}).Where("id = ?", address.ID).Count(&countWithDeleted)
		assert.Equal(t, int64(1), countWithDeleted)
	})

	// 批量查询测试
	t.Run("BatchQuery", func(t *testing.T) {
		// 创建多个地址
		userId := int32(10)
		for i := 0; i < 3; i++ {
			address := &Address{
				User:         userId,
				Province:     "测试省",
				City:         "测试市",
				District:     "测试区",
				Address:      "测试路",
				SignerName:   "测试人",
				SignerMobile: "13500135000",
			}
			global.DB.Create(address)
		}

		// 查询该用户的所有地址
		var addresses []Address
		err := global.DB.Where("user = ?", userId).Find(&addresses).Error
		assert.NoError(t, err)
		assert.Len(t, addresses, 3)
	})
}

// TestLeavingMessageCRUD 测试留言的CRUD操作
func TestLeavingMessageCRUD(t *testing.T) {
	setupTestDB(t)
	defer cleanTestData()

	// 创建测试
	t.Run("Create", func(t *testing.T) {
		message := &LeavingMessage{
			User:        1,
			MessageType: 1, // 留言
			Subject:     "产品咨询",
			Message:     "请问这个产品有货吗？",
			File:        "http://example.com/file.jpg",
		}
		err := global.DB.Create(message).Error
		assert.NoError(t, err)
		assert.NotZero(t, message.ID)
	})

	// 查询测试
	t.Run("Read", func(t *testing.T) {
		// 先创建数据
		message := &LeavingMessage{
			User:        2,
			MessageType: 2, // 投诉
			Subject:     "服务投诉",
			Message:     "对服务不满意",
			File:        "",
		}
		global.DB.Create(message)

		// 查询单条
		var found LeavingMessage
		err := global.DB.Where("user = ? AND message_type = ?", 2, 2).First(&found).Error
		assert.NoError(t, err)
		assert.Equal(t, "服务投诉", found.Subject)
		assert.Equal(t, int32(2), found.MessageType)

		// 查询用户的所有留言
		var messageList []LeavingMessage
		err = global.DB.Where("user = ?", 2).Find(&messageList).Error
		assert.NoError(t, err)
		assert.Len(t, messageList, 1)
	})

	// 更新测试
	t.Run("Update", func(t *testing.T) {
		// 先创建数据
		message := &LeavingMessage{
			User:        3,
			MessageType: 3, // 询问
			Subject:     "价格询问",
			Message:     "这个产品的价格是多少？",
			File:        "",
		}
		global.DB.Create(message)

		// 更新留言内容
		err := global.DB.Model(&LeavingMessage{}).Where("id = ?", message.ID).
			Update("message", "这个产品的批发价是多少？").Error
		assert.NoError(t, err)

		// 验证更新
		var updated LeavingMessage
		global.DB.First(&updated, message.ID)
		assert.Equal(t, "这个产品的批发价是多少？", updated.Message)
	})

	// 删除测试
	t.Run("Delete", func(t *testing.T) {
		// 先创建数据
		message := &LeavingMessage{
			User:        4,
			MessageType: 4, // 售后
			Subject:     "退货申请",
			Message:     "产品有质量问题，申请退货",
			File:        "http://example.com/evidence.jpg",
		}
		global.DB.Create(message)

		// 删除
		err := global.DB.Delete(&LeavingMessage{}, message.ID).Error
		assert.NoError(t, err)

		// 验证删除
		var count int64
		global.DB.Model(&LeavingMessage{}).Where("id = ?", message.ID).Count(&count)
		assert.Equal(t, int64(0), count)
	})

	// 分页查询测试
	t.Run("Pagination", func(t *testing.T) {
		// 创建多条留言
		userId := int32(20)
		for i := 1; i <= 5; i++ {
			message := &LeavingMessage{
				User:        userId,
				MessageType: int32(i),
				Subject:     "测试主题",
				Message:     "测试内容",
				File:        "",
			}
			global.DB.Create(message)
		}

		// 分页查询
		var messages []LeavingMessage
		err := global.DB.Where("user = ?", userId).Limit(3).Offset(0).Find(&messages).Error
		assert.NoError(t, err)
		assert.Len(t, messages, 3)

		// 查询第二页
		var page2 []LeavingMessage
		err = global.DB.Where("user = ?", userId).Limit(3).Offset(3).Find(&page2).Error
		assert.NoError(t, err)
		assert.Len(t, page2, 2)
	})

	// 按类型查询测试
	t.Run("QueryByType", func(t *testing.T) {
		// 创建不同类型的留言
		types := []int32{1, 2, 3, 4, 5} // 留言、投诉、询问、售后、求购
		for _, msgType := range types {
			message := &LeavingMessage{
				User:        30,
				MessageType: msgType,
				Subject:     "测试主题",
				Message:     "测试内容",
				File:        "",
			}
			global.DB.Create(message)
		}

		// 查询特定类型
		var complaints []LeavingMessage
		err := global.DB.Where("user = ? AND message_type = ?", 30, 2).Find(&complaints).Error
		assert.NoError(t, err)
		assert.Len(t, complaints, 1)
		assert.Equal(t, int32(2), complaints[0].MessageType)
	})
}

// TestTableNames 测试表名设置
func TestTableNames(t *testing.T) {
	assert.Equal(t, "userfav", UserFav{}.TableName())
	assert.Equal(t, "address", Address{}.TableName())
	assert.Equal(t, "leavingmessages", LeavingMessage{}.TableName())
}

// TestBaseModel 测试基础模型字段
func TestBaseModel(t *testing.T) {
	setupTestDB(t)
	defer cleanTestData()

	// 测试时间字段自动填充
	t.Run("Timestamps", func(t *testing.T) {
		userFav := &UserFav{
			User:  100,
			Goods: 200,
		}
		
		// 创建时自动填充创建时间和更新时间
		err := global.DB.Create(userFav).Error
		assert.NoError(t, err)
		assert.NotZero(t, userFav.CreatedAt)
		assert.NotZero(t, userFav.UpdatedAt)
		
		// 记录创建时间
		createdAt := userFav.CreatedAt
		
		// 等待一小段时间
		time.Sleep(100 * time.Millisecond)
		
		// 更新记录
		err = global.DB.Model(userFav).Update("goods", 201).Error
		assert.NoError(t, err)
		
		// 重新查询
		var updated UserFav
		global.DB.First(&updated, userFav.ID)
		
		// 创建时间不变，更新时间改变
		assert.Equal(t, createdAt.Unix(), updated.CreatedAt.Unix())
		assert.True(t, updated.UpdatedAt.After(createdAt))
	})

	// 测试软删除
	t.Run("SoftDelete", func(t *testing.T) {
		address := &Address{
			User:         200,
			Province:     "测试省",
			City:         "测试市",
			District:     "测试区",
			Address:      "测试地址",
			SignerName:   "测试",
			SignerMobile: "13800138000",
		}
		
		// 创建记录
		err := global.DB.Create(address).Error
		assert.NoError(t, err)
		addressID := address.ID
		
		// 软删除
		err = global.DB.Delete(address).Error
		assert.NoError(t, err)
		
		// 普通查询找不到（因为软删除）
		var notFound Address
		err = global.DB.First(&notFound, addressID).Error
		assert.Error(t, err)
		
		// 使用Unscoped可以找到
		var found Address
		err = global.DB.Unscoped().First(&found, addressID).Error
		assert.NoError(t, err)
		assert.NotNil(t, found.DeletedAt)
	})
}