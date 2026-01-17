package models

import "time"

// UserPointLog
// 积分事实表（event-sourced）
type UserPointLog struct {
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	ChainID         int64  `gorm:"not null;index:uniq_point_log,unique,priority:1"`
	ContractAddress string `gorm:"type:char(42);not null;index:uniq_point_log,unique,priority:2"`
	Account         string `gorm:"type:char(42);not null;index:uniq_point_log,unique,priority:3"`

	//当前余额
	Balance string `gorm:"type:decimal(38,18);not null"`

	// 本次积分计算区间
	FromTime time.Time `gorm:"type:datetime(6);not null;index:uniq_point_log,unique,priority:4"`
	ToTime   time.Time `gorm:"type:datetime(6);not null;index:uniq_point_log,unique,priority:5"`

	// 本区间产生的积分（decimal 字符串）
	Points string `gorm:"type:decimal(38,18);not null"`

	// 本次计算使用的积分规则
	RateNumerator   int64 `gorm:"not null"`
	RateDenominator int64 `gorm:"not null"`

	CreatedAt time.Time `gorm:"type:datetime(6);not null;autoCreateTime"`
}

func (UserPointLog) TableName() string {
	return "user_point_log"
}
