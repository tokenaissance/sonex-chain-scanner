package chain_scanner

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"

	pkg_common "github.com/sonex-dex/sonex-chain-scanner/pkg/common"
	"github.com/sonex-dex/sonex-chain-scanner/pkg/eth_client"
	"github.com/sonex-dex/sonex-chain-scanner/pkg/queue"
	"github.com/sonex-dex/sonex-chain-scanner/scanner/internal/svc"
	"github.com/sonex-dex/sonex-chain-scanner/scanner/repository"

	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
)

type ChainScanner struct {
	svcCtx                  svc.ServiceContext
	ethClientBalancer       *eth_client.EthClientBalancer
	minConfirmations        int64
	minDelayWaitMillisecond int64
	maxLogsNumber           int64
	handlers                []IHandler
	waitForRetryHandlers    []IHandler
	cancel                  context.CancelFunc
	replays                 queue.Queue
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
	minDelayWaitMillisecond := conf.Chain.MinDelayWaitMillisecond
	if minDelayWaitMillisecond < 0 {
		minDelayWaitMillisecond = 2000
	}
	maxLogsNumber := conf.Chain.MaxLogsNumber
	if maxLogsNumber < 0 {
		maxLogsNumber = 1000
	}

	chainScanner := &ChainScanner{
		svcCtx:                  svcCtx,
		ethClientBalancer:       ethClientBalancer,
		minConfirmations:        minConfirmations,
		minDelayWaitMillisecond: minDelayWaitMillisecond,
		maxLogsNumber:           maxLogsNumber,
	}

	return chainScanner, nil
}

func (chainScanner *ChainScanner) AddHandler(handler IHandler) {
	chainScanner.handlers = append(chainScanner.handlers, handler)
}

