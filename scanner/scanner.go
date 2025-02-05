package scanner

import (
	"fmt"

	"github.com/sonex-dex/sonex-chain-scanner/scanner/config"
	sonexserviceServer "github.com/sonex-dex/sonex-chain-scanner/scanner/internal/server/sonexservice"
	"github.com/sonex-dex/sonex-chain-scanner/scanner/internal/svc"
	"github.com/sonex-dex/sonex-chain-scanner/scanner/pb"

	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func Run(conf config.Config) {
	ctx := svc.NewServiceContext(conf)

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
