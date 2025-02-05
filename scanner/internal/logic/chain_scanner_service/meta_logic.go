package chain_scanner_service

import (
	"context"

	"github.com/sonex-dex/sonex-chain-scanner/scanner/internal/svc"
	"github.com/sonex-dex/sonex-chain-scanner/scanner/pb"
	"github.com/sonex-dex/sonex-chain-scanner/scanner/repository"

	"github.com/zeromicro/go-zero/core/logx"
)

type MetaLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMetaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MetaLogic {
	return &MetaLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MetaLogic) Meta(in *pb.EmptyReq) (*pb.ChainScannerMetaRes, error) {
	conf := l.svcCtx.Config
	conn := repository.NewConn(conf.PostgreSQL.GetConnectString())
	scanPrograssModel := repository.NewScanPrograssModel(conn)
	record, err := scanPrograssModel.FindOneByChainId(l.ctx, conf.Chain.ChainId)
	if err != nil {
		return nil, err
	}

	return &pb.ChainScannerMetaRes{
		ChainId:   record.ChainId,
		Height:    record.Height,
		UpdatedAt: record.UpdatedAt.Unix(),
	}, nil
}
