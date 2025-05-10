package util

import (
	"net"
	"os"
)

// GetFreePort 获取一个可用的端口
func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

// GetServerPort 根据环境获取服务端口
// 如果是测试环境（通过环境变量 GO_ENV=test 判断），则使用配置文件中的端口
// 如果是生产环境，则自动获取可用端口
func GetServerPort(configPort int) (int, error) {
	// 检查环境变量
	if os.Getenv("GO_ENV") == "test" {
		// 测试环境：使用配置文件中的端口
		return configPort, nil
	}
	// 生产环境：自动获取可用端口
	return GetFreePort()
}
