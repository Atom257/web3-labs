package models

import "time"

type PointRate struct {
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	ChainID         int64  `gorm:"not null;index:uniq_rate_time,unique;index:idx_rate_time,priority:1"`
	ContractAddress string `gorm:"type:char(42);not null;index:uniq_rate_time,unique;index:idx_rate_time,priority:2"`

	RateNumerator   int64 `gorm:"not null"`
	RateDenominator int64 `gorm:"not null"`

	EffectiveTime time.Time `gorm:"type:datetime(6);not null;index:uniq_rate_time,unique;index:idx_rate_time,priority:3"`
	CreatedAt     time.Time `gorm:"type:datetime(6);not null;autoCreateTime"`
}

func (PointRate) TableName() string { return "point_rate" }
