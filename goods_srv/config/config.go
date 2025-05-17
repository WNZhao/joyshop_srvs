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

// nacosConfig 是 Nacos 配置的结构体
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
