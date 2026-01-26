package indexer

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/big"
	"sort"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/Atom257/web3-labs/timeledger-backend/internal/config"
	"github.com/Atom257/web3-labs/timeledger-backend/internal/models"
	"github.com/Atom257/web3-labs/timeledger-backend/internal/repository"
	erc20 "github.com/Atom257/web3-labs/timeledger-backend/pkg/contract/erc20"
)

/*
Indexer
- 按 chain / contract 扫描 ERC20 Transfer
- 写 block header / balance log / user balance
- 维护 block cursor
- 所有 RPC 统一走 limiter + retry
*/
type Indexer struct {
	db    *gorm.DB
	cfg   *config.Config
	redis *redis.Client

	// 全局 RPC 限速器
	rpcLimiter *RPCLimiter

	// 内存中的 scan cursor（key = chainID + contract）
	scanCache map[string]int64
	scanMu    sync.Mutex
}

func New(db *gorm.DB, cfg *config.Config, rdb *redis.Client) *Indexer {
	return &Indexer{
		db:         db,
		cfg:        cfg,
		redis:      rdb,
		rpcLimiter: NewRPCLimiter(3), // Alchemy 测试账号：3 RPS
		scanCache:  make(map[string]int64),
	}
}

/*
RunOnceConcurrent
-------
并行执行：
1. 从数据库获取 Active 的合约
2. 从 Config 获取 RPC URL (loader.go 已注入)
3. 按链分组并行执行
*/
func (ix *Indexer) RunOnceConcurrent(ctx context.Context) error {
	// 1. 获取所有待索引的合约 (DB)
	contracts, err := repository.GetActiveContracts(ctx, ix.db)
	if err != nil {
		return fmt.Errorf("load active contracts failed: %w", err)
	}

	if len(contracts) == 0 {
		return nil
	}

	// 2. 获取所有链配置 (DB) 用于获取 Type, ReorgWindow 等
	var chains []models.SysChain
	if err := ix.db.Find(&chains).Error; err != nil {
		return fmt.Errorf("load chains failed: %w", err)
	}
	sysChainMap := make(map[int64]models.SysChain)
	for _, c := range chains {
		sysChainMap[c.ChainID] = c
	}

	// 3. 构建 RPC URL 映射 (Memory Config)
	// loader.go 启动时已经把 env 里的 URL 填进 ix.cfg 了，直接用
	chainRpcMap := make(map[int64]string)
	for _, c := range ix.cfg.Chains {
		chainRpcMap[c.ChainID] = c.RPCURL
	}

	// 4. 按 ChainID 分组合约
	contractGroups := make(map[int64][]models.SysContract)
	for _, c := range contracts {
		contractGroups[c.ChainID] = append(contractGroups[c.ChainID], c)
	}

	// 5. 并发执行
	g, ctx := errgroup.WithContext(ctx)

	for chainID, targets := range contractGroups {
		chainID := chainID
		targets := targets

		// 获取 DB 里的链配置
		sysChain, ok := sysChainMap[chainID]
		if !ok {
			log.Printf("[Indexer] warn: chain_id=%d not found in sys_chains, skipping", chainID)
			continue
		}

		g.Go(func() error {
			// 获取 RPC URL
			rpcURL := chainRpcMap[chainID]
			if rpcURL == "" {
				// 如果 Config 里找不到 URL，说明 loader 没加载到，可能是新加的链没配 env
				return fmt.Errorf("chain %s (id=%d) missing rpc url in config", sysChain.Name, chainID)
			}

			// 初始化 Adapter
			adapter, err := AdapterFor(sysChain.Type)
			if err != nil {
				return err
			}

			// 连接 RPC
			client, err := ethclient.Dial(rpcURL)
			if err != nil {
				return fmt.Errorf("dial rpc failed chain=%d: %w", chainID, err)
			}
			defer client.Close()

			// 链内串行处理合约
			for _, contract := range targets {
				// 【关键】这里传的是 models.SysChain 和 models.SysContract
				if err := ix.syncContract(ctx, client, adapter, sysChain, contract); err != nil {
					// 遇到限流，中断该链
					if errors.Is(err, ErrRateLimited) {
						log.Printf("[indexer.exit] rate limited on chain %d", chainID)
						return err
					}
					// 其他错误，记录日志但不中断其他合约
					log.Printf("[Indexer] sync failed contract=%s: %v", contract.Address, err)
				}
			}
			return nil
		})
	}

	return g.Wait()
}

