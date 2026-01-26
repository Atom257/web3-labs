package models

import (
	"fmt"
	"time"
)

// SysChain 链配置表
// 用于存储每条链的通用配置 (RPC, 确认数等)
type SysChain struct {
	ID      uint64 `gorm:"primaryKey;autoIncrement"`
	ChainID int64  `gorm:"not null;uniqueIndex"`      // 真实的链ID (如  11155111)
	Name    string `gorm:"type:varchar(64);not null"` // 链名称 (如 sepolia)
	Type    string `gorm:"type:varchar(32);not null"` // 链类型 (ethereum | opstack)

	// RPC 配置
	RpcEnvKey string `gorm:"type:varchar(64)"`  // 对应环境变量名
	RpcUrl    string `gorm:"type:varchar(255)"` // 允许直接存 URL (未来扩展用)

	// 同步参数
	Confirmations  int `gorm:"default:6"`   // 确认区块数
	ChunkSize      int `gorm:"default:10"`  // 每次扫描块数
	RequestDelayMs int `gorm:"default:100"` // 请求间隔

	//OP Stack 必须要用的回滚窗口
	ReorgWindow int `gorm:"default:200"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (SysChain) TableName() string { return "sys_chains" }

// SysContract 合约配置表
// 用于存储需要索引的合约。
// 使用这张表的 ID 来命名分表。
type SysContract struct {
	// 这是一个自增 ID (1, 2, 3...)
	// 用来生成日志表名：user_point_log_1, user_point_log_2
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	ChainID int64 `gorm:"not null;index:uniq_chain_addr,unique,priority:1"` // 关联到 SysChain

	Name    string `gorm:"type:varchar(64)"`                                               // 合约别名 (如 "USDT-Pool")
	Address string `gorm:"type:char(42);not null;index:uniq_chain_addr,unique,priority:2"` // 合约地址

	// 业务配置
	StartBlock    int64 `gorm:"not null"`
	TokenDecimals int   `gorm:"default:18"`

	// 状态开关 (方便单独暂停某个合约的索引/计算)
	IsEnabled bool `gorm:"default:true;index"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (SysContract) TableName() string { return "sys_contracts" }

// ---------------------------------------------------------
// 辅助函数：生成动态表名
// ---------------------------------------------------------

// GetLogTableName 根据合约 ID 生成唯一的积分流水表名
// 例如 ID=5 -> "user_point_log_5"
func (c *SysContract) GetLogTableName() string {
	return fmt.Sprintf("user_point_log_%d", c.ID)
}
