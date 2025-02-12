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
	"github.com/ethereum/go-ethereum/core/types"
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
			err = func(ctx context.Context, client *eth_client.EthClient) (result error) {
				defer func() {
					if err := recover(); err != nil {
						result = fmt.Errorf("%s, %s", pkg_common.ERROR_PANIC, err)
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

					retrys := 0
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
						retry.Delay(time.Millisecond*time.Duration(chainScanner.minDelayWaitMillisecond)),
						retry.DelayType(retry.FixedDelay),
						retry.OnRetry(func(n uint, err error) {
							retrys += 1
							logx.Infof("chain scanner retry %d times: %s", retrys, err)
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

				var wg sync.WaitGroup

				numbers := make([]int64, int(*opt.End-opt.Start))
				queryResults := make(chan struct {
					Err error
				}, len(numbers))

				// Query blocks
				var blocks []*types.Block
				wg.Add(1)
				go func(ch chan struct {
					Err error
				}) {
					defer wg.Done()
					defer func() {
						if err := recover(); err != nil {
							ch <- struct {
								Err error
							}{
								Err: fmt.Errorf("%s(panic)", err),
							}
						}
					}()

					for i := opt.Start; i < *opt.End; i++ {
						numbers[i-opt.Start] = int64(i)
					}
					blockMap, err := client.BatchGetBlockByNumber(ctx, numbers, true)
					if err != nil {
						ch <- struct {
							Err error
						}{
							Err: err,
						}
						return
					}
					blocks = make([]*types.Block, len(blockMap))
					for i := opt.Start; i < *opt.End; i++ {
						blocks[i-opt.Start] = blockMap[int64(i)]
					}
				}(queryResults)

				// Query block receipts
				var receipts map[int64][]*types.Receipt
				wg.Add(1)
				go func(ch chan struct {
					Err error
				}) {
					defer wg.Done()
					defer func() {
						if err := recover(); err != nil {
							ch <- struct {
								Err error
							}{
								Err: fmt.Errorf("%s(panic)", err),
							}
						}
					}()

					for i := opt.Start; i < *opt.End; i++ {
						numbers[i-opt.Start] = int64(i)
					}
					receiptsMap, err := client.BatchGetBlockReceiptsByNumber(ctx, numbers, true)
					if err != nil {
						ch <- struct {
							Err error
						}{
							Err: err,
						}
						return
					}
					receipts = receiptsMap
				}(queryResults)

				wg.Wait()
				close(queryResults)

				// Handle results
				for result := range queryResults {
					if result.Err != nil {
						return err
					}
				}

				// Execute handlers
				results := make(chan struct {
					Err     error
					Skip    bool
					Handler IHandler
				}, len(handlers))

				for _, handler := range handlers {
					wg.Add(1)
					go func(
						handler IHandler,
						client *eth_client.EthClient,
						blocks []*types.Block,
						ch chan struct {
							Err     error
							Skip    bool
							Handler IHandler
						},
					) {
						defer wg.Done()
						defer func() {
							if err := recover(); err != nil {
								ch <- struct {
									Err     error
									Skip    bool
									Handler IHandler
								}{
									Err:     fmt.Errorf("%s(panic)", err),
									Skip:    false,
									Handler: handler,
								}
							}
						}()

						err, skip := handler.Execute(client, blocks, receipts)
						ch <- struct {
							Err     error
							Skip    bool
							Handler IHandler
						}{
							Err:     err,
							Skip:    skip,
							Handler: handler,
						}
					}(handler, client, blocks, results)
				}

				wg.Wait()
				close(results)

				// Handle results
				for result := range results {
					if result.Err != nil && !result.Skip {
						chainScanner.waitForRetryHandlers = append(chainScanner.waitForRetryHandlers, result.Handler)
						logx.Errorf("chain scanner handler[%s] error: %s\n", reflect.TypeOf(result.Handler).String(), result.Err)
					}
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
			logx.Infof("chain scanner: from %d to %d, cost %d ms", opt.Start, *opt.End, time.Since(startTime).Milliseconds())
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
