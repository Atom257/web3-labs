package indexer

import (
	"context"

	"github.com/Atom257/web3-labs/timeledger-backend/internal/config"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EthereumAdapter struct{}

func (a *EthereumAdapter) SafeBlock(ctx context.Context, client *ethclient.Client, chain config.ChainConfig) (uint64, error) {
	head, err := client.BlockNumber(ctx)
	if err != nil {
		return 0, err
	}
	confirmations := uint64(chain.Confirmations)
	if head <= confirmations {
		return 0, nil
	}
	return head - confirmations, nil
}

func (a *EthereumAdapter) NeedBlockHeader() bool { return false }