// syncContract 是 indexer 的核心编排函数
// 参数已全部修改为 models.SysChain 和 models.SysContract
func (ix *Indexer) syncContract(
	ctx context.Context,
	client *ethclient.Client,
	adapter ChainAdapter,
	chain models.SysChain,
	contract models.SysContract,
) error {

	log.Printf(
		"[ENTER syncContract] chain=%d type=%s contract=%s",
		chain.ChainID,
		chain.Type,
		contract.Address,
	)

	// 加载或初始化 cursor
	cursor, err := ix.loadOrInitCursor(chain.ChainID, contract)
	if err != nil {
		return err
	}

	// 确保函数结束时清理内存 scanCache
	defer func() {
		key := fmt.Sprintf(
			"%s:%d:%s",
			ix.cfg.Redis.KeyPrefix,
			chain.ChainID,
			contract.Address,
		)

		ix.scanMu.Lock()
		delete(ix.scanCache, key)
		ix.scanMu.Unlock()
	}()

	// 补齐初始化状态下缺失的 block_hash
	if err := ix.ensureCursorHash(ctx, client, chain, cursor); err != nil {
		return err
	}

	// 已确认的 canonical block
	dbBlock := cursor.BlockNumber

	// 已落库的 scan_block_number（用于控制 flush 间隔）
	scanFlushed := cursor.ScanBlockNumber

	// OP Stack：记录 pending head，仅用于观测
	if chain.Type == "opstack" && ix.redis != nil {
		_ = ix.UpdatePendingHead(ctx, ix.redis, client, chain.ChainID, contract.Address)
	}

	// OP Stack：reorg 检测与回滚
	if chain.Type == "opstack" {
		//  ReorgWindow 是 int，需转 int64
		if err := ix.EnsureCanonicalOrRollback(
			ctx,
			client,
			chain.ChainID,
			contract.Address,
			int64(chain.ReorgWindow),
		); err != nil {
			return err
		}
	}

	// 当前 safe block (Adapter 已适配 models.SysChain)
	safeBlock, err := adapter.SafeBlock(ctx, client, chain)
	if err != nil {
		return err
	}

	log.Printf(
		"[indexer.cursor] chain_id=%d contract=%s cursor=%d safe=%d",
		chain.ChainID,
		contract.Address,
		cursor.BlockNumber,
		safeBlock,
	)

	// 没有可扫描区间
	if uint64(cursor.BlockNumber) >= safeBlock {
		return nil
	}

	// ERC20 合约实例
	token, err := erc20.NewTimeLedgerToken(
		common.HexToAddress(contract.Address),
		client,
	)
	if err != nil {
		return err
	}

	// 计算扫描区间
	from, to := ix.computeScanRange(cursor, safeBlock)

	for start := from; start <= to; {

		end := ix.computeChunkEnd(start, to, chain)

		// RPC 限速
		if chain.RequestDelayMs > 0 {
			time.Sleep(time.Duration(chain.RequestDelayMs) * time.Millisecond)
		}

		// 拉取 Transfer 事件
		events, headers, err := ix.fetchTransfers(
			ctx,
			client,
			token,
			chain.ChainID,
			contract.Address,
			start,
			end,
		)
		if err != nil {
			return err
		}

		// 更新内存 scan 进度
		ix.updateScanProgress(chain, contract, end)

		// 周期性落库 scan_block_number
		if err := ix.flushScanCursor(ctx, client, chain, contract, &scanFlushed); err != nil {
			return err
		}

		if chain.Type == "opstack" && ix.redis != nil {
			// OP Stack：写 Redis pending
			if err := ix.handleOpStackChunk(
				ctx,
				client,
				adapter,
				chain,
				contract,
				headers,
				events,
				&dbBlock,
			); err != nil {
				return err
			}
		} else {
			// 非 OP Stack：直接落库
			if err := ix.applyNormalChunk(
				ctx,
				client,
				chain,
				contract,
				start,
				end,
				events,
				headers,
				&dbBlock,
			); err != nil {
				return err
			}
		}

		start = end + 1
	}

	// 扫描结束后兜底刷新 scan cursor
	if err := ix.flushScanCursor(ctx, client, chain, contract, &scanFlushed); err != nil {
		return err
	}

	// OP Stack：兜底 flush safe pending
	if chain.Type == "opstack" && ix.redis != nil {
		if err := ix.FlushSafePending(
			ctx,
			client,
			adapter,
			chain,
			contract,
			&dbBlock,
		); err != nil {
			return err
		}
	}

	return nil
}

