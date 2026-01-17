package indexer

import (
	"context"
	"fmt"

	"github.com/Atom257/web3-labs/timeledger-backend/internal/config"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ChainAdapter interface {
	SafeBlock(ctx context.Context, client *ethclient.Client, chain config.ChainConfig) (uint64, error)
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
