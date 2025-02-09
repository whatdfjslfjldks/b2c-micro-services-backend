package config

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v3"
	"log"
	"micro-services/pkg/utils"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// MysqlConfig  配置结构体
type MysqlConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Dbname   string `yaml:"dbname"`
}

type TitleMConfig struct {
	Mysql MysqlConfig `yaml:"mysql-user-server"`
}

var (
	MSConfig    *TitleMConfig
	CtxM        = context.Background()
	MySqlClient *sql.DB
)

func InitMysqlConfig() error {
	rootPath := utils.GetCurrentPath(2)
	configPath := filepath.Join(rootPath, "../pkg/config", "config.yml")
	MSConfig = &TitleMConfig{}

	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, MSConfig)
	if err != nil {
		return err
	}
	return nil
}

func InitMySql() {
	dsn := MSConfig.Mysql.User + ":" + MSConfig.Mysql.Password + "@tcp(" + MSConfig.Mysql.Host + ":" + strconv.Itoa(MSConfig.Mysql.Port) + ")/" + MSConfig.Mysql.Dbname
	var err error
	//fmt.Println("sdfd : ", dsn)
	MySqlClient, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	// 设置连接池参数
	MySqlClient.SetMaxOpenConns(25)
	MySqlClient.SetMaxIdleConns(25)
	MySqlClient.SetConnMaxLifetime(5 * time.Minute)

	// 测试连接
	err = MySqlClient.Ping()
	if err != nil {
		log.Fatal(err)
	}
}
