package logger

import (
	"fmt"
	"github.com/saichler/l8types/go/ifs"
)

type FmtLogMethod struct {
}

func (fmtLogMethod *FmtLogMethod) Log(level ifs.LogLevel, msg string) {
	fmt.Println(msg)
}
