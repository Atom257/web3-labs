package repository

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/Atom257/web3-labs/timeledger-backend/internal/config"
)

// InitDB 初始化数据库连接
// 从环境变量读取数据库配置，建立连接并自动创建表结构
func InitDB(cfg config.DatabaseConfig) (*gorm.DB, error) {
	// 从环境变量读取数据库连接信息
	dbUser, err := getEnv("DB_USER")
	if err != nil {
		return nil, err
	}
	dbPassword, err := getEnv("DB_PASSWORD")
	if err != nil {
		return nil, err
	}
	dbHost, err := getEnv("DB_HOST")
	if err != nil {
		return nil, err
	}
	dbPort, err := getEnv("DB_PORT")
	if err != nil {
		return nil, err
	}
	dbName, err := getEnv("DB_NAME")
	if err != nil {
		return nil, err
	}

	// 构造 DSN（数据源名称）
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=UTC",
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
		dbName,
	)

	// 打开数据库连接
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("打开数据库连接失败: %w", err)
	}

	// 获取底层的 sql.DB 对象，用于设置连接池参数
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 设置连接池参数
	if cfg.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	sqlDB.SetConnMaxLifetime(1 * time.Hour)

	return db, nil
}

// getEnv 从环境变量获取配置值
// 如果环境变量不存在，返回错误而不是 panic
func getEnv(key string) (string, error) {
	v := os.Getenv(key)
	if v == "" {
		return "", fmt.Errorf("缺少必需的环境变量: %s", key)
	}
	return v, nil
}