// handleOpStackChunk 处理 OP Stack 扫描到的一个 chunk
func (ix *Indexer) handleOpStackChunk(
	ctx context.Context,
	client *ethclient.Client,
	adapter ChainAdapter,
	chain models.SysChain,
	contract models.SysContract,
	headers map[uint64]*blockHeaderMini,
	events []TransferEvent,
	dbBlock *int64,
) error {

	// 按 blockNumber 排序
	var blockNums []uint64
	for bn := range headers {
		blockNums = append(blockNums, bn)
	}
	sort.Slice(blockNums, func(i, j int) bool {
		return blockNums[i] < blockNums[j]
	})

	for _, bn := range blockNums {
		h := headers[bn]

		var blockEvents []TransferEvent
		for _, ev := range events {
			if ev.BlockNumber == h.Number {
				blockEvents = append(blockEvents, ev)
			}
		}

		pb := &PendingBlock{
			ChainID:         chain.ChainID,
			ContractAddress: contract.Address,
			BlockNumber:     h.Number,
			BlockHash:       h.Hash.Hex(),
			ParentHash:      h.Parent.Hex(),
			BlockTime:       h.Time,
			Events:          blockEvents,
			CreatedAt:       time.Now().UTC(),
		}

		if err := ix.StagePendingBlock(ctx, ix.redis, pb); err != nil {
			return err
		}
	}

	//  递归调用 FlushSafePending
	return ix.FlushSafePending(
		ctx,
		client,
		adapter,
		chain,
		contract,
		dbBlock,
	)
}

