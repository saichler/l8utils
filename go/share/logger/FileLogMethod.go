package logger

import (
	"github.com/saichler/types/go/common"
	"os"
)

type FileLogMethod struct {
	filename string
	file     *os.File
}

func NewFileLogMethod(filename string) *FileLogMethod {
	return &FileLogMethod{filename: filename}
}

func (this *FileLogMethod) Log(level common.LogLevel, msg string) {
	if this.file == nil {
		_, err := os.Stat(this.filename)
		if err != nil {
			os.Create(this.filename)
		}
		f, e := os.OpenFile(this.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if e != nil {
			panic(e)
		}
		this.file = f
	}
	this.file.WriteString(msg)
	this.file.WriteString("\n")
}
