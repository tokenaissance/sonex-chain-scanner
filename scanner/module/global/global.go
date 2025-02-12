package global

import (
	"context"
	"fmt"

	"github.com/sonex-dex/sonex-chain-scanner/pkg/common"
	"github.com/sonex-dex/sonex-chain-scanner/pkg/eth_client"
	"github.com/sonex-dex/sonex-chain-scanner/pkg/id_generater"
	"github.com/sonex-dex/sonex-chain-scanner/scanner/internal/svc"
	"github.com/sonex-dex/sonex-chain-scanner/scanner/module/chain_scanner"
)

var Ctx *svc.ServiceContext

var idGenerater *id_generater.IdGenerater
var chainScanner *chain_scanner.ChainScanner
var ethClientBalancer *eth_client.EthClientBalancer

func InitIdGenerater(workerId int64) (err error) {
	idGenerater, err = id_generater.NewIdGenerater(workerId)
	if err != nil {
		return err
	}
	return nil
}

func InitChainScanner() (err error) {
	chainScanner, err = chain_scanner.NewChainScanner(*Ctx)
	if err != nil {
		return err
	}
	chainScanner.Start()
	return nil
}

func InitEthClientBalancer() (err error) {
	ethClientBalancer, err = eth_client.NewEthClientBalancer(Ctx.Config.Chain)
	if err != nil {
		return err
	}
	return nil
}

func NewId() (int64, error) {
	if idGenerater == nil {
		return 0, fmt.Errorf("%s, %s", common.ERROR_NOT_INITIALIZED, "id generater")
	}
	return idGenerater.Generate()
}

func AddChainHandler(handler chain_scanner.IHandler) error {
	if chainScanner == nil {
		return fmt.Errorf("%s, %s", common.ERROR_NOT_INITIALIZED, "chain scanner")
	}
	chainScanner.AddHandler(handler)
	return nil
}

func Replay(request chain_scanner.Replay) error {
	if chainScanner == nil {
		return fmt.Errorf("%s, %s", common.ERROR_NOT_INITIALIZED, "chain scanner")
	}
	chainScanner.Replay(request)
	return nil
}

func GetEthClient(ctx context.Context) (*eth_client.EthClient, error) {
	if ethClientBalancer == nil {
		return nil, fmt.Errorf("%s, %s", common.ERROR_NOT_INITIALIZED, "eth client balancer")
	}
	return ethClientBalancer.GetClient(ctx)
}
