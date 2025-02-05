package chain_scanner

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
)

type IHandler interface {
	Execute(client *ethclient.Client, opt *bind.FilterOpts, isReplay bool) (error, bool)
}
