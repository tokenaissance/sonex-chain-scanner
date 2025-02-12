package common

import (
	"fmt"

	pkg_common "github.com/sonex-dex/sonex-chain-scanner/pkg/common"
)

var (
	ERROR_NO_AVAILABLE_NODE pkg_common.Error = pkg_common.Error{Code: 20000, Raw: fmt.Errorf("no available node")}
)
