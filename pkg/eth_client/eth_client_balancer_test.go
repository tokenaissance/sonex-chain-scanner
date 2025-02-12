package eth_client

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

	balancer, err := NewEthClientBalancer(c.Chain)
	if err != nil {
		t.Error(err)
		return
	}
	client, err := balancer.GetClient(context.Background())
	if err != nil {
		t.Error(err)
		return
	}
	blocks, err := client.BatchGetBlockByNumber(context.Background(), []int64{8079828}, true)
	if err != nil {
		t.Error(err)
		return
	}
	for number, block := range blocks {
		for i, tx := range block.Transactions() {
			sender, _ := types.LatestSignerForChainID(big.NewInt(1946)).Sender(tx)
			fmt.Printf("%d[%d] ==> %s %s %s to %s\n", number, i, tx.Hash(), time.Unix(int64(block.Time()), 0).UTC(), sender, tx.To())
		}
	}
}

func TestBatchGetBlockReceiptsByNumber(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../test/scanner_local.yaml", &c)

	balancer, err := NewEthClientBalancer(c.Chain)
	if err != nil {
		t.Error(err)
		return
	}
	client, err := balancer.GetClient(context.Background())
	if err != nil {
		t.Error(err)
		return
	}
	receipts, err := client.BatchGetBlockReceiptsByNumber(context.Background(), []int64{8079828}, true)
	if err != nil {
		t.Error(err)
		return
	}
	for number, receiptArr := range receipts {
		for _, receipt := range receiptArr {
			fmt.Printf("%d[%d] ==> %s %+v\n", number, receipt.TransactionIndex, receipt.TxHash.Hex(), receipt.Logs)
		}
	}
}