// FlushSafePending 将 Redis 中 <= 当前 safe block 的 pending 区块落库
func (ix *Indexer) FlushSafePending(
	ctx context.Context,
	client *ethclient.Client,
	adapter ChainAdapter,
	chain models.SysChain,
	contract models.SysContract,
	dbBlock *int64,
) error {

	// 1. 获取最新 safe block
	safeBlock, err := adapter.SafeBlock(ctx, client, chain)
	if err != nil {
		return err
	}

	// 2. 从 Redis 读取 pending blocks
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

	sort.Slice(pendingBlocks, func(i, j int) bool {
		return pendingBlocks[i].BlockNumber < pendingBlocks[j].BlockNumber
	})

	for _, pb := range pendingBlocks {
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

		// 3. 落库
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

		*dbBlock = int64(pb.BlockNumber)

		// 4. 清理 Redis
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

// applyNormalChunk 处理非 OP Stack 链的一个 chunk
func (ix *Indexer) applyNormalChunk(
	ctx context.Context,
	client *ethclient.Client,
	chain models.SysChain,
	contract models.SysContract,
	start, end uint64,
	events []TransferEvent,
	headers map[uint64]*blockHeaderMini,
	dbBlock *int64,
) error {

	if err := ix.applyChunkTx(
		ctx,
		client,
		chain.ChainID,
		contract.Address,
		start,
		end,
		events,
		headers,
	); err != nil {
		return err
	}

	// 非 OP Stack：apply 成功后立即推进 canonical block
	*dbBlock = int64(end)
	return nil
}

// computeScanRange 根据 cursor 状态计算本次需要扫描的区块范围
func (ix *Indexer) computeScanRange(
	cursor *models.BlockCursor,
	safeBlock uint64,
) (uint64, uint64) {

	startFrom := cursor.BlockNumber
	if cursor.ScanBlockNumber > startFrom {
		startFrom = cursor.ScanBlockNumber
	}

	return uint64(startFrom + 1), safeBlock
}

// computeChunkEnd 根据 chunkSize 计算当前扫描 chunk 的结束区块
func (ix *Indexer) computeChunkEnd(
	start uint64,
	applyTo uint64,
	chain models.SysChain,
) uint64 {

	// ChunkSize 在 DB model 里是 int，需转换
	end := start + uint64(chain.ChunkSize) - 1
	if end > applyTo {
		end = applyTo
	}
	return end
}

// updateScanProgress 更新内存中的扫描进度
func (ix *Indexer) updateScanProgress(
	chain models.SysChain,
	contract models.SysContract,
	end uint64,
) {
	key := fmt.Sprintf("%s:%d:%s",
		ix.cfg.Redis.KeyPrefix,
		chain.ChainID,
		contract.Address,
	)

	ix.scanMu.Lock()
	ix.scanCache[key] = int64(end)
	ix.scanMu.Unlock()
}

// flushScanCursor 将内存中的扫描进度按需写入数据库
func (ix *Indexer) flushScanCursor(
	ctx context.Context,
	client *ethclient.Client,
	chain models.SysChain,
	contract models.SysContract,
	scanFlushed *int64,
) error {
	key := fmt.Sprintf(
		"%s:%d:%s",
		ix.cfg.Redis.KeyPrefix,
		chain.ChainID,
		contract.Address,
	)

	ix.scanMu.Lock()
	scan := ix.scanCache[key]
	ix.scanMu.Unlock()

	newFlushed, err := ix.maybeFlushScanCursor(
		ctx,
		client,
		chain.ChainID,
		contract.Address,
		*scanFlushed,
		scan,
		chain.Type,
	)
	if err != nil {
		return err
	}

	*scanFlushed = newFlushed
	return nil
}

// ensureCursorHash 在初始化 cursor 缺失 block_hash 时补齐
func (ix *Indexer) ensureCursorHash(
	ctx context.Context,
	client *ethclient.Client,
	chain models.SysChain,
	cursor *models.BlockCursor,
) error {

	if cursor.BlockNumber <= 0 || cursor.BlockHash != "" {
		return nil
	}

	h, err := callRPCWithRetry(
		ctx,
		ix.rpcLimiter,
		"eth_getBlockByNumber",
		chain.ChainID,
		uint64(cursor.BlockNumber),
		func() (*types.Header, error) {
			return client.HeaderByNumber(ctx, big.NewInt(cursor.BlockNumber))
		},
	)
	if err != nil {
		return err
	}

	return ix.db.Model(&models.BlockCursor{}).
		Where("id=?", cursor.ID).
		Updates(map[string]any{
			"block_hash": h.Hash().Hex(),
			"updated_at": time.Now().UTC(),
		}).Error
}

/*
====================
Cursor
====================
*/

func (ix *Indexer) loadOrInitCursor(
	chainID int64,
	contract models.SysContract,
) (*models.BlockCursor, error) {

	var cursor models.BlockCursor
	err := ix.db.
		Where("chain_id=? AND contract_address=?", chainID, contract.Address).
		First(&cursor).Error

	if err == nil {
		return &cursor, nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		cursor = models.BlockCursor{
			ChainID:         chainID,
			ContractAddress: contract.Address,
			// DB 中 StartBlock 是 int64，可以直接赋值
			BlockNumber:   contract.StartBlock - 1,
			BlockHash:     "",
			LastBlockTime: time.Unix(0, 0).UTC(),
			UpdatedAt:     time.Now().UTC(),
		}
		if err := ix.db.Create(&cursor).Error; err != nil {
			return nil, err
		}
		return &cursor, nil
	}

	return nil, err
}

/*
====================
Transfer Fetch
====================
*/

type TransferEvent struct {
	BlockNumber uint64
	LogIndex    uint
	TxHash      common.Hash
	From        common.Address
	To          common.Address
	Value       *big.Int
	BlockTime   time.Time
}

type blockHeaderMini struct {
	Number uint64
	Hash   common.Hash
	Parent common.Hash
	Time   time.Time
}

func (ix *Indexer) fetchTransfers(
	ctx context.Context,
	client *ethclient.Client,
	token *erc20.TimeLedgerToken,
	chainID int64,
	contract string,
	start, end uint64,
) ([]TransferEvent, map[uint64]*blockHeaderMini, error) {

	// eth_getLogs（带 limiter + retry）
	it, err := callRPCWithRetry(
		ctx,
		ix.rpcLimiter,
		"eth_getLogs",
		chainID,
		start,
		func() (*erc20.TimeLedgerTokenTransferIterator, error) {
			return token.FilterTransfer(&bind.FilterOpts{
				Start: start,
				End:   &end,
			}, nil, nil)
		},
	)
	if err != nil {
		return nil, nil, err
	}

	headers := make(map[uint64]*blockHeaderMini)

	getHeader := func(bn uint64) (*blockHeaderMini, error) {
		if h, ok := headers[bn]; ok {
			return h, nil
		}

		h, err := callRPCWithRetry(
			ctx,
			ix.rpcLimiter,
			"eth_getBlockByNumber",
			chainID,
			bn,
			func() (*types.Header, error) {
				return client.HeaderByNumber(ctx, big.NewInt(int64(bn)))
			},
		)
		if err != nil {
			return nil, err
		}

		bh := &blockHeaderMini{
			Number: bn,
			Hash:   h.Hash(),
			Parent: h.ParentHash,
			Time:   time.Unix(int64(h.Time), 0).UTC(),
		}
		headers[bn] = bh
		return bh, nil
	}

	var events []TransferEvent

	for it.Next() {
		ev := it.Event
		if ev.Raw.Removed {
			continue
		}

		h, err := getHeader(ev.Raw.BlockNumber)
		if err != nil {
			return nil, nil, err
		}

		events = append(events, TransferEvent{
			BlockNumber: ev.Raw.BlockNumber,
			LogIndex:    ev.Raw.Index,
			TxHash:      ev.Raw.TxHash,
			From:        ev.From,
			To:          ev.To,
			Value:       ev.Value,
			BlockTime:   h.Time,
		})
	}

	sort.Slice(events, func(i, j int) bool {
		if events[i].BlockNumber != events[j].BlockNumber {
			return events[i].BlockNumber < events[j].BlockNumber
		}
		return events[i].LogIndex < events[j].LogIndex
	})

	return events, headers, nil
}

/*
====================
Apply Chunk
====================
*/

var zeroAddr = common.HexToAddress("0x0000000000000000000000000000000000000000")

func (ix *Indexer) applyChunkTx(
	ctx context.Context,
	client *ethclient.Client,
	chainID int64,
	contract string,
	start, end uint64,
	events []TransferEvent,
	headers map[uint64]*blockHeaderMini,
) error {

	return ix.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		//	写 block_header（幂等 + hash 一致性校验）
		for _, h := range headers {
			bh := models.BlockHeader{
				ChainID:         chainID,
				ContractAddress: contract,
				BlockNumber:     int64(h.Number),
				BlockHash:       h.Hash.Hex(),
				ParentHash:      h.Parent.Hex(),
				BlockTime:       h.Time,
				CreatedAt:       time.Now().UTC(),
			}

			res := tx.
				Clauses(clause.OnConflict{DoNothing: true}).
				Create(&bh)
			if res.Error != nil {
				return res.Error
			}

			// 已存在同高度 block，校验 hash
			if res.RowsAffected == 0 {
				var old models.BlockHeader
				if err := tx.
					Where(
						"chain_id=? AND contract_address=? AND block_number=?",
						chainID, contract, h.Number,
					).
					First(&old).Error; err != nil {
					return err
				}

				if old.BlockHash != h.Hash.Hex() {
					return fmt.Errorf(
						"block hash mismatch at block=%d old=%s new=%s",
						h.Number, old.BlockHash, h.Hash.Hex(),
					)
				}
			}
		}

		//	应用 Transfer 事件
		for _, ev := range events {
			if ev.From != zeroAddr {
				if err := ix.applyAccountDelta(
					tx, chainID, contract, ev, ev.From, new(big.Int).Neg(ev.Value),
				); err != nil {
					return err
				}
			}
			if ev.To != zeroAddr {
				if err := ix.applyAccountDelta(
					tx, chainID, contract, ev, ev.To, ev.Value,
				); err != nil {
					return err
				}
			}
		}

		//	chunk 末尾 block header 兜底（用于 cursor）
		endHeader := headers[end]
		if endHeader == nil {
			h, err := callRPCWithRetry(
				ctx,
				ix.rpcLimiter,
				"eth_getBlockByNumber",
				chainID,
				end,
				func() (*types.Header, error) {
					return client.HeaderByNumber(ctx, big.NewInt(int64(end)))
				},
			)
			if err != nil {
				return err
			}

			endHeader = &blockHeaderMini{
				Number: end,
				Hash:   h.Hash(),
				Parent: h.ParentHash,
				Time:   time.Unix(int64(h.Time), 0).UTC(),
			}
		}

		//	推进 cursor
		return tx.Model(&models.BlockCursor{}).
			Where("chain_id=? AND contract_address=?", chainID, contract).
			Updates(map[string]any{
				"block_number":    int64(end),
				"block_hash":      endHeader.Hash.Hex(),
				"last_block_time": endHeader.Time,
				"updated_at":      time.Now().UTC(),
			}).Error
	})
}

