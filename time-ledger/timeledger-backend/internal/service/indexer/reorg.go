package indexer

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"gorm.io/gorm"

	"github.com/Atom257/web3-labs/timeledger-backend/internal/models"
)

// EnsureCanonicalOrRollback：
// - 对 opstack：用 block_header / chain header hash 比较
// - 不一致就 rollback 到 common ancestor
func (ix *Indexer) EnsureCanonicalOrRollback(
	ctx context.Context,
	client *ethclient.Client,
	chainID int64,
	contractAddr string,
	reorgWindow int64,
) error {

	if reorgWindow <= 0 {
		return fmt.Errorf("reorgWindow must be > 0")
	}

	// 读 cursor
	var cur models.BlockCursor
	err := ix.db.
		Where("chain_id=? AND contract_address=?", chainID, contractAddr).
		First(&cur).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 第一次启动，还没 cursor，不可能发生 reorg
		return nil
	}
	if err != nil {
		return err
	}

	// cursor 还没真正同步过区块，但已经有 block_number，没有 block_hash
	// 这种状态不能做 reorg 判断
	if cur.BlockNumber > 0 && cur.BlockHash == "" {
		return nil
	}

	n := cur.BlockNumber
	if n <= 0 {
		// 还没处理到任何块，无需校验
		return nil
	}

	// 确认 DB 里已经有 block_header
	var cnt int64
	if err := ix.db.
		Model(&models.BlockHeader{}).
		Where(
			"chain_id=? AND contract_address=? AND block_number=?",
			chainID,
			contractAddr,
			n,
		).
		Count(&cnt).Error; err != nil {
		return err
	}

	// DB 里还没有任何 block_header，说明仍处于 pending 阶段
	// OP Stack 下这是正常情况，不能做 reorg 判断
	if cnt == 0 {
		return nil
	}

	// 找 common ancestor
	ancestor, ancestorHash, err := ix.findCommonAncestor(ctx, client, chainID, contractAddr, n, reorgWindow)
	if err != nil {
		return err
	}

	if ancestor == n {
		// 没有分叉
		return nil
	}

	// 记录 OP Stack reorg 发生
	if ix.redis != nil {
		key := fmt.Sprintf(
			"reorg:seen:%d:%s",
			chainID,
			contractAddr,
		)
		_ = ix.redis.Set(
			ctx,
			key,
			time.Now().UTC().Format(time.RFC3339),
			time.Hour,
		)
	}

	// 发生分叉，执行 rollback
	return ix.rollbackTo(ctx, client, chainID, contractAddr, ancestor, ancestorHash)
}

func (ix *Indexer) findCommonAncestor(
	ctx context.Context,
	client *ethclient.Client,
	chainID int64,
	contractAddr string,
	cursor int64,
	reorgWindow int64,
) (int64, string, error) {

	low := cursor - reorgWindow
	if low < 0 {
		low = 0
	}

	for bn := cursor; bn >= low; bn-- {
		// DB hash
		var bh models.BlockHeader
		err := ix.db.
			Select("block_hash").
			Where("chain_id=? AND contract_address=? AND block_number=?",
				chainID, contractAddr, bn,
			).
			First(&bh).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			continue
		}
		if err != nil {
			return 0, "", err
		}

		// chain hash
		hdr, err := callRPCWithRetry(
			ctx,
			ix.rpcLimiter,
			"eth_getBlockByNumber",
			chainID,
			uint64(bn),
			func() (*types.Header, error) {
				return client.HeaderByNumber(ctx, big.NewInt(bn))
			},
		)

		if err != nil {
			return 0, "", err
		}

		chainHash := hdr.Hash().Hex()
		if bh.BlockHash == chainHash {
			return bn, chainHash, nil
		}
	}

	return 0, "", fmt.Errorf("no common ancestor within reorg_window=%d blocks", reorgWindow)
}

