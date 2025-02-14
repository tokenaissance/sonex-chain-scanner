package chain_scanner

import (
	"github.com/sonex-dex/sonex-chain-scanner/pkg/eth_client"
)

type IHandler interface {
	Execute(client *eth_client.EthClient, from int64, to int64, transactions []*eth_client.Transaction) (error, bool)
}
