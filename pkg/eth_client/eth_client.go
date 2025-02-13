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
)

type EthClient struct {
	*ethclient.Client
}

func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	if number.Sign() >= 0 {
		return hexutil.EncodeBig(number)
	}
	// It's negative.
	if number.IsInt64() {
		return rpc.BlockNumber(number.Int64()).String()
	}
	// It's negative and large, which is invalid.
	return fmt.Sprintf("<invalid %d>", number)
}

func (ethClient *EthClient) BatchGetBlockByNumber(ctx context.Context, numbers []int64, full bool) (map[int64]*types.Block, error) {
	if len(numbers) == 0 {
		return map[int64]*types.Block{}, nil
	}

	batchCall := func(ctx context.Context, numbers []int64, full bool) (map[int64]*types.Block, error) {
		requests := make([]rpc.BatchElem, len(numbers))
		rawMessages := make([]*json.RawMessage, len(numbers))
		for i, number := range numbers {
			requests[i] = rpc.BatchElem{
				Method: "eth_getBlockByNumber",
				Args:   []any{toBlockNumArg(big.NewInt(number)), full},
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

	go func() {
		wg.Wait()
		close(results)
	}()

	// Handle results
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

func (ethClient *EthClient) BatchGetBlockReceiptsByNumber(ctx context.Context, numbers []int64, full bool) (map[int64][]*types.Receipt, error) {
	if len(numbers) == 0 {
		return map[int64][]*types.Receipt{}, nil
	}

	batchCall := func(ctx context.Context, numbers []int64) (map[int64][]*types.Receipt, error) {
		requests := make([]rpc.BatchElem, len(numbers))
		responses := make([][]*types.Receipt, len(numbers))
		for i, number := range numbers {
			requests[i] = rpc.BatchElem{
				Method: "eth_getBlockReceipts",
				Args:   []any{rpc.BlockNumberOrHashWithNumber(rpc.BlockNumber(number))},
				Result: &responses[i],
			}
		}
		err := ethClient.Client.Client().BatchCallContext(ctx, requests)
		if err != nil {
			return nil, err
		}
		result := make(map[int64][]*types.Receipt, len(numbers))
		for i, number := range numbers {
			if requests[i].Error != nil {
				return nil, requests[i].Error
			}
			response := responses[i]
			if response == nil {
				return nil, fmt.Errorf("get block[%d] receipts failed", number)
			}
			result[number] = response
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
		Receipts map[int64][]*types.Receipt
		Err      error
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
			Receipts map[int64][]*types.Receipt
			Err      error
		}) {
			defer wg.Done()
			defer func() {
				if err := recover(); err != nil {
					ch <- struct {
						Receipts map[int64][]*types.Receipt
						Err      error
					}{
						Err: fmt.Errorf("%s(panic)", err),
					}
				}
			}()

			receipts, err := batchCall(ctx, currentNumbers)
			if err != nil {
				ch <- struct {
					Receipts map[int64][]*types.Receipt
					Err      error
				}{
					Err: err,
				}
				return
			}

			ch <- struct {
				Receipts map[int64][]*types.Receipt
				Err      error
			}{
				Receipts: receipts,
			}
		}(ctx, currentNumbers, results)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	// Handle results
	receipts := map[int64][]*types.Receipt{}
	for result := range results {
		if result.Err != nil {
			return nil, result.Err
		}
		for key, value := range result.Receipts {
			receipts[key] = value
		}
	}

	return receipts, nil
}
