package svc

import "github.com/sonex-dex/sonex-chain-scanner/worker/config"

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