func (ix *Indexer) rollbackTo(
	ctx context.Context,
	client *ethclient.Client,
	chainID int64,
	contractAddr string,
	ancestor int64,
	ancestorHash string,
) error {

	// 1️⃣ 先执行数据库回滚（原有逻辑，完全不动）
	err := ix.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		// 删除 fork 段 block_header
		if err := tx.Where(
			"chain_id=? AND contract_address=? AND block_number > ?",
			chainID, contractAddr, ancestor,
		).Delete(&models.BlockHeader{}).Error; err != nil {
			return err
		}

		// 删除 fork 段 balance_log
		if err := tx.Where(
			"chain_id=? AND contract_address=? AND block_number > ?",
			chainID, contractAddr, ancestor,
		).Delete(&models.BalanceLog{}).Error; err != nil {
			return err
		}

		// 删除 fork 段对应时间之后的积分日志
		var ancHeader models.BlockHeader
		if err := tx.
			Select("block_time").
			Where(
				"chain_id=? AND contract_address=? AND block_number=?",
				chainID, contractAddr, ancestor,
			).
			First(&ancHeader).Error; err != nil {
			return err
		}
		ancestorTime := ancHeader.BlockTime.UTC()

		if err := tx.Where(
			"chain_id=? AND contract_address=? AND to_time > ?",
			chainID, contractAddr, ancestorTime,
		).Delete(&models.UserPointLog{}).Error; err != nil {
			return err
		}

		// 重建 user_balance
		if err := tx.Where(
			"chain_id=? AND contract_address=?",
			chainID, contractAddr,
		).Delete(&models.UserBalance{}).Error; err != nil {
			return err
		}

		type row struct {
			Account      string
			BalanceAfter string
			BlockNumber  int64
			BlockTime    time.Time
		}

		var rows []row
		sub := tx.Model(&models.BalanceLog{}).
			Select("account, MAX(block_number) AS max_bn").
			Where("chain_id=? AND contract_address=?", chainID, contractAddr).
			Group("account")

		if err := tx.Table("balance_log bl").
			Select("bl.account, bl.balance_after, bl.block_number, bl.block_time").
			Joins("JOIN (?) m ON bl.account = m.account AND bl.block_number = m.max_bn", sub).
			Where("bl.chain_id=? AND bl.contract_address=?", chainID, contractAddr).
			Scan(&rows).Error; err != nil {
			return err
		}

		now := time.Now().UTC()
		for _, r := range rows {
			ub := models.UserBalance{
				ChainID:         chainID,
				ContractAddress: contractAddr,
				Account:         r.Account,
				Balance:         r.BalanceAfter,
				BlockNumber:     r.BlockNumber,
				BlockTime:       r.BlockTime.UTC(),
				UpdatedAt:       now,
			}
			if err := tx.Create(&ub).Error; err != nil {
				return err
			}
		}

		// 重建 user_point
		if err := tx.Where(
			"chain_id=? AND contract_address=?",
			chainID, contractAddr,
		).Delete(&models.UserPoint{}).Error; err != nil {
			return err
		}

		type pointAggRow struct {
			Account      string
			TotalPoints  string
			LastCalcTime time.Time
		}

		var pointRows []pointAggRow
		if err := tx.Table("user_point_log AS upl").
			Select(`
				upl.account AS account,
				SUM(upl.points) AS total_points,
				MAX(upl.to_time) AS last_calc_time
			`).
			Where("upl.chain_id=? AND upl.contract_address=?", chainID, contractAddr).
			Group("upl.account").
			Scan(&pointRows).Error; err != nil {
			return err
		}

		for _, r := range pointRows {
			up := models.UserPoint{
				ChainID:         chainID,
				ContractAddress: contractAddr,
				Account:         r.Account,
				TotalPoints:     r.TotalPoints,
				LastCalcTime:    r.LastCalcTime.UTC(),
				UpdatedAt:       now,
			}
			if err := tx.Create(&up).Error; err != nil {
				return err
			}
		}

		// 更新 cursor（最后一步）
		if err := tx.Model(&models.BlockCursor{}).
			Where("chain_id=? AND contract_address=?", chainID, contractAddr).
			Updates(map[string]any{
				"block_number":      ancestor,
				"block_hash":        ancestorHash,
				"scan_block_number": ancestor, // 必须回退
				"updated_at":        now,
			}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	//	DB 回滚成功后，再清理 Redis（新增代码）
	if ix.redis != nil {
		ix.CleanupPendingAfterReorg(
			ctx,
			ix.redis,
			chainID,
			contractAddr,
			ancestor,
		)
	}

	return nil
}
