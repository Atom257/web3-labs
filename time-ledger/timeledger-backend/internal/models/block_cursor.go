package models

import "time"

type BlockCursor struct {
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	ChainID         int64  `gorm:"not null;index:uniq_chain_contract,unique"`
	ContractAddress string `gorm:"type:char(42);not null;index:uniq_chain_contract,unique"`

	//	DB 中已经确定的最高块

	BlockNumber int64  `gorm:"not null"`
	BlockHash   string `gorm:"type:char(66);not null"`

	//	已完成 RPC 扫描的最高块（允许回退）
	ScanBlockNumber int64 `gorm:"not null;default:0"`

	UpdatedAt time.Time `gorm:"type:datetime(6);not null;autoUpdateTime"`
}

func (BlockCursor) TableName() string { return "block_cursor" }
