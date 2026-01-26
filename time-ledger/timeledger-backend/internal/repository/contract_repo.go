package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/Atom257/web3-labs/timeledger-backend/internal/models"
)

// GetActiveContracts 获取所有启用的合约任务
// Calculator 和 Indexer 都会用到这个方法
func GetActiveContracts(ctx context.Context, db *gorm.DB) ([]models.SysContract, error) {
	var contracts []models.SysContract

	// 因为 Indexer 和 Calculator 都会单独去查 chain 配置，或者只需要 ChainID 就够了
	err := db.WithContext(ctx).
		Where("is_enabled = ?", true).
		Find(&contracts).Error

	return contracts, err
}
