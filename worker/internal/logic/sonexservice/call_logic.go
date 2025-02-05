package sonexservicelogic

import (
	"context"

	"github.com/sonex-dex/sonex-chain-scanner/worker/internal/svc"
	"github.com/sonex-dex/sonex-chain-scanner/worker/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CallLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCallLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CallLogic {
	return &CallLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CallLogic) Call(in *pb.SonexReq) (*pb.SonexRes, error) {
	// todo: add your logic here and delete this line

	return &pb.SonexRes{}, nil
}
