package chain_scanner

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/sonex-dex/sonex-chain-scanner/pkg/eth_client"
	scanner_common "github.com/sonex-dex/sonex-chain-scanner/scanner/common"
	"github.com/sonex-dex/sonex-chain-scanner/scanner/internal/svc"
	"github.com/sonex-dex/sonex-chain-scanner/scanner/repository"

	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
)

type ChainScanner struct {
	svcCtx                   svc.ServiceContext
	ethClientBalancer        *eth_client.EthClientBalancer
	minConfirmations         int64
	minDelayWaitMilliseconds int64
	maxDelayWaitMilliseconds int64
	maxBlockNumber           int64
	handlers                 []IHandler
	waitForRetryHandlers     []IHandler
	callbacks                []func(ctx context.Context, client *eth_client.EthClient)
	cancel                   context.CancelFunc
	replay                   *Replay
}

func NewChainScanner(svcCtx svc.ServiceContext) (*ChainScanner, error) {
	conf := svcCtx.Config

	ethClientBalancer, err := eth_client.NewEthClientBalancer(conf.Chain)
	if err != nil {
		return nil, err
	}

	minConfirmations := conf.Chain.MinConfirmations
	if minConfirmations < 0 {
		minConfirmations = 5
	}
	minDelayWaitMilliseconds := conf.Chain.MinDelayWaitMilliseconds
	if minDelayWaitMilliseconds < 0 {
		minDelayWaitMilliseconds = 2000
	}
	maxDelayWaitMilliseconds := conf.Chain.MaxDelayWaitMilliseconds
	if maxDelayWaitMilliseconds < 0 {
		maxDelayWaitMilliseconds = 10000
	}
	maxBlockNumber := conf.Chain.MaxBlocksNumber
	if maxBlockNumber < 0 {
		maxBlockNumber = 1000
	}

	chainScanner := &ChainScanner{
		svcCtx:                   svcCtx,
		ethClientBalancer:        ethClientBalancer,
		minConfirmations:         minConfirmations,
		minDelayWaitMilliseconds: minDelayWaitMilliseconds,
		maxDelayWaitMilliseconds: maxDelayWaitMilliseconds,
		maxBlockNumber:           maxBlockNumber,
	}

	return chainScanner, nil
}

