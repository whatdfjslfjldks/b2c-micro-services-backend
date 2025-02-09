package config

import (
	"context"
	"github.com/go-redis/redis/v8"
	"gopkg.in/yaml.v3"
	"micro-services/pkg/utils"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// RedisConfig Redis 配置结构体
type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	Db       int    `yaml:"db"`
}

type TitleConfig struct {
	Redis RedisConfig `yaml:"redis"`
}

var (
	RdConfig *TitleConfig
	Ctx      = context.Background()
	RdClient *redis.Client
)

// InitRedisConfig 初始化配置 Redis
func InitRedisConfig() error {
	rootPath := utils.GetCurrentPath(2)
	configPath := filepath.Join(rootPath, "../pkg/config", "config.yml")
	RdConfig = &TitleConfig{}

	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, RdConfig)
	if err != nil {
		return err
	}
	return nil
}

// InitRedis 初始化 Redis 客户端
func InitRedis() {
	// 设置超时时间 5 秒
	dialTimeout := 5 * time.Second
	readTimeout := 5 * time.Second
	writeTimeout := 5 * time.Second

	// 创建 Redis-0 客户端
	RdClient = redis.NewClient(&redis.Options{
		Addr:         RdConfig.Redis.Host + ":" + strconv.Itoa(RdConfig.Redis.Port), // Redis 服务器地址
		Password:     RdConfig.Redis.Password,                                       // Redis 密码，如果没有就留空
		DB:           5,                                                             // 使用的数据库索引，5
		DialTimeout:  dialTimeout,                                                   // 连接超时
		ReadTimeout:  readTimeout,                                                   // 读取超时
		WriteTimeout: writeTimeout,                                                  // 写入超时
	})

	// 测试连接
	_, err := RdClient.Ping(Ctx).Result()
	if err != nil {
		panic(err)
	}
}
