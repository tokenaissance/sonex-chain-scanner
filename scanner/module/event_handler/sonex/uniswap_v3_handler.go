package sonex

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/sonex-dex/sonex-chain-scanner/pkg/contract/sonex"
	"github.com/sonex-dex/sonex-chain-scanner/pkg/eth_client"
	"github.com/sonex-dex/sonex-chain-scanner/scanner/config"
	"github.com/sonex-dex/sonex-chain-scanner/scanner/module/global"
	"github.com/sonex-dex/sonex-chain-scanner/scanner/repository"
	"github.com/zeromicro/go-zero/core/logx"
)

type UniswapV3Handler struct {
	conf      config.Config
	sonexConf config.SonexConfig
}

func NewUniswapV3Handler(conf config.Config, sonexConf config.SonexConfig) *UniswapV3Handler {
	return &UniswapV3Handler{
		conf:      conf,
		sonexConf: sonexConf,
	}
}

func (handler *UniswapV3Handler) Execute(ctx context.Context, client *eth_client.EthClient, from int64, to int64, transactions []*eth_client.Transaction) ([]func(ctx context.Context, client *eth_client.EthClient), error, bool) {
	positionManagerAddress := common.HexToAddress(handler.sonexConf.PositionManager)

	factory, err := sonex.NewUniswapV3Factory(common.Address{}, client)
	if err != nil {
		return nil, err, false
	}
	pool, err := sonex.NewUniswapV3Pool(common.Address{}, client)
	if err != nil {
		return nil, err, false
	}
	manager, err := sonex.NewNonfungiblePositionManager(common.Address{}, client)
	if err != nil {
		return nil, err, false
	}

	conn := repository.NewConn(handler.conf.PostgreSQL.GetConnectString())
	sonexEventFactoryPoolCreatedModel := repository.NewSonexEventFactoryPoolCreatedModel(conn)
	sonexEventPoolMintModel := repository.NewSonexEventPoolMintModel(conn)
	sonexEventPoolBurnModel := repository.NewSonexEventPoolBurnModel(conn)
	sonexEventPoolCollectModel := repository.NewSonexEventPoolCollectModel(conn)
	sonexEventPoolSwapModel := repository.NewSonexEventPoolSwapModel(conn)
	sonexEventPositionManagerTransferModel := repository.NewSonexEventPositionManagerTransferModel(conn)
	sonexEventPositionManagerIncreaseLiquidityModel := repository.NewSonexEventPositionManagerIncreaseLiquidityModel(conn)
	sonexEventPositionManagerDecreaseLiquidityModel := repository.NewSonexEventPositionManagerDecreaseLiquidityModel(conn)
	sonexEventPositionManagerCollectModel := repository.NewSonexEventPositionManagerCollectModel(conn)

	callbacks := []func(ctx context.Context, client *eth_client.EthClient){}
	swapsInTransaction := map[string]*struct {
		Transaction *eth_client.Transaction
		Swaps       []*sonex.UniswapV3PoolSwap
	}{}

	for _, transaction := range transactions {
		for _, log := range transaction.Receipt.Logs {
			chainId := transaction.ChainId()
			blockNumber := transaction.BlockNumber()
			txHash := transaction.Hash().Hex()
			txAt := transaction.Time().UTC()
			txFrom := transaction.From().Hex()
			txTo := "0x"
			{
				to := transaction.To()
				if to != nil {
					txTo = to.Hex()
				}
			}
			txValue := transaction.Value().String()
			logIndex := int64(log.Index)
			logAddress := log.Address.Hex()

			// Create pool event
			{
				if event, err := factory.ParsePoolCreated(*log); err == nil {
					eventToken0 := event.Token0.Hex()
					eventToken1 := event.Token1.Hex()
					eventFee := event.Fee.Int64()
					eventTickSpacing := event.TickSpacing.Int64()
					eventPool := event.Pool.Hex()

					// If transaction has been handle, ignore
					_, err := sonexEventFactoryPoolCreatedModel.FindOneByChainIdTxHashLogIndex(ctx, chainId, txHash, logIndex)
					switch err {
					case repository.ErrNotFound:
					case nil:
						continue
					default:
						return nil, err, false
					}

					id, err := global.NewId()
					if err != nil {
						return nil, err, false
					}
					_, err = sonexEventFactoryPoolCreatedModel.Insert(ctx, &repository.SonexEventFactoryPoolCreated{
						Id:               id,
						ChainId:          chainId,
						BlockNumber:      blockNumber,
						TxHash:           txHash,
						TxAt:             txAt,
						TxFrom:           txFrom,
						TxTo:             txTo,
						TxValue:          txValue,
						LogIndex:         logIndex,
						LogAddress:       logAddress,
						EventToken0:      eventToken0,
						EventToken1:      eventToken1,
						EventFee:         eventFee,
						EventTickSpacing: eventTickSpacing,
						EventPool:        eventPool,
					})
					if err != nil {
						return nil, err, false
					}

					continue
				}
			}

			// Mint position event
			{
				if event, err := pool.ParseMint(*log); err == nil {
					eventSender := event.Sender.Hex()
					eventOwner := event.Owner.Hex()
					eventTickLower := event.TickLower.Int64()
					eventTickUpper := event.TickUpper.Int64()
					eventAmount := event.Amount.String()
					eventAmount0 := event.Amount0.String()
					eventAmount1 := event.Amount1.String()

					// If owner address not match, ignore
					if eventOwner != positionManagerAddress.Hex() {
						continue
					}

					// If transaction has been handle, ignore
					_, err := sonexEventPoolMintModel.FindOneByChainIdTxHashLogIndex(ctx, chainId, txHash, logIndex)
					switch err {
					case repository.ErrNotFound:
					case nil:
						continue
					default:
						return nil, err, false
					}

					id, err := global.NewId()
					if err != nil {
						return nil, err, false
					}
					_, err = sonexEventPoolMintModel.Insert(ctx, &repository.SonexEventPoolMint{
						Id:             id,
						ChainId:        chainId,
						BlockNumber:    blockNumber,
						TxHash:         txHash,
						TxAt:           txAt,
						TxFrom:         txFrom,
						TxTo:           txTo,
						TxValue:        txValue,
						LogIndex:       logIndex,
						LogAddress:     logAddress,
						EventSender:    eventSender,
						EventOwner:     eventOwner,
						EventTickLower: eventTickLower,
						EventTickUpper: eventTickUpper,
						EventAmount:    eventAmount,
						EventAmount0:   eventAmount0,
						EventAmount1:   eventAmount1,
					})
					if err != nil {
						return nil, err, false
					}

					continue
				}
			}

			// Burn position event
			{
				if event, err := pool.ParseBurn(*log); err == nil {
					eventOwner := event.Owner.Hex()
					eventTickLower := event.TickLower.Int64()
					eventTickUpper := event.TickUpper.Int64()
					eventAmount := event.Amount.String()
					eventAmount0 := event.Amount0.String()
					eventAmount1 := event.Amount1.String()

					// If owner address not match, ignore
					if eventOwner != positionManagerAddress.Hex() {
						continue
					}

					// If transaction has been handle, ignore
					_, err := sonexEventPoolBurnModel.FindOneByChainIdTxHashLogIndex(ctx, chainId, txHash, logIndex)
					switch err {
					case repository.ErrNotFound:
					case nil:
						continue
					default:
						return nil, err, false
					}

					id, err := global.NewId()
					if err != nil {
						return nil, err, false
					}
					_, err = sonexEventPoolBurnModel.Insert(ctx, &repository.SonexEventPoolBurn{
						Id:             id,
						ChainId:        chainId,
						BlockNumber:    blockNumber,
						TxHash:         txHash,
						TxAt:           txAt,
						TxFrom:         txFrom,
						TxTo:           txTo,
						TxValue:        txValue,
						LogIndex:       logIndex,
						LogAddress:     logAddress,
						EventOwner:     eventOwner,
						EventTickLower: eventTickLower,
						EventTickUpper: eventTickUpper,
						EventAmount:    eventAmount,
						EventAmount0:   eventAmount0,
						EventAmount1:   eventAmount1,
					})
					if err != nil {
						return nil, err, false
					}

					continue
				}
			}

			// Pool collect event
			{
				if event, err := pool.ParseCollect(*log); err == nil {
					eventOwner := event.Owner.Hex()
					eventRecipient := event.Recipient.Hex()
					eventTickLower := event.TickLower.Int64()
					eventTickUpper := event.TickUpper.Int64()
					eventAmount0 := event.Amount0.String()
					eventAmount1 := event.Amount1.String()

					// If owner address not match, ignore
					if eventOwner != positionManagerAddress.Hex() {
						continue
					}

					// If transaction has been handle, ignore
					_, err := sonexEventPoolCollectModel.FindOneByChainIdTxHashLogIndex(ctx, chainId, txHash, logIndex)
					switch err {
					case repository.ErrNotFound:
					case nil:
						continue
					default:
						return nil, err, false
					}

					id, err := global.NewId()
					if err != nil {
						return nil, err, false
					}
					_, err = sonexEventPoolCollectModel.Insert(ctx, &repository.SonexEventPoolCollect{
						Id:             id,
						ChainId:        chainId,
						BlockNumber:    blockNumber,
						TxHash:         txHash,
						TxAt:           txAt,
						TxFrom:         txFrom,
						TxTo:           txTo,
						TxValue:        txValue,
						LogIndex:       logIndex,
						LogAddress:     logAddress,
						EventOwner:     eventOwner,
						EventRecipient: eventRecipient,
						EventTickLower: eventTickLower,
						EventTickUpper: eventTickUpper,
						EventAmount0:   eventAmount0,
						EventAmount1:   eventAmount1,
					})
					if err != nil {
						return nil, err, false
					}

					continue
				}
			}

			// Swap event
			{
				if event, err := pool.ParseSwap(*log); err == nil {
					eventSender := event.Sender.Hex()
					eventRecipient := event.Recipient.Hex()
					eventAmount0 := event.Amount0.String()
					eventAmount1 := event.Amount1.String()
					eventSqrtPriceX96 := event.SqrtPriceX96.String()
					eventLiquidity := event.Liquidity.String()
					eventTick := event.Tick.Int64()

					// If transaction has been handle, ignore
					_, err := sonexEventPoolSwapModel.FindOneByChainIdTxHashLogIndex(ctx, chainId, txHash, logIndex)
					switch err {
					case repository.ErrNotFound:
					case nil:
						continue
					default:
						return nil, err, false
					}

					id, err := global.NewId()
					if err != nil {
						return nil, err, false
					}
					_, err = sonexEventPoolSwapModel.Insert(ctx, &repository.SonexEventPoolSwap{
						Id:                id,
						ChainId:           chainId,
						BlockNumber:       blockNumber,
						TxHash:            txHash,
						TxAt:              txAt,
						TxFrom:            txFrom,
						TxTo:              txTo,
						TxValue:           txValue,
						LogIndex:          logIndex,
						LogAddress:        logAddress,
						EventSender:       eventSender,
						EventRecipient:    eventRecipient,
						EventAmount0:      eventAmount0,
						EventAmount1:      eventAmount1,
						EventSqrtPriceX96: eventSqrtPriceX96,
						EventLiquidity:    eventLiquidity,
						EventTick:         eventTick,
					})
					if err != nil {
						return nil, err, false
					}

					// Callback
					if _, ok := swapsInTransaction[txHash]; !ok {
						swapsInTransaction[txHash] = &struct {
							Transaction *eth_client.Transaction
							Swaps       []*sonex.UniswapV3PoolSwap
						}{
							Transaction: transaction,
						}
					}
					swapsInTransaction[txHash].Swaps = append(swapsInTransaction[txHash].Swaps, event)

					continue
				}
			}

			// Transfer position event
			{
				if event, err := manager.ParseTransfer(*log); err == nil {
					// If address not match, ignore
					if logAddress != positionManagerAddress.Hex() {
						continue
					}

					eventFrom := event.From.Hex()
					eventTo := event.To.Hex()
					eventTokenId := event.TokenId.String()

					// If transaction has been handle, ignore
					_, err := sonexEventPositionManagerTransferModel.FindOneByChainIdTxHashLogIndex(ctx, chainId, txHash, logIndex)
					switch err {
					case repository.ErrNotFound:
					case nil:
						continue
					default:
						return nil, err, false
					}

					id, err := global.NewId()
					if err != nil {
						return nil, err, false
					}
					_, err = sonexEventPositionManagerTransferModel.Insert(ctx, &repository.SonexEventPositionManagerTransfer{
						Id:           id,
						ChainId:      chainId,
						BlockNumber:  blockNumber,
						TxHash:       txHash,
						TxAt:         txAt,
						TxFrom:       txFrom,
						TxTo:         txTo,
						TxValue:      txValue,
						LogIndex:     logIndex,
						LogAddress:   logAddress,
						EventFrom:    eventFrom,
						EventTo:      eventTo,
						EventTokenId: eventTokenId,
					})
					if err != nil {
						return nil, err, false
					}

					continue
				}
			}

			// Increase liquidity event
			{
				if event, err := manager.ParseIncreaseLiquidity(*log); err == nil {
					// If address not match, ignore
					if logAddress != positionManagerAddress.Hex() {
						continue
					}

					eventTokenId := event.TokenId.String()
					eventLiquidity := event.Liquidity.String()
					eventAmount0 := event.Amount0.String()
					eventAmount1 := event.Amount1.String()

					// If transaction has been handle, ignore
					_, err := sonexEventPositionManagerIncreaseLiquidityModel.FindOneByChainIdTxHashLogIndex(ctx, chainId, txHash, logIndex)
					switch err {
					case repository.ErrNotFound:
					case nil:
						continue
					default:
						return nil, err, false
					}

					id, err := global.NewId()
					if err != nil {
						return nil, err, false
					}
					_, err = sonexEventPositionManagerIncreaseLiquidityModel.Insert(ctx, &repository.SonexEventPositionManagerIncreaseLiquidity{
						Id:             id,
						ChainId:        chainId,
						BlockNumber:    blockNumber,
						TxHash:         txHash,
						TxAt:           txAt,
						TxFrom:         txFrom,
						TxTo:           txTo,
						TxValue:        txValue,
						LogIndex:       logIndex,
						LogAddress:     logAddress,
						EventTokenId:   eventTokenId,
						EventLiquidity: eventLiquidity,
						EventAmount0:   eventAmount0,
						EventAmount1:   eventAmount1,
					})
					if err != nil {
						return nil, err, false
					}

					continue
				}
			}

			// Decrease liquidity event
			{
				if event, err := manager.ParseDecreaseLiquidity(*log); err == nil {
					// If address not match, ignore
					if logAddress != positionManagerAddress.Hex() {
						continue
					}

					eventTokenId := event.TokenId.String()
					eventLiquidity := event.Liquidity.String()
					eventAmount0 := event.Amount0.String()
					eventAmount1 := event.Amount1.String()

					// If transaction has been handle, ignore
					_, err := sonexEventPositionManagerDecreaseLiquidityModel.FindOneByChainIdTxHashLogIndex(ctx, chainId, txHash, logIndex)
					switch err {
					case repository.ErrNotFound:
					case nil:
						continue
					default:
						return nil, err, false
					}

					id, err := global.NewId()
					if err != nil {
						return nil, err, false
					}
					_, err = sonexEventPositionManagerDecreaseLiquidityModel.Insert(ctx, &repository.SonexEventPositionManagerDecreaseLiquidity{
						Id:             id,
						ChainId:        chainId,
						BlockNumber:    blockNumber,
						TxHash:         txHash,
						TxAt:           txAt,
						TxFrom:         txFrom,
						TxTo:           txTo,
						TxValue:        txValue,
						LogIndex:       logIndex,
						LogAddress:     logAddress,
						EventTokenId:   eventTokenId,
						EventLiquidity: eventLiquidity,
						EventAmount0:   eventAmount0,
						EventAmount1:   eventAmount1,
					})
					if err != nil {
						return nil, err, false
					}

					continue
				}
			}

			// Manager collect event
			{
				if event, err := manager.ParseCollect(*log); err == nil {
					// If address not match, ignore
					if logAddress != positionManagerAddress.Hex() {
						continue
					}

					eventTokenId := event.TokenId.String()
					eventRecipient := event.Recipient.String()
					eventAmount0 := event.Amount0.String()
					eventAmount1 := event.Amount1.String()

					// If transaction has been handle, ignore
					_, err := sonexEventPositionManagerCollectModel.FindOneByChainIdTxHashLogIndex(ctx, chainId, txHash, logIndex)
					switch err {
					case repository.ErrNotFound:
					case nil:
						continue
					default:
						return nil, err, false
					}

					id, err := global.NewId()
					if err != nil {
						return nil, err, false
					}
					_, err = sonexEventPositionManagerCollectModel.Insert(ctx, &repository.SonexEventPositionManagerCollect{
						Id:             id,
						ChainId:        chainId,
						BlockNumber:    blockNumber,
						TxHash:         txHash,
						TxAt:           txAt,
						TxFrom:         txFrom,
						TxTo:           txTo,
						TxValue:        txValue,
						LogIndex:       logIndex,
						LogAddress:     logAddress,
						EventTokenId:   eventTokenId,
						EventRecipient: eventRecipient,
						EventAmount0:   eventAmount0,
						EventAmount1:   eventAmount1,
					})
					if err != nil {
						return nil, err, false
					}

					continue
				}
			}
		}
	}

	// Swap callback handle
	for k, v := range swapsInTransaction {
		callbacks = append(callbacks, func(ctx context.Context, client *eth_client.EthClient) {
			logx.Infof("Tx: %s, %d swaps, ", k, len(v.Swaps))
		})
	}

	return callbacks, nil, false
}
