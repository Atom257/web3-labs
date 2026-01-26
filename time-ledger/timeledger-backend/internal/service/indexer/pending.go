package indexer

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"

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

	// pending 区块设置不设置过期时间，靠主动清除
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
