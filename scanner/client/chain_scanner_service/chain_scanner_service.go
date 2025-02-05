package chain_scanner_service

import (
	"context"

	"github.com/sonex-dex/sonex-chain-scanner/scanner/pb"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	EmptyReq              = pb.EmptyReq
	EmptyRes              = pb.EmptyRes
	ChainScannerMetaRes   = pb.ChainScannerMetaRes
	ChainScannerReplayReq = pb.ChainScannerReplayReq

	ChainScannerService interface {
		Meta(ctx context.Context, in *EmptyReq, opts ...grpc.CallOption) (*ChainScannerMetaRes, error)
		Replay(ctx context.Context, in *ChainScannerReplayReq, opts ...grpc.CallOption) (*EmptyRes, error)
	}

	defaultChainScannerService struct {
		cli zrpc.Client
	}
)

func NewChainScannerService(cli zrpc.Client) ChainScannerService {
	return &defaultChainScannerService{
		cli: cli,
	}
}

func (m *defaultChainScannerService) Meta(ctx context.Context, in *EmptyReq, opts ...grpc.CallOption) (*ChainScannerMetaRes, error) {
	client := pb.NewChainScannerServiceClient(m.cli.Conn())
	return client.Meta(ctx, in, opts...)
}

func (m *defaultChainScannerService) Replay(ctx context.Context, in *ChainScannerReplayReq, opts ...grpc.CallOption) (*EmptyRes, error) {
	client := pb.NewChainScannerServiceClient(m.cli.Conn())
	return client.Replay(ctx, in, opts...)
}
