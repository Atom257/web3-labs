package models

import "time"

type UserPoint struct {
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	ChainID         int64  `gorm:"not null;index:uniq_account,unique;index:idx_last_calc_time,priority:1"`
	ContractAddress string `gorm:"type:char(42);not null;index:uniq_account,unique;index:idx_last_calc_time,priority:2"`
	Account         string `gorm:"type:char(42);not null;index:uniq_account,unique"`

	TotalPoints  string    `gorm:"type:decimal(38,18);not null"`
	LastCalcTime time.Time `gorm:"type:datetime(6);not null"`

	UpdatedAt time.Time `gorm:"type:datetime(6);not null;autoUpdateTime"`
}

func (UserPoint) TableName() string { return "user_point" }
