/*
 * @Author: Will nanan_zhao@163.com
 * @Date: 2025-05-18 18:29:58
 * @LastEditors: Will nanan_zhao@163.com
 * @LastEditTime: 2025-05-18 19:26:53
 * @FilePath: /joyshop_srvs/goods_srv/tests/init_test.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package tests

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"goods_srv/global"
	"goods_srv/initialize"
	"goods_srv/proto"
	"goods_srv/util"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var goodsClient proto.GoodsClient
var conn *grpc.ClientConn

func initTestConfig() error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("获取工作目录失败: %v", err)
	}
	configFile := filepath.Join(dir, "..", "config", "config-develop.yaml")
	zap.S().Infof("正在加载配置文件: %s", configFile)
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return fmt.Errorf("配置文件不存在: %s", configFile)
	}
	v := viper.New()
	v.SetConfigFile(configFile)
	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败: %v", err)
	}
	if err := v.Unmarshal(&global.ServerConfig); err != nil {
		return fmt.Errorf("解析配置文件失败: %v", err)
	}
	return nil
}

func NewClient(addr string) (*grpc.ClientConn, error) {
	return grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
}

func TestMain(m *testing.M) {
	initialize.InitLogger()
	if err := initTestConfig(); err != nil {
		zap.S().Fatalf("初始化测试配置失败: %v", err)
	}
	initialize.InitDB()

	var err error
	conn, err = NewClient(fmt.Sprintf("%s:%d", global.ServerConfig.Host, global.ServerConfig.Port))
	if err != nil {
		zap.S().Errorw("[GetGoodsSrvClient] 连接商品服务失败", "msg", err.Error())
		os.Exit(1)
	}
	goodsClient = proto.NewGoodsClient(conn)
	code := m.Run()
	if err := conn.Close(); err != nil {
		zap.S().Errorf("关闭gRPC连接失败: %v", err)
	}
	os.Exit(code)
}
func TestInit(t *testing.T) {
	// TestMain(nil)
	err := util.DropAllTables()
	if err != nil {
		t.Fatalf("删除所有测试表结构失败: %v", err)
	}
	err = util.SetupTestTables()
	if err != nil {
		t.Fatalf("设置测试表结构失败: %v", err)
	}
}