/*
====================
Balance
====================
*/

func (ix *Indexer) applyAccountDelta(
	tx *gorm.DB,
	chainID int64,
	contract string,
	ev TransferEvent,
	account common.Address,
	delta *big.Int,
) error {

	var ub models.UserBalance
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where(
			"chain_id=? AND contract_address=? AND account=?",
			chainID, contract, account.Hex(),
		).
		First(&ub).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		ub = models.UserBalance{
			ChainID:         chainID,
			ContractAddress: contract,
			Account:         account.Hex(),
			Balance:         "0",
		}
	} else if err != nil {
		return err
	}

	//	计算新余额
	cur, _ := new(big.Int).SetString(ub.Balance, 10)
	cur.Add(cur, delta)

	//	负数保护（必须）
	if cur.Sign() < 0 {
		return fmt.Errorf(
			"negative balance: acct=%s bal=%s delta=%s",
			account.Hex(), ub.Balance, delta.String(),
		)
	}

	//	写 balance_log（幂等）
	res := tx.
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&models.BalanceLog{
			ChainID:         chainID,
			ContractAddress: contract,
			Account:         account.Hex(),
			Delta:           delta.String(),
			BalanceAfter:    cur.String(),
			BlockNumber:     int64(ev.BlockNumber),
			BlockTime:       ev.BlockTime,
			TxHash:          ev.TxHash.Hex(),
			LogIndex:        int64(ev.LogIndex),
			CreatedAt:       time.Now().UTC(),
		})

	if res.Error != nil {
		return res.Error
	}

	//	如果这条 log 已存在（重复事件），直接退出
	if res.RowsAffected == 0 {
		return nil
	}

	//	仅在 log 真正插入成功后，更新 user_balance
	ub.Balance = cur.String()
	ub.BlockNumber = int64(ev.BlockNumber)
	ub.BlockTime = ev.BlockTime
	ub.UpdatedAt = time.Now().UTC()

	if ub.ID == 0 {
		return tx.Create(&ub).Error
	}

	return tx.Model(&models.UserBalance{}).
		Where("id=?", ub.ID).
		Updates(map[string]any{
			"balance":      ub.Balance,
			"block_number": ub.BlockNumber,
			"block_time":   ub.BlockTime,
			"updated_at":   ub.UpdatedAt,
		}).Error
}

