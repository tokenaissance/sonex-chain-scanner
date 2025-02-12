package eth_client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	pkg_common "github.com/sonex-dex/sonex-chain-scanner/pkg/common"
)

type EthClient struct {
	*ethclient.Client
}

func toBlockNumArg(number *big.Int) (string, error) {
	if number == nil {
		return "latest", nil
	}
	if number.Sign() >= 0 {
		return hexutil.EncodeBig(number), nil
	}
	if number.IsInt64() {
		return rpc.BlockNumber(number.Int64()).String(), nil
	}
	return "", fmt.Errorf("%s", pkg_common.ERROR_INVAILD_PARAMETERS)
}

func (ethClient *EthClient) BatchGetBlockByNumber(ctx context.Context, numbers []int64, full bool) (map[int64]*types.Block, error) {
	if len(numbers) == 0 {
		return map[int64]*types.Block{}, nil
	}

	batchCall := func(ctx context.Context, numbers []int64, full bool) (map[int64]*types.Block, error) {
		requests := make([]rpc.BatchElem, len(numbers))
		rawMessages := make([]*json.RawMessage, len(numbers))
		for i, number := range numbers {
			numberArg, err := toBlockNumArg(big.NewInt(number))
			if err != nil {
				return nil, fmt.Errorf("%s, number: %d", err, number)
			}
			requests[i] = rpc.BatchElem{
				Method: "eth_getBlockByNumber",
				Args:   []any{numberArg, full},
				Result: &rawMessages[i],
			}
		}
		err := ethClient.Client.Client().BatchCallContext(ctx, requests)
		if err != nil {
			return nil, err
		}
		result := make(map[int64]*types.Block, len(numbers))
		for i, number := range numbers {
			if requests[i].Error != nil {
				return nil, requests[i].Error
			}
			rawMessage := rawMessages[i]
			if rawMessage == nil {
				return nil, fmt.Errorf("get block[%d] raw data failed", number)
			}

			// Decode header and transactions.
			var head *types.Header
			err := json.Unmarshal(*rawMessages[i], &head)
			if err != nil {
				return nil, err
			}
			// When the block is not found, the API returns JSON null.
			if head == nil {
				return nil, ethereum.NotFound
			}

			var body rpcBlock
			if err := json.Unmarshal(*rawMessages[i], &body); err != nil {
				return nil, err
			}
			// Quick-verify transaction and uncle lists. This mostly helps with debugging the server.
			if head.UncleHash == types.EmptyUncleHash && len(body.UncleHashes) > 0 {
				return nil, errors.New("server returned non-empty uncle list but block header indicates no uncles")
			}
			if head.UncleHash != types.EmptyUncleHash && len(body.UncleHashes) == 0 {
				return nil, errors.New("server returned empty uncle list but block header indicates uncles")
			}
			if head.TxHash == types.EmptyTxsHash && len(body.Transactions) > 0 {
				return nil, errors.New("server returned non-empty transaction list but block header indicates no transactions")
			}
			if head.TxHash != types.EmptyTxsHash && len(body.Transactions) == 0 {
				return nil, errors.New("server returned empty transaction list but block header indicates transactions")
			}
			// Load uncles because they are not included in the block response.
			var uncles []*types.Header
			if len(body.UncleHashes) > 0 {
				uncles = make([]*types.Header, len(body.UncleHashes))
				reqs := make([]rpc.BatchElem, len(body.UncleHashes))
				for i := range reqs {
					reqs[i] = rpc.BatchElem{
						Method: "eth_getUncleByBlockHashAndIndex",
						Args:   []interface{}{body.Hash, hexutil.EncodeUint64(uint64(i))},
						Result: &uncles[i],
					}
				}
				err = ethClient.Client.Client().BatchCallContext(ctx, reqs)
				if err != nil {
					return nil, err
				}
				for i := range reqs {
					if reqs[i].Error != nil {
						return nil, reqs[i].Error
					}
					if uncles[i] == nil {
						return nil, fmt.Errorf("got null header for uncle %d of block %x", i, body.Hash[:])
					}
				}
			}
			// Fill the sender cache of transactions in the block.
			txs := make([]*types.Transaction, len(body.Transactions))
			for i, tx := range body.Transactions {
				if tx.From != nil {
					setSenderFromServer(tx.tx, *tx.From, body.Hash)
				}
				txs[i] = tx.tx
			}
			result[number] = types.NewBlockWithHeader(head).WithBody(
				types.Body{
					Transactions: txs,
					Uncles:       uncles,
					Withdrawals:  body.Withdrawals,
					Requests:     body.Requests,
				})
		}
		return result, nil
	}

	batchSize := 100
	batchCount := len(numbers) / batchSize
	if len(numbers)%batchSize > 0 {
		batchCount++
	}

	var wg sync.WaitGroup
	results := make(chan struct {
		Blocks map[int64]*types.Block
		Err    error
	}, batchCount)

	for i := 0; i < batchCount; i++ {
		start := i * batchSize
		end := (i + 1) * batchSize
		if end > len(numbers) {
			end = len(numbers)
		}
		currentNumbers := numbers[start:end]

		wg.Add(1)
		go func(ctx context.Context, currentNumbers []int64, ch chan struct {
			Blocks map[int64]*types.Block
			Err    error
		}) {
			defer wg.Done()
			defer func() {
				if err := recover(); err != nil {
					ch <- struct {
						Blocks map[int64]*types.Block
						Err    error
					}{
						Err: fmt.Errorf("%s(panic)", err),
					}
				}
			}()

			blocks, err := batchCall(ctx, currentNumbers, full)
			if err != nil {
				ch <- struct {
					Blocks map[int64]*types.Block
					Err    error
				}{
					Err: err,
				}
				return
			}

			ch <- struct {
				Blocks map[int64]*types.Block
				Err    error
			}{
				Blocks: blocks,
			}
		}(ctx, currentNumbers, results)
	}

	wg.Wait()
	close(results)

	blocks := map[int64]*types.Block{}
	for result := range results {
		if result.Err != nil {
			return nil, result.Err
		}
		for key, value := range result.Blocks {
			blocks[key] = value
		}
	}

	return blocks, nil
}
