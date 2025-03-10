package sonex

import (
	"context"
	"testing"

	"github.com/sonex-dex/sonex-chain-scanner/pkg/eth_client"
	"github.com/sonex-dex/sonex-chain-scanner/scanner/config"
	"github.com/sonex-dex/sonex-chain-scanner/scanner/module/global"
	"github.com/zeromicro/go-zero/core/conf"
)

func TestUniswapV3Handler(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../../../test/scanner_local.yaml", &c)
	_ = global.InitIdGenerater(0)

	var blockNumber int64 = 6697634

	balancer, err := eth_client.NewEthClientBalancer(c.Chain)
	if err != nil {
		t.Error(err)
		return
	}
	client, err := balancer.GetClient(context.Background())
	if err != nil {
		t.Error(err)
		return
	}
	transactionMap, err := client.BatchGetTransactions(context.Background(), []int64{blockNumber}, 50)
	if err != nil {
		t.Error(err)
		return
	}

	handler := NewUniswapV3Handler(c, *c.Sonex)
	callbacks, err, _ := handler.Execute(context.Background(), client, blockNumber, blockNumber, transactionMap[blockNumber])
	if err != nil {
		t.Error(err)
		return
	}
	_ = callbacks
}
