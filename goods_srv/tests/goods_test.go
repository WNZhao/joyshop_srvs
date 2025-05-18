package tests

import (
	"context"
	"fmt"
	"goods_srv/global"
	"goods_srv/initialize"
	"goods_srv/proto"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	goodsClient proto.GoodsClient
	conn        *grpc.ClientConn
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
	initialize.InitDB()

	// 初始化gRPC连接
	var err error
	conn, err = grpc.NewClient(
		fmt.Sprintf("%s:%d", global.ServerConfig.Host, global.ServerConfig.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		zap.S().Errorw("[GetGoodsSrvClient] 连接商品服务失败", "msg", err.Error())
		return
	}
	goodsClient = proto.NewGoodsClient(conn)

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

// 测试获取商品列表
func TestGoodsList(t *testing.T) {
	// 测试基本分页
	t.Run("基本分页", func(t *testing.T) {
		rsp, err := goodsClient.GoodsList(context.Background(), &proto.GoodsFilterRequest{
			Pages:       1,
			PagePerNums: 10,
		})
		if err != nil {
			t.Errorf("获取商品列表失败: %v", err)
			return
		}
		t.Logf("商品总数: %d", rsp.Total)
		for _, goods := range rsp.Data {
			t.Logf("商品信息 - ID: %d, 名称: %s, 价格: %.2f", goods.Id, goods.Name, goods.ShopPrice)
		}
	})

	// 测试按品牌筛选
	t.Run("按品牌筛选", func(t *testing.T) {
		rsp, err := goodsClient.GoodsList(context.Background(), &proto.GoodsFilterRequest{
			Pages:       1,
			PagePerNums: 10,
			BrandId:     1, // 假设品牌ID为1
		})
		if err != nil {
			t.Errorf("按品牌获取商品列表失败: %v", err)
			return
		}
		t.Logf("品牌商品总数: %d", rsp.Total)
		for _, goods := range rsp.Data {
			t.Logf("品牌商品信息 - ID: %d, 名称: %s, 品牌ID: %d", goods.Id, goods.Name, goods.BrandId)
		}
	})

	// 测试按分类筛选
	t.Run("按分类筛选", func(t *testing.T) {
		rsp, err := goodsClient.GoodsList(context.Background(), &proto.GoodsFilterRequest{
			Pages:       1,
			PagePerNums: 10,
			CategoryId:  1, // 假设分类ID为1
		})
		if err != nil {
			t.Errorf("按分类获取商品列表失败: %v", err)
			return
		}
		t.Logf("分类商品总数: %d", rsp.Total)
		for _, goods := range rsp.Data {
			t.Logf("分类商品信息 - ID: %d, 名称: %s, 分类IDs: %v", goods.Id, goods.Name, goods.CategoryIds)
		}
	})

	// 测试按关键词搜索
	t.Run("按关键词搜索", func(t *testing.T) {
		rsp, err := goodsClient.GoodsList(context.Background(), &proto.GoodsFilterRequest{
			Pages:       1,
			PagePerNums: 10,
			Keywords:    "手机", // 假设搜索关键词为"手机"
		})
		if err != nil {
			t.Errorf("按关键词获取商品列表失败: %v", err)
			return
		}
		t.Logf("搜索结果总数: %d", rsp.Total)
		for _, goods := range rsp.Data {
			t.Logf("搜索结果 - ID: %d, 名称: %s", goods.Id, goods.Name)
		}
	})

	// 测试热门商品
	t.Run("热门商品", func(t *testing.T) {
		rsp, err := goodsClient.GoodsList(context.Background(), &proto.GoodsFilterRequest{
			Pages:       1,
			PagePerNums: 10,
			IsHot:       true,
		})
		if err != nil {
			t.Errorf("获取热门商品列表失败: %v", err)
			return
		}
		t.Logf("热门商品总数: %d", rsp.Total)
		for _, goods := range rsp.Data {
			t.Logf("热门商品信息 - ID: %d, 名称: %s, 是否热门: %v", goods.Id, goods.Name, goods.IsHot)
		}
	})

	// 测试新品
	t.Run("新品", func(t *testing.T) {
		rsp, err := goodsClient.GoodsList(context.Background(), &proto.GoodsFilterRequest{
			Pages:       1,
			PagePerNums: 10,
			IsNew:       true,
		})
		if err != nil {
			t.Errorf("获取新品列表失败: %v", err)
			return
		}
		t.Logf("新品总数: %d", rsp.Total)
		for _, goods := range rsp.Data {
			t.Logf("新品信息 - ID: %d, 名称: %s, 是否新品: %v", goods.Id, goods.Name, goods.IsNew)
		}
	})

	// 测试在售商品
	t.Run("在售商品", func(t *testing.T) {
		rsp, err := goodsClient.GoodsList(context.Background(), &proto.GoodsFilterRequest{
			Pages:       1,
			PagePerNums: 10,
			OnSale:      true,
		})
		if err != nil {
			t.Errorf("获取在售商品列表失败: %v", err)
			return
		}
		t.Logf("在售商品总数: %d", rsp.Total)
		for _, goods := range rsp.Data {
			t.Logf("在售商品信息 - ID: %d, 名称: %s, 是否在售: %v", goods.Id, goods.Name, goods.OnSale)
		}
	})
}
