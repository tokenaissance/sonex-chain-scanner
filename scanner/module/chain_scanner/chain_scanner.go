package chain_scanner

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/sonex-dex/sonex-chain-scanner/pkg/queue"
	"github.com/sonex-dex/sonex-chain-scanner/scanner/internal/svc"
	"github.com/sonex-dex/sonex-chain-scanner/scanner/repository"

	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zeromicro/go-zero/core/logx"
)

type ChainScanner struct {
	svcCtx                  svc.ServiceContext
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
		minConfirmations:        minConfirmations,
		minDelayWaitMillisecond: minDelayWaitMillisecond,
		maxLogsNumber:           maxLogsNumber,
	}
	return chainScanner, nil
}

func (cs *ChainScanner) AddHandler(handler IHandler) {
	cs.handlers = append(cs.handlers, handler)
}

func (cs *ChainScanner) Start() {
	conf := cs.svcCtx.Config
	ctx, cancel := context.WithCancel(context.Background())
	cs.cancel = cancel

	go func(ctx context.Context) {
		defer func() {
			if err := recover(); err != nil {
				logx.Errorf("ChainScanner: ChainId[%d] stopped\nError: %s", conf.Chain.ChainId, err)
			}
		}()

		opt := &bind.FilterOpts{}
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			var startTime time.Time
			err := func(ctx context.Context) (result error) {
				defer func() {
					if err := recover(); err != nil {
						result = fmt.Errorf("%s(panic)", err)
					}
				}()

				client, err := conf.Chain.GetClient()
				if err != nil {
					return err
				}

				var handlers []IHandler
				isReplay := false
				if len(cs.waitForRetryHandlers) == 0 {
					handlers = cs.handlers
					var start uint64
					var end *uint64

					// Check replay
					if replayA, ok := cs.replays.Dequeue(); ok {
						replay, ok := replayA.(Replay)
						if ok {
							start = replay.Start
							end = replay.End
							if end != nil {
								isReplay = true
							}
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
						newEnd := height - uint64(cs.minConfirmations)
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
						newEnd := newHeight - uint64(cs.minConfirmations)
						end = &newEnd
						return fmt.Errorf("%s", "Waiting block confirm...")
					},
						// retry.RetryIf(func(err error) bool { return retrys < 5 }),
						retry.Delay(time.Millisecond*time.Duration(cs.minDelayWaitMillisecond)),
						retry.DelayType(retry.FixedDelay),
						retry.OnRetry(func(n uint, err error) {
							retrys += 1
							logx.Infof("ChainScanner retry %d times: %s", retrys, err)
						}),
					)
					if err != nil {
						return err
					}
					if *end > start+uint64(cs.maxLogsNumber) {
						*end = start + uint64(cs.maxLogsNumber)
					}
					opt.Start = start
					opt.End = end
				} else {
					handlers = cs.waitForRetryHandlers
					cs.waitForRetryHandlers = []IHandler{}
				}

				startTime = time.Now()

				var wg sync.WaitGroup
				results := make(chan struct {
					Err     error
					Skip    bool
					Handler IHandler
				}, len(handlers))

				for _, handler := range handlers {
					wg.Add(1)
					go func(
						handler IHandler,
						client *ethclient.Client,
						opt *bind.FilterOpts,
						isReplay bool,
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

						err, skip := handler.Execute(client, opt, isReplay)
						ch <- struct {
							Err     error
							Skip    bool
							Handler IHandler
						}{
							Err:     err,
							Skip:    skip,
							Handler: handler,
						}
					}(handler, client, opt, isReplay, results)
				}
				wg.Wait()
				close(results)

				// Handle results
				for result := range results {
					if result.Err != nil && !result.Skip {
						cs.waitForRetryHandlers = append(cs.waitForRetryHandlers, result.Handler)
						logx.Errorf("Handler[%s] error: %s\n", reflect.TypeOf(result.Handler).String(), result.Err)
					}
				}
				if !isReplay && len(cs.waitForRetryHandlers) == 0 {
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
			}(ctx)
			if err != nil {
				logx.Errorf("ChainScanner error: %s\n", err)
				continue
			}
			logx.Infof("ChainScanner: ChainId[%d] From %d to %d, cost %d ms", conf.Chain.ChainId, opt.Start, *opt.End, time.Since(startTime).Milliseconds())
		}
	}(ctx)
}

func (cs *ChainScanner) Stop() {
	cs.cancel()
}

func (cs *ChainScanner) Replay(request Replay) {
	cs.replays.Enqueue(request)
}

func (cs *ChainScanner) ReplayByBlockNumber(number uint64) {
	cs.Replay(Replay{
		Start: number,
		End:   &number,
	})
}

func (cs *ChainScanner) ReplayByTxHash(ctx context.Context, txHash string) error {
	conf := cs.svcCtx.Config
	client, err := conf.Chain.GetClient()
	if err != nil {
		return err
	}
	receipt, err := client.TransactionReceipt(ctx, common.HexToHash(txHash))
	if err != nil {
		return err
	}
	number := receipt.BlockNumber.Uint64()
	cs.Replay(Replay{
		Start: number,
		End:   &number,
	})
	return nil
}
