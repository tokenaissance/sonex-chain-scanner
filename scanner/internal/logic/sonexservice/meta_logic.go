package sonexservicelogic

import (
	"context"

	"github.com/sonex-dex/sonex-chain-scanner/scanner/internal/svc"
	"github.com/sonex-dex/sonex-chain-scanner/scanner/pb"

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

func (l *MetaLogic) Meta(in *pb.EmptyReq) (*pb.EmptyRes, error) {
	// todo: add your logic here and delete this line

	return &pb.EmptyRes{}, nil
}
