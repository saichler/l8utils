package logger

import (
	"fmt"
	"github.com/saichler/types/go/common"
)

type FmtLogMethod struct {
}

func (fmtLogMethod *FmtLogMethod) Log(level common.LogLevel, msg string) {
	fmt.Println(msg)
}
