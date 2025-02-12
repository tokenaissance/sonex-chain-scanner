package common

import "fmt"

type Error struct {
	Raw  error
	Code int64
}

func (error Error) String() string {
	return fmt.Sprintf("%s[%d]", error.Raw, error.Code)
}

var (
	ERROR_PANIC              Error = Error{Code: 10000, Raw: fmt.Errorf("panic")}
	ERROR_INVAILD_PARAMETERS Error = Error{Code: 10001, Raw: fmt.Errorf("invalid parameters")}
	ERROR_NOT_INITIALIZED    Error = Error{Code: 10002, Raw: fmt.Errorf("not initialized")}
	ERROR_MISS_CINFIGURE     Error = Error{Code: 10003, Raw: fmt.Errorf("miss configure")}
	ERROR_NO_AVAILABLE_NODE  Error = Error{Code: 10004, Raw: fmt.Errorf("no available node")}
)
