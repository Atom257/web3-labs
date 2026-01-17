package models

import "time"

type BlockHeader struct {
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	ChainID         int64  `gorm:"not null;index:idx_chain_contract_block,priority:1"`
	ContractAddress string `gorm:"type:char(42);not null;index:idx_chain_contract_block,priority:2"`

	BlockNumber int64  `gorm:"not null;index:idx_chain_contract_block,priority:3"`
	BlockHash   string `gorm:"type:char(66);not null"`
	ParentHash  string `gorm:"type:char(66);not null"`

	BlockTime time.Time `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
}

func (BlockHeader) TableName() string { return "block_header" }
