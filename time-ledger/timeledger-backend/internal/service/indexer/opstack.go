package indexer

import (
	"context"
	"fmt"

	"github.com/Atom257/web3-labs/timeledger-backend/internal/config"
	"github.com/ethereum/go-ethereum/ethclient"
)

type OpStackAdapter struct{}

func (a *OpStackAdapter) SafeBlock(ctx context.Context, client *ethclient.Client, chain config.ChainConfig) (uint64, error) {
	head, err := client.BlockNumber(ctx)
	if err != nil {
		return 0, err
	}

	if chain.ReorgWindow <= 0 {
		return 0, fmt.Errorf("opstack requires reorg_window > 0")
	}

	w := uint64(chain.ReorgWindow)
	if head <= w {
		return 0, nil
	}
	return head - w, nil
}

func (a *OpStackAdapter) NeedBlockHeader() bool { return true }