func (chainScanner *ChainScanner) scanBlock(ctx context.Context, client *eth_client.EthClient, start int64, end int64, handlers []IHandler, isReplay bool) ([]func(ctx context.Context, client *eth_client.EthClient), []IHandler, error) {
	prefix := "chain scanner"
	if isReplay {
		prefix = "chain replay scanner"
	}

	if start > end || len(handlers) == 0 {
		return nil, []IHandler{}, nil
	}

	startTime := time.Now()

	// Query transactions
	numbers := make([]int64, end-start)
	for i := start; i < end; i++ {
		numbers[i-start] = i
	}
	transactionMap, err := client.BatchGetTransactions(ctx, numbers, 50)
	if err != nil {
		return nil, []IHandler{}, err
	}
	transactions := []*eth_client.Transaction{}
	for blockNumber := start; blockNumber < end; blockNumber++ {
		txs, ok := transactionMap[blockNumber]
		if !ok {
			return nil, []IHandler{}, fmt.Errorf("%s: miss block[%d] transactions data", scanner_common.ERROR_EVENT_HANDLE, blockNumber)
		}
		transactions = append(transactions, txs...)
	}

	// Execute handlers
	var wg sync.WaitGroup
	results := make(chan struct {
		Callbacks []func(ctx context.Context, client *eth_client.EthClient)
		Err       error
		Skip      bool
		Handler   IHandler
		Cost      time.Duration
	}, len(handlers))

	for _, handler := range handlers {
		wg.Add(1)
		go func(
			handler IHandler,
			client *eth_client.EthClient,
			transactions []*eth_client.Transaction,
			ch chan struct {
				Callbacks []func(ctx context.Context, client *eth_client.EthClient)
				Err       error
				Skip      bool
				Handler   IHandler
				Cost      time.Duration
			},
		) {
			defer wg.Done()
			defer func() {
				if err := recover(); err != nil {
					ch <- struct {
						Callbacks []func(ctx context.Context, client *eth_client.EthClient)
						Err       error
						Skip      bool
						Handler   IHandler
						Cost      time.Duration
					}{
						Err:     fmt.Errorf("%s(panic)", err),
						Skip:    false,
						Handler: handler,
					}
				}
			}()

			startTime := time.Now()
			callbacks, err, skip := handler.Execute(ctx, client, start, end, transactions)
			if err != nil {
				callbacks = []func(ctx context.Context, client *eth_client.EthClient){}
			}
			ch <- struct {
				Callbacks []func(ctx context.Context, client *eth_client.EthClient)
				Err       error
				Skip      bool
				Handler   IHandler
				Cost      time.Duration
			}{
				Callbacks: callbacks,
				Err:       err,
				Skip:      skip,
				Handler:   handler,
				Cost:      time.Since(startTime),
			}
		}(handler, client, transactions, results)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	// Handle results
	waitForRetryHandlers := []IHandler{}
	callbacks := []func(ctx context.Context, client *eth_client.EthClient){}
	for result := range results {
		if result.Err != nil && !result.Skip {
			waitForRetryHandlers = append(chainScanner.waitForRetryHandlers, result.Handler)
			logx.Errorf("handler[%s], %s", reflect.TypeOf(result.Handler).String(), result.Err)
		}
		callbacks = append(callbacks, result.Callbacks...)
		logx.Infof("%s handler[%s] executed: cost %s", prefix, reflect.TypeOf(result.Handler).String(), result.Cost)
	}

	logx.Infof("%s: from %d to %d, tx count %d, retry count %d, cost %s", prefix, start, end, len(transactions), len(waitForRetryHandlers), time.Since(startTime))

	return callbacks, waitForRetryHandlers, nil
}

func (chainScanner *ChainScanner) AddHandler(handler IHandler) {
	chainScanner.handlers = append(chainScanner.handlers, handler)
}

func (chainScanner *ChainScanner) Start() {
	conf := chainScanner.svcCtx.Config
	ctx, cancel := context.WithCancel(context.Background())
	chainScanner.cancel = cancel

	// Main scan
	go func(ctx context.Context) {
		defer func() {
			logx.Errorf("%s", scanner_common.ERROR_EVENT_HANDLE_STOPPED)
		}()

		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			var handlers []IHandler
			var start int64
			var end int64

			client, err := chainScanner.ethClientBalancer.GetClient(ctx)
			if err != nil {
				logx.Errorf("%s: %s", scanner_common.ERROR_EVENT_HANDLE, err)
				continue
			}

			if len(chainScanner.waitForRetryHandlers) == 0 {
				handlers = chainScanner.handlers

				// Check start and end
				height, err := client.BlockNumber(ctx)
				if err != nil {
					logx.Errorf("%s: %s", scanner_common.ERROR_EVENT_HANDLE, err)
					continue
				}
				if start <= 0 {
					conn := repository.NewConn(conf.PostgreSQL.GetConnectString())
					scanPrograssModel := repository.NewScanPrograssModel(conn)
					prograss, err := scanPrograssModel.FindOneByChainId(ctx, conf.Chain.ChainId)
					switch err {
					case repository.ErrNotFound:
						prograss = &repository.ScanPrograss{Height: int64(height - 1)}
					case nil:
					default:
						logx.Errorf("%s: %s", scanner_common.ERROR_EVENT_HANDLE, err)
						continue
					}
					start = prograss.Height
				}

				err = retry.Do(func() error {
					height, err := client.BlockNumber(ctx)
					if err != nil {
						return err
					}
					end = int64(height) - chainScanner.minConfirmations
					if end > start {
						return nil
					}
					// end <= start, waiting block confirm.
					return fmt.Errorf("%s", scanner_common.ERROR_EVENT_HANDLE_WAITING_BLOCK)
				},
					retry.Attempts(0),
					retry.DelayType(retry.DelayTypeFunc(func(n uint, err error, config *retry.Config) time.Duration {
						delay := chainScanner.minDelayWaitMilliseconds * int64(n)
						if delay > chainScanner.maxDelayWaitMilliseconds {
							delay = chainScanner.maxDelayWaitMilliseconds
						}
						return time.Millisecond * time.Duration(delay)
					})),
					retry.OnRetry(func(n uint, err error) {
						logx.Infof("%s: %d times, %s", scanner_common.ERROR_EVENT_HANDLE_RETRY, n+1, err)
					}),
				)
				if err != nil {
					logx.Errorf("%s: %s", scanner_common.ERROR_EVENT_HANDLE_RETRY, err)
				}
				if end > start+chainScanner.maxBlockNumber {
					end = start + chainScanner.maxBlockNumber
				}
			} else {
				handlers = chainScanner.waitForRetryHandlers
				chainScanner.waitForRetryHandlers = []IHandler{}
			}

			callbacks, waitForRetryHandlers, err := chainScanner.scanBlock(ctx, client, start, end, handlers, false)
			if err != nil {
				logx.Errorf("%s: %s", scanner_common.ERROR_EVENT_HANDLE, err)
				continue
			}
			chainScanner.callbacks = append(chainScanner.callbacks, callbacks...)
			chainScanner.waitForRetryHandlers = waitForRetryHandlers

			if len(chainScanner.waitForRetryHandlers) == 0 {
				// Execute callback
				for _, callback := range chainScanner.callbacks {
					go callback(ctx, client)
				}
				chainScanner.callbacks = []func(ctx context.Context, client *eth_client.EthClient){}

				// Update scan prograss
				conn := repository.NewConn(conf.PostgreSQL.GetConnectString())
				scanPrograssModel := repository.NewScanPrograssModel(conn)
				record, err := scanPrograssModel.FindOneByChainId(ctx, conf.Chain.ChainId)
				switch err {
				case repository.ErrNotFound:
					scanPrograssModel.New(ctx, &repository.ScanPrograss{ChainId: conf.Chain.ChainId, Height: end + 1})
				case nil:
					err = scanPrograssModel.UpdateHeight(ctx, &repository.ScanPrograss{Id: record.Id, ChainId: conf.Chain.ChainId, Height: end + 1})
					if err != nil {
						logx.Errorf("%s: %s", scanner_common.ERROR_EVENT_HANDLE, err)
						continue
					}
				default:
					logx.Errorf("%s: %s", scanner_common.ERROR_EVENT_HANDLE, err)
					continue
				}
			}
		}
	}(ctx)

	// Replay scan
	go func(ctx context.Context) {
		defer func() {
			logx.Errorf("%s", scanner_common.ERROR_REPLAY_HANDLE_STOPPED)
		}()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				time.Sleep(1000)
			}

			if chainScanner.replay == nil {
				continue
			}

			// Check start/end
			start := chainScanner.replay.Start
			end := chainScanner.replay.End

			if start > end {
				chainScanner.replay = nil
				continue
			}
			if end > start+chainScanner.maxBlockNumber {
				end = start + chainScanner.maxBlockNumber
			}

			client, err := chainScanner.ethClientBalancer.GetClient(ctx)
			if err != nil {
				logx.Errorf("%s: %s", scanner_common.ERROR_REPLAY_HANDLE, err)
				continue
			}
			callbacks, waitForRetryHandlers, err := chainScanner.scanBlock(ctx, client, start, end, chainScanner.handlers, true)
			if err != nil {
				logx.Errorf("%s: %s", scanner_common.ERROR_REPLAY_HANDLE, err)
				continue
			}
			// Execute callback
			for _, callback := range callbacks {
				go callback(ctx, client)
			}

			if len(waitForRetryHandlers) > 0 {
				logx.Errorf("%s: %s", scanner_common.ERROR_REPLAY_HANDLE, "some handler executed failed")
				chainScanner.replay = nil
			} else {
				// Update scan prograss
				chainScanner.replay.Start = end + 1
			}
		}
	}(ctx)
}

func (chainScanner *ChainScanner) Stop() {
	chainScanner.cancel()
}

func (chainScanner *ChainScanner) Replay(request Replay) {
	chainScanner.replay = &request
}

func (chainScanner *ChainScanner) ReplayByBlockNumber(number int64) {
	chainScanner.Replay(Replay{
		Start: number,
		End:   number,
	})
}

func (chainScanner *ChainScanner) ReplayByTxHash(ctx context.Context, txHash string) error {
	client, err := chainScanner.ethClientBalancer.GetClient(ctx)
	if err != nil {
		return err
	}
	receipt, err := client.TransactionReceipt(ctx, common.HexToHash(txHash))
	if err != nil {
		return err
	}
	number := receipt.BlockNumber.Int64()
	chainScanner.Replay(Replay{
		Start: number,
		End:   number,
	})

	return nil
}
