package chain_scanner

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/sonex-dex/sonex-chain-scanner/pkg/eth_client"
)

type IHandler interface {
	Execute(client *eth_client.EthClient, blocks []*types.Block) (error, bool)
}
