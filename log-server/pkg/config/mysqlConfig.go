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

// MysqlConfig 配置结构体
type MysqlConfig struct {
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	User        string `yaml:"user"`
	Password    string `yaml:"password"`
	DbnameInfo  string `yaml:"dbname_info"`
	DbnameWarn  string `yaml:"dbname_warn"`
	DbnameError string `yaml:"dbname_error"`
}

type TitleMConfig struct {
	Mysql MysqlConfig `yaml:"mysql-log-server"`
}

var (
	MSConfig         *TitleMConfig
	CtxM             = context.Background()
	MySqlInfoClient  *sql.DB
	MySqlWarnClient  *sql.DB
	MySqlErrorClient *sql.DB
)

// InitMysqlConfig 读取配置文件并初始化配置结构体
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

// InitMySql 初始化 MySQL 客户端连接
func InitMySql() {
	// 配置数据库 DSN 模板
	dsnTemplate := MSConfig.Mysql.User + ":" + MSConfig.Mysql.Password + "@tcp(" + MSConfig.Mysql.Host + ":" + strconv.Itoa(MSConfig.Mysql.Port) + ")/"

	// 初始化 info、warn、error 数据库客户端
	var err error
	MySqlInfoClient, err = initMySqlClient(dsnTemplate + MSConfig.Mysql.DbnameInfo)
	if err != nil {
		log.Fatal("Error initializing Info DB: ", err)
	}

	MySqlWarnClient, err = initMySqlClient(dsnTemplate + MSConfig.Mysql.DbnameWarn)
	if err != nil {
		log.Fatal("Error initializing Warn DB: ", err)
	}

	MySqlErrorClient, err = initMySqlClient(dsnTemplate + MSConfig.Mysql.DbnameError)
	if err != nil {
		log.Fatal("Error initializing Error DB: ", err)
	}
}

// initMySqlClient 用于初始化数据库客户端的通用方法
func initMySqlClient(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// 设置连接池参数
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// 测试连接
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
