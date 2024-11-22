package logger

import (
	"fmt"
	"github.com/saichler/shared/go/share/interfaces"
)

type FmtLogMethod struct {
}

func (fmtLogMethod *FmtLogMethod) Log(level interfaces.LogLevel, msg string) {
	fmt.Println(msg)
}
