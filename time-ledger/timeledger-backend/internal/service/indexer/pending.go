package indexer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/Atom257/web3-labs/timeledger-backend/internal/config"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/redis/go-redis/v9"
)

type PendingHead struct {
	ChainID         int64     `json:"chain_id"`
	ContractAddress string    `json:"contract_address"`
	BlockNumber     uint64    `json:"block_number"`
	BlockHash       string    `json:"block_hash"`
	ParentHash      string    `json:"parent_hash"`
	BlockTime       time.Time `json:"block_time"`
	Finalized       bool      `json:"finalized"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// PendingBlock 表示 OP Stack 上一个尚未确认的区块数据，仅用于 Redis 缓存
type PendingBlock struct {
	ChainID         int64  `json:"chain_id"`
	ContractAddress string `json:"contract_address"`

	BlockNumber uint64    `json:"block_number"`
	BlockHash   string    `json:"block_hash"`
	ParentHash  string    `json:"parent_hash"`
	BlockTime   time.Time `json:"block_time"`

	// 区块内的 Transfer 事件
	Events []TransferEvent `json:"events"`

	CreatedAt time.Time `json:"created_at"`
}

// UpdatePendingHead 将 OP Stack 当前最新 head 写入 Redis，仅用于观测链状态
func (ix *Indexer) UpdatePendingHead(
	ctx context.Context,
	rdb *redis.Client,
	client *ethclient.Client,
	chainID int64,
	contractAddr string,
) error {

	// 获取链上最新的区块头
	h, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		return err
	}

	ph := PendingHead{
		ChainID:         chainID,
		ContractAddress: contractAddr,
		BlockNumber:     h.Number.Uint64(),
		BlockHash:       h.Hash().Hex(),
		ParentHash:      h.ParentHash.Hex(),
		BlockTime:       time.Unix(int64(h.Time), 0).UTC(),
		Finalized:       false,
		UpdatedAt:       time.Now().UTC(),
	}

	b, _ := json.Marshal(ph)
	key := fmt.Sprintf(
		"%s:pending:head:%d:%s",
		ix.cfg.Redis.KeyPrefix,
		chainID,
		contractAddr,
	)

	return rdb.Set(ctx, key, b, 10*time.Second).Err()
}

// StagePendingBlock 将未确认区块写入 Redis，用于 OP Stack 数据暂存
func (ix *Indexer) StagePendingBlock(
	ctx context.Context,
	rdb *redis.Client,
	block *PendingBlock,
) error {

	key := fmt.Sprintf(
		"%s:pending:block:%d:%s:%d",
		ix.cfg.Redis.KeyPrefix, // 需要前缀
		block.ChainID,
		block.ContractAddress,
		block.BlockNumber,
	)

	data, err := json.Marshal(block)
	if err != nil {
		return err
	}

	// pending 区块设置短 不过期。 防止时间过短还没落库已经被清楚
	//  redis的清除 靠主动清除。
	return rdb.Set(ctx, key, data, 0).Err()
}

// ListPendingBlocksUpTo 读取指定高度之前的 pending 区块
func (ix *Indexer) ListPendingBlocksUpTo(
	ctx context.Context,
	rdb *redis.Client,
	chainID int64,
	contract string,
	maxBlock uint64,
) ([]*PendingBlock, error) {

	pattern := fmt.Sprintf(
		"%s:pending:block:%d:%s:*",
		ix.cfg.Redis.KeyPrefix,
		chainID,
		contract,
	)

	var blocks []*PendingBlock

	// 用 SCAN 代替 KEYS，避免 Redis 阻塞
	iter := rdb.Scan(ctx, 0, pattern, 200).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()

		val, err := rdb.Get(ctx, key).Bytes()
		if err != nil {
			continue
		}

		var pb PendingBlock
		if err := json.Unmarshal(val, &pb); err != nil {
			continue
		}

		if pb.BlockNumber <= maxBlock {
			blocks = append(blocks, &pb)
		}
	}

	if err := iter.Err(); err != nil {
		return nil, err
	}

	sort.Slice(blocks, func(i, j int) bool {
		return blocks[i].BlockNumber < blocks[j].BlockNumber
	})

	return blocks, nil
}

// FlushSafePending 将 Redis 中 <= 当前 safe block 的 pending 区块
// 按区块高度顺序 apply 到数据库。
//
// 设计保证：
// - 可重入（多次调用安全）
// - 仅推进 canonical block（block_number / block_hash）
// - 不修改 scan_block_number（扫描进度）
// - 与 reorg 逻辑兼容
func (ix *Indexer) FlushSafePending(
	ctx context.Context,
	client *ethclient.Client,
	adapter ChainAdapter,
	chain config.ChainConfig,
	contract config.ContractConfig,
	dbBlock *int64,
) error {

	// 重新获取最新 safe block（不能复用旧 snapshot）
	safeBlock, err := adapter.SafeBlock(ctx, client, chain)
	if err != nil {
		return err
	}

	// 从 Redis 读取 <= safeBlock 的 pending 区块
	pendingBlocks, err := ix.ListPendingBlocksUpTo(
		ctx,
		ix.redis,
		chain.ChainID,
		contract.Address,
		safeBlock,
	)
	if err != nil {
		return err
	}

	if len(pendingBlocks) == 0 {
		return nil
	}

	// 确保按 block_number 升序（ListPendingBlocksUpTo 已排序，这里防御性再排一次）
	sort.Slice(pendingBlocks, func(i, j int) bool {
		return pendingBlocks[i].BlockNumber < pendingBlocks[j].BlockNumber
	})

	for _, pb := range pendingBlocks {

		// 已经 apply 过的 block，直接跳过（防御）
		if int64(pb.BlockNumber) <= *dbBlock {
			continue
		}

		headers := map[uint64]*blockHeaderMini{
			pb.BlockNumber: {
				Number: pb.BlockNumber,
				Hash:   common.HexToHash(pb.BlockHash),
				Parent: common.HexToHash(pb.ParentHash),
				Time:   pb.BlockTime,
			},
		}

		log.Printf(
			"[opstack.flush.safe] chain=%d contract=%s block=%d events=%d safe=%d",
			chain.ChainID,
			contract.Address,
			pb.BlockNumber,
			len(pb.Events),
			safeBlock,
		)

		// apply 到数据库（会推进 block_cursor.block_number / block_hash）
		if err := ix.applyChunkTx(
			ctx,
			client,
			chain.ChainID,
			contract.Address,
			pb.BlockNumber,
			pb.BlockNumber,
			pb.Events,
			headers,
		); err != nil {
			return err
		}

		// 更新外部 dbBlock（供调用方使用）
		*dbBlock = int64(pb.BlockNumber)

		// apply 成功后，删除 Redis pending
		key := fmt.Sprintf(
			"%s:pending:block:%d:%s:%d",
			ix.cfg.Redis.KeyPrefix,
			chain.ChainID,
			contract.Address,
			pb.BlockNumber,
		)
		_ = ix.redis.Del(ctx, key).Err()
	}

	return nil
}

// CleanupPendingAfterReorg
// 清理 Redis 中 reorg 发生点之后的 pending 区块数据
func (ix *Indexer) CleanupPendingAfterReorg(
	ctx context.Context,
	rdb *redis.Client,
	chainID int64,
	contract string,
	ancestor int64,
) {

	if rdb == nil {
		return
	}

	pattern := fmt.Sprintf(
		"pending:block:%d:%s:*",
		chainID,
		contract,
	)

	// 用 SCAN 代替 KEYS，避免 Redis 阻塞
	iter := rdb.Scan(ctx, 0, pattern, 200).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()

		val, err := rdb.Get(ctx, key).Bytes()
		if err != nil {
			continue
		}

		var pb PendingBlock
		if err := json.Unmarshal(val, &pb); err != nil {
			continue
		}

		if int64(pb.BlockNumber) > ancestor {
			_ = rdb.Del(ctx, key).Err()
		}
	}

	_ = iter.Err()
}
