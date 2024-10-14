package config

import (
	"flag"
	"fmt"

	"github.com/spf13/viper"
)

type RedisInstanceConfig struct {
	Name     string     `yaml:"name"`     // 实例名称
	Addr     string     `yaml:"addr"`     // 地址
	Port     int        `yaml:"port"`     // 端口
	Password string     `yaml:"password"` // 密码
	DBs      []DBConfig `yaml:"dbs"`      // 数据库
}

type DBConfig struct {
	DB       int `yaml:"db"`
	PoolSize int `yaml:"pool_size"`
}

type Zap struct {
	Director   string // 日志文件夹
	Level      string // 日志级别
	MaxAge     int    // 日志保存天数
	MaxSize    int    // 日志大小(MB)
	MaxBackups int    // 日志备份数量
	Format     string // 输出日志格式
}

type Config struct {
	Debug        string   // 调试模式
	Port         int      // 端口
	Limit        float64  // 限流
	AllowOrigins []string `yaml:"allowOrigins"` // 允许跨域origin
	Redis        struct { // Redis配置
		Instances []RedisInstanceConfig // Redis实例配置
	}
	Zap Zap // 日志
}

const configFilePath = "./config.yaml"

var (
	AppConfig = &Config{}
	port      int
)

func InitConfig() error {
	// 设置命令行参数
	flag.IntVar(&port, "port", 0, "指定端口号")
	flag.Parse()

	// 读取配置文件
	viper.SetConfigFile(configFilePath)
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
		fmt.Println(fmt.Sprintf("使用命令行设置的端口:%v", port))
	}
	return nil
}
