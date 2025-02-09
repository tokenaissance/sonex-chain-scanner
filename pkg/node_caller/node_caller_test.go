package node_caller

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/sonex-dex/sonex-chain-scanner/scanner/config"
	"github.com/zeromicro/go-zero/core/conf"
)

func TestBatchGetBlockByNumber(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../test/scanner_local.yaml", &c)

	nodeCaller, err := NewNodeCallerWithUrl(*c.Chain.GetNode())
	if err != nil {
		t.Errorf("%s\n", err.Error())
		return
	}
	blocks, err := nodeCaller.BatchGetBlockByNumber(context.Background(), []int64{7945194}, true)
	if err != nil {
		t.Errorf("%s\n", err.Error())
		return
	}
	for number, block := range blocks {
		for i, tx := range block.Transactions() {
			sender, _ := types.LatestSignerForChainID(big.NewInt(1946)).Sender(tx)
			fmt.Printf("%d[%d] ==> %s %s %s to %s\n", number, i, tx.Hash(), time.Unix(int64(block.Time()), 0).UTC(), sender, tx.To())
		}
	}
}
