package global

import (
	"fmt"

	"github.com/sonex-dex/sonex-chain-scanner/pkg/id_generater"
	"github.com/sonex-dex/sonex-chain-scanner/scanner/internal/svc"
	"github.com/sonex-dex/sonex-chain-scanner/scanner/module/chain_scanner"
)

var Ctx *svc.ServiceContext

var idGenerater *id_generater.IdGenerater
var chainScanner *chain_scanner.ChainScanner

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

func NewId() (int64, error) {
	if idGenerater == nil {
		return 0, fmt.Errorf("%s", "Id generater not initialized")
	}
	return idGenerater.Generate()
}
