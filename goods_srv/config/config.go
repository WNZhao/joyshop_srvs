/*
 * @Author: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @Date: 2025-07-16 16:35:24
 * @LastEditors: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @LastEditTime: 2025-07-25 13:26:14
 * @FilePath: /joyshop_srvs/goods_srv/config/config.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package config

type ServerConfig struct {
	Name string   `mapstructure:"name"`
	Host string   `mapstructure:"host"`
	Port int      `mapstructure:"port"`
	Tags []string `mapstructure:"tags"`

	Consul struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	} `mapstructure:"consul"`

	MySQL struct {
		Host         string `mapstructure:"host"`
		Port         int    `mapstructure:"port"`
		User         string `mapstructure:"user"`
		Password     string `mapstructure:"password"`
		DBName       string `mapstructure:"dbname"`
		MaxIdleConns int    `mapstructure:"max_idle_conns"`
		MaxOpenConns int    `mapstructure:"max_open_conns"`
		LogMode      string `mapstructure:"log_mode"`
		LogZap       bool   `mapstructure:"log_zap"`
	} `mapstructure:"mysql"`

	Log struct {
		Level      string `mapstructure:"level"`
		Filename   string `mapstructure:"filename"`
		MaxSize    int    `mapstructure:"max_size"`
		MaxAge     int    `mapstructure:"max_age"`
		MaxBackups int    `mapstructure:"max_backups"`
	} `mapstructure:"log"`
}

// NacosConfig 是 Nacos 配置的结构体
type NacosConfig struct {
	Host      string `mapstructure:"host"`
	Port      uint64 `mapstructure:"port"`
	Namespace string `mapstructure:"namespace"`
	Timeout   uint64 `mapstructure:"timeout"`
	LogDir    string `mapstructure:"logDir"`
	CacheDir  string `mapstructure:"cacheDir"`
	LogLevel  string `mapstructure:"logLevel"`
	DataId    string `mapstructure:"dataId"`
	Group     string `mapstructure:"group"`
	User     string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
}
