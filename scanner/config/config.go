package config

import (
	"github.com/sonex-dex/sonex-chain-scanner/pkg/common/config"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	PostgreSQL PostgreSQLConfig   `json:","`
	Chain      config.ChainConfig `json:","`
}
