/*
 * @Author: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @Date: 2025-05-08 17:58:32
 * @LastEditors: Will nanan_zhao@163.com
 * @LastEditTime: 2025-05-12 16:58:25
 * @FilePath: /joyshop_srvs/user_srv/config/config.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package config

type ServerConfig struct {
	ConsulInfo ConsulConfig `mapstructure:"consul"`
	DBConfig   DBConfig     `mapstructure:"db"`
	ServerInfo ServerInfo   `mapstructure:"server"`
}

type DBConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	User     string `mapstructure:"user" json:"user"`
	Password string `mapstructure:"password" json:"password"`
	DBName   string `mapstructure:"dbname" json:"dbname"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type ServerInfo struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
	Name string `mapstructure:"name" json:"name"`
}

type NacosConfig struct {
	Host      string `mapstructure:"host" yaml:"host"`
	Port      uint64 `mapstructure:"port" yaml:"port"`
	Namespace string `mapstructure:"namespace" yaml:"namespace"`
	Timeout   uint64 `mapstructure:"timeout" yaml:"timeout"`
	LogDir    string `mapstructure:"logDir" yaml:"logDir"`
	CacheDir  string `mapstructure:"cacheDir" yaml:"cacheDir"`
	LogLevel  string `mapstructure:"logLevel" yaml:"logLevel"`
	DataId    string `mapstructure:"dataId" yaml:"dataId"`
	Group     string `mapstructure:"group" yaml:"group"`
}
