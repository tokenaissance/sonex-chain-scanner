package scanner

import (
	"fmt"

	"github.com/sonex-dex/sonex-chain-scanner/scanner/config"
	"github.com/sonex-dex/sonex-chain-scanner/scanner/internal/server/chain_scanner_service"
	"github.com/sonex-dex/sonex-chain-scanner/scanner/internal/svc"
	"github.com/sonex-dex/sonex-chain-scanner/scanner/module/event_handler/sonex"
	"github.com/sonex-dex/sonex-chain-scanner/scanner/module/global"
	"github.com/sonex-dex/sonex-chain-scanner/scanner/pb"

	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func Run(conf config.Config) {
	ctx := svc.NewServiceContext(conf)
	global.Ctx = ctx

	// Set global configure
	sqlx.DisableStmtLog()
	err := global.InitIdGenerater(0)
	if err != nil {
		panic(err)
	}
	err = global.InitChainScanner()
	if err != nil {
		panic(err)
	}
	err = global.InitEthClientBalancer()
	if err != nil {
		panic(err)
	}

	// Register chain handler
	sonexConf := conf.Sonex
	if sonexConf != nil {
		err = global.AddChainHandler(sonex.NewUniswapV3Handler(conf, *sonexConf))
		if err != nil {
			panic(err)
		}
	}

	// Rpc configure
	s := zrpc.MustNewServer(conf.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterChainScannerServiceServer(grpcServer, chain_scanner_service.NewChainScannerServiceServer(ctx))

		if conf.Mode == service.DevMode || conf.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", conf.ListenOn)
	s.Start()
}
