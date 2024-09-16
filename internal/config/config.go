package config

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"time"
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
	MaxAge        time.Duration
	Format        string
	StackTraceKey string
	EncodeLevel   string
	Prefix        string
	LoginConsole  bool
	ShowLine      bool
}

type Config struct {
	Debug    bool
	Port     int
	Language string `yaml:"language"` // 语言
	Limit    float64
	Cors     Cors
	Redis    struct {
		Instances []RedisInstanceConfig
	}
	Zap Zap
}

var (
	AppConfig Config
	port      int
)

func init() {
	// 设置命令行参数
	flag.IntVar(&port, "port", 0, "Port number")
	flag.Parse()
	println(port)

	viper.SetConfigFile("./config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file%v", err)
		//logger.Fatal("Error reading config file", zap.Error(err))
	}

	if err := viper.Unmarshal(&AppConfig); err != nil {
		fmt.Printf("Unable to decode into struct%v", err)
		//logger.Fatal("Unable to decode into struct", zap.Error(err))
	}

	// 如果命令行参数中有指定端口，则更新配置文件中的端口
	if port != 0 {
		AppConfig.Port = port
		fmt.Printf("使用命令行端口:%v", port)
		//logger.Info("使用命令行端口:", zap.Int("port", port))
	}
}
