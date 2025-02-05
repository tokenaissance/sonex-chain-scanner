package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	PostgreSQL PostgreSQLConfig `json:","`
	Chain      ChainConfig      `json:","`
}
