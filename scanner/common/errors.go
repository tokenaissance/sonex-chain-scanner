package common

import (
	"fmt"

	pkg_common "github.com/sonex-dex/sonex-chain-scanner/pkg/common"
)

var (
	ERROR_NO_AVAILABLE_NODE          pkg_common.Error = pkg_common.Error{Code: 20000, Raw: fmt.Errorf("no available node")}
	ERROR_EVENT_HANDLE               pkg_common.Error = pkg_common.Error{Code: 30000, Raw: fmt.Errorf("main scan error")}
	ERROR_EVENT_HANDLE_WAITING_BLOCK pkg_common.Error = pkg_common.Error{Code: 30001, Raw: fmt.Errorf("waiting block confirm")}
	ERROR_EVENT_HANDLE_STOPPED       pkg_common.Error = pkg_common.Error{Code: 30002, Raw: fmt.Errorf("main scan stopped")}
	ERROR_EVENT_HANDLE_RETRY         pkg_common.Error = pkg_common.Error{Code: 30002, Raw: fmt.Errorf("main scan retry")}
	ERROR_REPLAY_HANDLE              pkg_common.Error = pkg_common.Error{Code: 30010, Raw: fmt.Errorf("replay scan error")}
	ERROR_REPLAY_HANDLE_STOPPED      pkg_common.Error = pkg_common.Error{Code: 30011, Raw: fmt.Errorf("replay scan stopped")}
)
