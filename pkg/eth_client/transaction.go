package eth_client

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Transaction struct {
	Block   *types.Block
	Tx      *types.Transaction
	Receipt *types.Receipt
}

func (transaction *Transaction) ChainId() int64 {
	return transaction.Tx.ChainId().Int64()
}

func (transaction *Transaction) BlockNumber() int64 {
	return transaction.Block.Number().Int64()
}

func (transaction *Transaction) Hash() common.Hash {
	return transaction.Tx.Hash()
}

func (transaction *Transaction) Time() time.Time {
	return time.Unix(int64(transaction.Block.Time()), 0)
}

func (transaction *Transaction) From() common.Address {
	from, _ := types.LatestSignerForChainID(transaction.Tx.ChainId()).Sender(transaction.Tx)
	return from
}

func (transaction *Transaction) To() *common.Address {
	return transaction.Tx.To()
}

func (transaction *Transaction) Value() *big.Int {
	return transaction.Tx.Value()
}