func (ix *Indexer) maybeFlushScanCursor(
	ctx context.Context,
	client *ethclient.Client,
	chainID int64,
	contract string,
	lastFlushedScan int64,
	scanBlock int64,
	chainType string,
) (int64, error) {

	// 默认 gap
	scanFlushGap := int64(100)

	// Ethereum 可设为 0（立即更新）或保持 100
	if chainType == "ethereum" {
		scanFlushGap = 0
	}

	if scanBlock-lastFlushedScan < scanFlushGap {
		return lastFlushedScan, nil
	}

	// OP Stack：主动查时间
	var lastBlockTime time.Time

	if chainType == "opstack" {
		h, err := callRPCWithRetry(
			ctx,
			ix.rpcLimiter,
			"eth_getBlockByNumber",
			chainID,
			uint64(scanBlock),
			func() (*types.Header, error) {
				return client.HeaderByNumber(ctx, big.NewInt(scanBlock))
			},
		)
		if err != nil {
			return lastFlushedScan, err
		}
		lastBlockTime = time.Unix(int64(h.Time), 0).UTC()
	}

	updates := map[string]any{
		"scan_block_number": scanBlock,
		"updated_at":        time.Now().UTC(),
	}

	if !lastBlockTime.IsZero() {
		updates["last_block_time"] = lastBlockTime
	}

	if err := ix.db.Model(&models.BlockCursor{}).
		Where("chain_id=? AND contract_address=?", chainID, contract).
		Updates(updates).Error; err != nil {
		return lastFlushedScan, err
	}

	return scanBlock, nil
}
