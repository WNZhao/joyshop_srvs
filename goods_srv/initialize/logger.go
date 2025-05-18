/*
 * @Author: Will nanan_zhao@163.com
 * @Date: 2025-05-12 16:58:08
 * @LastEditors: Will nanan_zhao@163.com
 * @LastEditTime: 2025-05-17 17:50:25
 * @FilePath: /joyshop_srvs/goods_srv/initialize/logger.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package initialize

import (
	"os"

	"go.uber.org/zap"
)

func InitLogger() {
	var logger *zap.Logger
	var err error

	// 根据环境变量设置日志模式
	if os.Getenv("APP_ENV") == "production" {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}

	if err != nil {
		panic(err)
	}

	zap.ReplaceGlobals(logger)
}
