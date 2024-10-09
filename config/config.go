package config

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
)

// Cors 跨域配置
type Cors struct {
	AllowOrigins     []string `yaml:"allowOrigins"`     // 允许跨域origin
	AllowMethods     string   `yaml:"allowMethods"`     // 方法
	AllowHeaders     string   `yaml:"allowHeaders"`     // 请求头
	ExposeHeaders    string   `yaml:"exposeHeaders"`    //
	AllowCredentials string   `yaml:"allowCredentials"` //
	MaxAge           string   `yaml:"maxAge"`           //
}

type RedisInstanceConfig struct {
	Name     string
	Addr     string
	Port     int
	Password string
	DBs      []int
}

type Zap struct {
	Director      string
	Level         string
	MaxAge        int
	MaxSize       int
	MaxBackups    int
	Format        string
	StackTraceKey string
	EncodeLevel   string
	Prefix        string
	LoginConsole  bool
	ShowLine      bool
}

type Config struct {
	Debug string
	Port  int
	Limit float64
	Cors  Cors
	Redis struct {
		Instances []RedisInstanceConfig
	}
	Zap Zap
}

var (
	AppConfig Config
	port      int
)

func InitConfig() error {
	// 设置命令行参数
	flag.IntVar(&port, "port", 0, "指定端口号")
	flag.Parse()

	// 读取配置文件
	viper.SetConfigFile("./config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("配置文件读取错误: %w", err)
	}

	// 配置转结构体
	if err := viper.Unmarshal(&AppConfig); err != nil {
		return fmt.Errorf("配置解码失败: %w", err)
	}

	// 如果命令行参数中有指定端口，则更新配置文件中的端口
	if port != 0 {
		AppConfig.Port = port
		fmt.Println(fmt.Sprintf("使用命令行端口:%v", port))
	}
	return nil
}
