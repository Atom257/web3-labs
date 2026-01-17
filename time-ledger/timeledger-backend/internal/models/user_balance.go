package models

import "time"

type UserBalance struct {
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	ChainID         int64  `gorm:"not null;index:uniq_account,unique;index:idx_block,priority:1"`
	ContractAddress string `gorm:"type:char(42);not null;index:uniq_account,unique;index:idx_block,priority:2"`
	Account         string `gorm:"type:char(42);not null;index:uniq_account,unique"`

	Balance string `gorm:"type:decimal(65,0);not null"`

	BlockNumber int64     `gorm:"not null;index:idx_block,priority:3"`
	BlockTime   time.Time `gorm:"type:datetime(6);not null"`

	UpdatedAt time.Time `gorm:"type:datetime(6);not null;autoUpdateTime"`
}

func (UserBalance) TableName() string { return "user_balance" }
