package eth_client

import "github.com/ethereum/go-ethereum/core/types"

type Transaction struct {
	Block   *types.Block
	Receipt *types.Receipt
}
