package scanner

import (
	"fmt"

	"github.com/sonex-dex/sonex-chain-scanner/scanner/config"
	sonexserviceServer "github.com/sonex-dex/sonex-chain-scanner/scanner/internal/server/sonexservice"
	"github.com/sonex-dex/sonex-chain-scanner/scanner/internal/svc"
	"github.com/sonex-dex/sonex-chain-scanner/scanner/module/chain_scanner"
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

	// Start scan chain
	chainScanner, err := chain_scanner.NewChainScanner()
	if err != nil {
		panic(err)
	}
	chainScanner.Start()

	// Rpc configure
	s := zrpc.MustNewServer(conf.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterSonexServiceServer(grpcServer, sonexserviceServer.NewSonexServiceServer(ctx))

		if conf.Mode == service.DevMode || conf.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", conf.ListenOn)
	s.Start()
}
