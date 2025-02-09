package common

import "fmt"

type Error struct {
	Code int64
	Raw  error
}

func (error Error) String() string {
	return fmt.Sprintf("%s[%d]", error.Raw, error.Code)
}

var (
	ERROR_INVAILD_PARAMETERS Error = Error{Code: 10000, Raw: fmt.Errorf("invalid parameters")}
)
