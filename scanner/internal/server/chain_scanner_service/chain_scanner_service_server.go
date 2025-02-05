package chain_scanner_service

import (
	"context"

	"github.com/sonex-dex/sonex-chain-scanner/scanner/internal/logic/chain_scanner_service"
	"github.com/sonex-dex/sonex-chain-scanner/scanner/internal/svc"
	"github.com/sonex-dex/sonex-chain-scanner/scanner/pb"
)

type ChainScannerServiceServer struct {
	svcCtx *svc.ServiceContext
	pb.UnimplementedChainScannerServiceServer
}

func NewChainScannerServiceServer(svcCtx *svc.ServiceContext) *ChainScannerServiceServer {
	return &ChainScannerServiceServer{
		svcCtx: svcCtx,
	}
}

func (s *ChainScannerServiceServer) Meta(ctx context.Context, in *pb.EmptyReq) (*pb.ChainScannerMetaRes, error) {
	l := chain_scanner_service.NewMetaLogic(ctx, s.svcCtx)
	return l.Meta(in)
}

func (s *ChainScannerServiceServer) Replay(ctx context.Context, in *pb.ChainScannerReplayReq) (*pb.EmptyRes, error) {
	l := chain_scanner_service.NewReplayLogic(ctx, s.svcCtx)
	return l.Replay(in)
}
