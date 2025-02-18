package chain_scanner_service

import (
	"context"
	"fmt"

	"github.com/sonex-dex/sonex-chain-scanner/scanner/internal/svc"
	"github.com/sonex-dex/sonex-chain-scanner/scanner/module/chain_scanner"
	"github.com/sonex-dex/sonex-chain-scanner/scanner/module/global"
	"github.com/sonex-dex/sonex-chain-scanner/scanner/pb"
	"github.com/sonex-dex/sonex-chain-scanner/scanner/repository"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReplayLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReplayLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReplayLogic {
	return &ReplayLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ReplayLogic) Replay(in *pb.ChainScannerReplayReq) (*pb.EmptyRes, error) {
	conf := l.svcCtx.Config
	conn := repository.NewConn(conf.PostgreSQL.GetConnectString())
	scanPrograssModel := repository.NewScanPrograssModel(conn)
	record, err := scanPrograssModel.FindOneByChainId(l.ctx, conf.Chain.ChainId)
	switch err {
	case repository.ErrNotFound:
		return nil, fmt.Errorf("%s", "Miss scan prograss record")
	case nil:
	default:
		return nil, err
	}
	start := in.Start
	var end int64
	if start > record.Height {
		start = record.Height
	}
	if in.End == nil || *in.End > record.Height {
		end = record.Height
	} else {
		end = *in.End
	}
	global.Replay(chain_scanner.Replay{Start: start, End: end})

	return &pb.EmptyRes{}, nil
}
