package indexer

import (
	"context"
	"fmt"

	"github.com/Atom257/web3-labs/timeledger-backend/internal/models" // 引入 models
	"github.com/ethereum/go-ethereum/ethclient"
)

type ChainAdapter interface {
	// 【修改点】第三个参数改为 models.SysChain
	SafeBlock(ctx context.Context, client *ethclient.Client, chain models.SysChain) (uint64, error)
	NeedBlockHeader() bool
}

func AdapterFor(chainType string) (ChainAdapter, error) {
	switch chainType {
	case "ethereum":
		return &EthereumAdapter{}, nil
	case "opstack":
		return &OpStackAdapter{}, nil
	default:
		return nil, fmt.Errorf("unknown chain type: %s", chainType)
	}
}