func (chainScanner *ChainScanner) Start() {
	conf := chainScanner.svcCtx.Config
	ctx, cancel := context.WithCancel(context.Background())
	chainScanner.cancel = cancel

	go func(ctx context.Context) {
		defer func() {
			logx.Errorf("chain scanner stopped")
		}()

		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			client, err := chainScanner.ethClientBalancer.GetClient(ctx)
			if err != nil {
				logx.Errorf("chain scanner error: %s\n", err)
				continue
			}
			opt := &bind.FilterOpts{}

			var startTime time.Time
			var transactions []*eth_client.Transaction
			err = func(ctx context.Context, client *eth_client.EthClient) (err error) {
				defer func() {
					if err := recover(); err != nil {
						err = fmt.Errorf("%s, %s", pkg_common.ERROR_PANIC, err)
					}
				}()

				var handlers []IHandler
				if len(chainScanner.waitForRetryHandlers) == 0 {
					handlers = chainScanner.handlers
					var start uint64
					var end *uint64

					// Check replay
					if replayA, ok := chainScanner.replays.Dequeue(); ok {
						replay, ok := replayA.(Replay)
						if ok {
							start = replay.Start
							end = replay.End
						}
					}

					// Check start and end
					height, err := client.BlockNumber(ctx)
					if err != nil {
						return err
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
							return err
						}
						start = uint64(prograss.Height)
					}
					if end == nil {
						newEnd := height - uint64(chainScanner.minConfirmations)
						end = &newEnd
					}

					err = retry.Do(func() error {
						if *end > start {
							return nil
						}
						newHeight, err := client.BlockNumber(ctx)
						if err != nil {
							return err
						}
						newEnd := newHeight - uint64(chainScanner.minConfirmations)
						end = &newEnd
						return fmt.Errorf("%s", "waiting block confirm...")
					},
						retry.Attempts(0),
						retry.DelayType(retry.DelayTypeFunc(func(n uint, err error, config *retry.Config) time.Duration {
							return time.Millisecond * time.Duration(chainScanner.minDelayWaitMillisecond*int64(n))
						})),
						retry.OnRetry(func(n uint, err error) {
							logx.Infof("chain scanner retry %d times: %s", n+1, err)
						}),
					)
					if err != nil {
						return err
					}
					if *end > start+uint64(chainScanner.maxLogsNumber) {
						*end = start + uint64(chainScanner.maxLogsNumber)
					}
					opt.Start = start
					opt.End = end
				} else {
					handlers = chainScanner.waitForRetryHandlers
					chainScanner.waitForRetryHandlers = []IHandler{}
				}

				startTime = time.Now()

				// Query transactions
				numbers := make([]int64, int(*opt.End-opt.Start))
				for i := opt.Start; i < *opt.End; i++ {
					numbers[i-opt.Start] = int64(i)
				}
				transactionMap, err := client.BatchGetTransactions(ctx, numbers, 50)
				if err != nil {
					return err
				}
				transactions = []*eth_client.Transaction{}
				for blockNumber := int64(opt.Start); blockNumber < int64(*opt.End); blockNumber++ {
					txs, ok := transactionMap[blockNumber]
					if !ok {
						return fmt.Errorf("miss block[%d] transactions data", blockNumber)
					}
					transactions = append(transactions, txs...)
				}

				// Execute handlers
				var wg sync.WaitGroup
				results := make(chan struct {
					Err     error
					Skip    bool
					Handler IHandler
					Cost    time.Duration
				}, len(handlers))

				for _, handler := range handlers {
					wg.Add(1)
					go func(
						handler IHandler,
						client *eth_client.EthClient,
						transactions []*eth_client.Transaction,
						ch chan struct {
							Err     error
							Skip    bool
							Handler IHandler
							Cost    time.Duration
						},
					) {
						defer wg.Done()
						defer func() {
							if err := recover(); err != nil {
								ch <- struct {
									Err     error
									Skip    bool
									Handler IHandler
									Cost    time.Duration
								}{
									Err:     fmt.Errorf("%s(panic)", err),
									Skip:    false,
									Handler: handler,
								}
							}
						}()

						startTime := time.Now()
						err, skip := handler.Execute(client, int64(opt.Start), int64(*opt.End), transactions)
						ch <- struct {
							Err     error
							Skip    bool
							Handler IHandler
							Cost    time.Duration
						}{
							Err:     err,
							Skip:    skip,
							Handler: handler,
							Cost:    time.Since(startTime),
						}
					}(handler, client, transactions, results)
				}

				go func() {
					wg.Wait()
					close(results)
				}()

				// Handle results
				for result := range results {
					if result.Err != nil && !result.Skip {
						chainScanner.waitForRetryHandlers = append(chainScanner.waitForRetryHandlers, result.Handler)
						logx.Errorf("chain scanner handler[%s] error: %s\n", reflect.TypeOf(result.Handler).String(), result.Err)
					}
					logx.Infof("chain scanner handler[%s] executed: cost %s \n", reflect.TypeOf(result.Handler).String(), result.Cost)
				}

				if len(chainScanner.waitForRetryHandlers) == 0 {
					conn := repository.NewConn(conf.PostgreSQL.GetConnectString())
					scanPrograssModel := repository.NewScanPrograssModel(conn)
					record, err := scanPrograssModel.FindOneByChainId(ctx, conf.Chain.ChainId)
					switch err {
					case repository.ErrNotFound:
						scanPrograssModel.New(ctx, &repository.ScanPrograss{ChainId: conf.Chain.ChainId, Height: int64(*opt.End) + 1})
					case nil:
						err = scanPrograssModel.UpdateHeight(ctx, &repository.ScanPrograss{Id: record.Id, ChainId: conf.Chain.ChainId, Height: int64(*opt.End) + 1})
						if err != nil {
							return err
						}
					default:
						return err
					}
				}

				return nil
			}(ctx, client)
			if err != nil {
				logx.Errorf("chain scanner error: %s\n", err)
				continue
			}
			logx.Infof("chain scanner: from %d to %d, tx count %d, cost %s", opt.Start, *opt.End, len(transactions), time.Since(startTime))
		}
	}(ctx)
}

func (chainScanner *ChainScanner) Stop() {
	chainScanner.cancel()
}

func (chainScanner *ChainScanner) Replay(request Replay) {
	chainScanner.replays.Enqueue(request)
}

func (chainScanner *ChainScanner) ReplayByBlockNumber(number uint64) {
	chainScanner.Replay(Replay{
		Start: number,
		End:   &number,
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
	number := receipt.BlockNumber.Uint64()
	chainScanner.Replay(Replay{
		Start: number,
		End:   &number,
	})

	return nil
}
