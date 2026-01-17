package models

import "time"

type BalanceLog struct {
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	ChainID         int64  `gorm:"not null;index:uniq_event,priority:1;index:idx_account_time,priority:1;index:idx_block,priority:1"`
	ContractAddress string `gorm:"type:char(42);not null;index:uniq_event,priority:2;index:idx_account_time,priority:2;index:idx_block,priority:2"`

	Account string `gorm:"type:char(42);not null;index:uniq_event,priority:3;index:idx_account_time,priority:3"`

	Delta        string `gorm:"type:decimal(65,0);not null"`
	BalanceAfter string `gorm:"type:decimal(65,0);not null"`

	BlockNumber int64     `gorm:"not null;index:uniq_event,priority:4;index:idx_block,priority:3"`
	BlockTime   time.Time `gorm:"type:datetime(6);not null"`

	TxHash   string `gorm:"type:char(66);not null"`
	LogIndex int64  `gorm:"not null;index:uniq_event,priority:5"`

	CreatedAt time.Time `gorm:"type:datetime(6);not null;autoCreateTime"`
}

func (BalanceLog) TableName() string { return "balance_log" }
