package main

import (
	"flag"

	"github.com/sonex-dex/sonex-chain-scanner/internal/config"
	"github.com/sonex-dex/sonex-chain-scanner/scanner"
	svcConfig "github.com/sonex-dex/sonex-chain-scanner/scanner/config"
	"github.com/zeromicro/go-zero/core/conf"
)

var configFile = flag.String("f", "etc/default.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	var svcConf svcConfig.Config
	err := c.GetSvcConfig(&svcConf)
	if err != nil {
		panic(err)
	}

	// Run service
	scanner.Run(svcConf)
}
