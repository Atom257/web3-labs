package repository

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/Atom257/web3-labs/timeledger-backend/internal/config"
	"github.com/Atom257/web3-labs/timeledger-backend/internal/models"
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

	// 自动创建或更新表结构
	if err := autoMigrate(db); err != nil {
		return nil, fmt.Errorf("自动建表失败: %w", err)
	}

	return db, nil
}

// autoMigrate 自动创建或更新数据库表结构
// 根据模型定义自动生成表和索引
func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.BlockCursor{},  // 区块同步游标
		&models.BlockHeader{},  // 区块头信息
		&models.BalanceLog{},   // 余额变动日志
		&models.UserBalance{},  // 用户余额快照
		&models.UserPoint{},    // 用户积分
		&models.UserPointLog{}, // 用户积分日志
		&models.PointRate{},    // 积分率配置
	)
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
