package common

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	contract_common "github.com/sonex-dex/sonex-chain-scanner/pkg/contract/common"
	"github.com/sonex-dex/sonex-chain-scanner/pkg/eth_client"
	"github.com/sonex-dex/sonex-chain-scanner/scanner/config"
	"github.com/sonex-dex/sonex-chain-scanner/scanner/module/global"
	"github.com/sonex-dex/sonex-chain-scanner/scanner/repository"
)

type WETHHandler struct {
	conf config.Config
}

func NewWETHHandler(conf config.Config) *WETHHandler {
	return &WETHHandler{
		conf: conf,
	}
}

func (handler *WETHHandler) Execute(ctx context.Context, client *eth_client.EthClient, from int64, to int64, transactions []*eth_client.Transaction) (error, bool) {
	wethAddress := common.HexToAddress(handler.conf.Chain.ETH.WETH)

	weth, err := contract_common.NewWETH(common.Address{}, client)
	if err != nil {
		return err, false
	}

	conn := repository.NewConn(handler.conf.PostgreSQL.GetConnectString())
	commonEventWethDepositModel := repository.NewCommonEventWethDepositModel(conn)
	commonEventWethWithdrawalModel := repository.NewCommonEventWethWithdrawalModel(conn)

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

			// WETH deposit event
			{
				if event, err := weth.ParseDeposit(*log); err == nil {
					// If address not match, ignore
					if logAddress != wethAddress.Hex() {
						continue
					}

					eventDst := event.Dst.String()
					eventWad := event.Wad.String()

					// If transaction has been handle, ignore
					_, err := commonEventWethDepositModel.FindOneByChainIdTxHashLogIndex(ctx, chainId, txHash, logIndex)
					switch err {
					case repository.ErrNotFound:
					case nil:
						continue
					default:
						return err, false
					}

					id, err := global.NewId()
					if err != nil {
						return err, false
					}
					_, err = commonEventWethDepositModel.Insert(ctx, &repository.CommonEventWethDeposit{
						Id:          id,
						ChainId:     chainId,
						BlockNumber: blockNumber,
						TxHash:      txHash,
						TxAt:        txAt,
						TxFrom:      txFrom,
						TxTo:        txTo,
						TxValue:     txValue,
						LogIndex:    logIndex,
						LogAddress:  logAddress,
						EventDst:    eventDst,
						EventWad:    eventWad,
					})
					if err != nil {
						return err, false
					}

					continue
				}
			}

			// WETH withdrawal event
			{
				if event, err := weth.ParseWithdrawal(*log); err == nil {
					// If address not match, ignore
					if logAddress != wethAddress.Hex() {
						continue
					}

					eventSrc := event.Src.String()
					eventWad := event.Wad.String()

					// If transaction has been handle, ignore
					_, err := commonEventWethWithdrawalModel.FindOneByChainIdTxHashLogIndex(ctx, chainId, txHash, logIndex)
					switch err {
					case repository.ErrNotFound:
					case nil:
						continue
					default:
						return err, false
					}

					id, err := global.NewId()
					if err != nil {
						return err, false
					}
					_, err = commonEventWethWithdrawalModel.Insert(ctx, &repository.CommonEventWethWithdrawal{
						Id:          id,
						ChainId:     chainId,
						BlockNumber: blockNumber,
						TxHash:      txHash,
						TxAt:        txAt,
						TxFrom:      txFrom,
						TxTo:        txTo,
						TxValue:     txValue,
						LogIndex:    logIndex,
						LogAddress:  logAddress,
						EventSrc:    eventSrc,
						EventWad:    eventWad,
					})
					if err != nil {
						return err, false
					}

					continue
				}
			}
		}
	}

	return nil, false
}
