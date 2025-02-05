package global

import (
	"fmt"

	"github.com/sonex-dex/sonex-chain-scanner/pkg/id_generater"
	"github.com/sonex-dex/sonex-chain-scanner/scanner/internal/svc"
)

var Ctx *svc.ServiceContext
var idGenerater *id_generater.IdGenerater

func InitIdGenerater(workerId int64) error {
	generater, err := id_generater.NewIdGenerater(workerId)
	if err != nil {
		return err
	}
	idGenerater = generater
	return nil
}

func NewId() (int64, error) {
	if idGenerater == nil {
		return 0, fmt.Errorf("%s", "Id generater not initialized")
	}
	return idGenerater.Generate()
}
