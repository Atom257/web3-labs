package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/Atom257/web3-labs/timeledger-backend/internal/config"
	"github.com/Atom257/web3-labs/timeledger-backend/internal/models"
)

// InitSystem 系统初始化入口
// 职责：统一建表 -> 同步配置 -> 初始化默认规则
func InitSystem(ctx context.Context, db *gorm.DB, cfg *config.Config) error {
	log.Println("[Init] 开始系统初始化...")

	// ---------------------------------------------------------
	// 1. 统一建表 (AutoMigrate) - 仅限静态表
	// ---------------------------------------------------------
	log.Println("[Init] 正在检查并迁移数据库表结构...")

	sysTables := []interface{}{
		&models.SysChain{},
		&models.SysContract{},
		&models.PointRate{},
		&models.BlockCursor{},
		&models.BlockHeader{},
		&models.BalanceLog{},
		&models.UserBalance{},
		&models.UserPoint{},
		// 注意：不包含 UserPointLog，因为它是动态表
	}

	if err := db.AutoMigrate(sysTables...); err != nil {
		return fmt.Errorf("数据库表结构迁移失败: %w", err)
	}

	// ---------------------------------------------------------
	// 2. 数据初始化与动态建表
	// ---------------------------------------------------------
	// 同步 config.toml -> DB，并创建 user_point_log_X
	if err := syncSysConfig(ctx, db, cfg); err != nil {
		return err
	}

	return nil
}

// syncSysConfig 同步配置表 + 动态创建分表
func syncSysConfig(ctx context.Context, db *gorm.DB, cfg *config.Config) error {
	for _, chainCfg := range cfg.Chains {
		// --- Sync Chain ---
		sysChain := models.SysChain{
			ChainID:        chainCfg.ChainID,
			Name:           chainCfg.Name,
			Type:           chainCfg.Type,
			RpcEnvKey:      chainCfg.RPCEnvKey,
			ReorgWindow:    int(chainCfg.ReorgWindow),
			ChunkSize:      int(chainCfg.ChunkSize),
			RequestDelayMs: int(chainCfg.RequestDelayMs),
		}

		if err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "chain_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"name", "type", "reorg_window", "chunk_size", "request_delay_ms"}),
		}).Create(&sysChain).Error; err != nil {
			return fmt.Errorf("同步 Chain %d 失败: %w", chainCfg.ChainID, err)
		}

		// --- Sync Contracts ---
		for _, contractCfg := range chainCfg.Contracts {
			// 1. 准备数据
			sysContract := models.SysContract{
				ChainID:       chainCfg.ChainID,
				Address:       contractCfg.Address,
				StartBlock:    contractCfg.StartBlock,
				TokenDecimals: int(contractCfg.TokenDecimals),
				Name:          "Default-Pool",
				IsEnabled:     true,
				CreatedAt:     time.Now().UTC(),
			}

			// 2. Upsert (我们需要拿到 ID)
			// 先查一下是否存在，为了获取 ID
			var existing models.SysContract
			err := db.Where("chain_id = ? AND address = ?", chainCfg.ChainID, contractCfg.Address).First(&existing).Error

			if err == nil {
				// 存在：更新
				sysContract.ID = existing.ID
				if err := db.Model(&existing).Updates(map[string]interface{}{
					"start_block":    contractCfg.StartBlock,
					"token_decimals": int(contractCfg.TokenDecimals),
					"is_enabled":     true,
				}).Error; err != nil {
					return err
				}
			} else {
				// 不存在：插入
				if err := db.Create(&sysContract).Error; err != nil {
					return err
				}
				// 此时 sysContract.ID 已经被 GORM 填充
			}

			// 3. 【关键】动态创建分表 user_point_log_{id}
			// -------------------------------------------------------
			logTableName := sysContract.GetLogTableName() // 获取表名，如 user_point_log_1

			// 检查表是否存在，不存在则创建
			if !db.Migrator().HasTable(logTableName) {
				log.Printf("[Init] 正在创建动态分表: %s", logTableName)
				// 使用 UserPointLog 结构体做模板，但指定 TableName
				if err := db.Table(logTableName).AutoMigrate(&models.UserPointLog{}); err != nil {
					return fmt.Errorf("create dynamic table %s failed: %w", logTableName, err)
				}
			}

			// 4. 初始化默认积分规则 (如果是新合约)
			if err := initDefaultRates(db, sysContract); err != nil {
				log.Printf("[Warn] init default rate failed: %v", err)
			}
		}
	}
	log.Println("[Init] 系统配置同步及动态建表完成")
	return nil
}

// initDefaultRates 确保默认积分规则存在
func initDefaultRates(db *gorm.DB, contract models.SysContract) error {
	var count int64
	db.Model(&models.PointRate{}).
		Where("chain_id = ? AND contract_address = ?", contract.ChainID, contract.Address).
		Count(&count)

	if count == 0 {
		log.Printf("[Init] 为合约 %s 创建初始积分规则", contract.Address)
		rate := models.PointRate{
			ChainID:         contract.ChainID,
			ContractAddress: contract.Address,
			EffectiveTime:   time.Unix(0, 0).UTC(),
			RateNumerator:   5,
			RateDenominator: 100,
			CreatedAt:       time.Now().UTC(),
		}
		return db.Create(&rate).Error
	}
	return nil
}
