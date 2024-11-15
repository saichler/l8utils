package logger

import (
	"bytes"
	"github.com/saichler/shared/go/share/string_utils"
	"strconv"
	"time"
)

type LogLevel string

const (
	Trace   LogLevel = "(  Trace)"
	Debug   LogLevel = "(  Debug)"
	Info    LogLevel = "(   Info)"
	Warning LogLevel = "(Warning)"
	Error   LogLevel = "(  Error)"
)

func FormatLog(level LogLevel, args ...interface{}) string {
	str := string_utils.New()
	str.Add(LogTimeFormat(time.Now().Unix(), level))
	if args != nil {
		for _, arg := range args {
			str.Add(str.StringOf(arg))
		}
	}
	return str.String()
}

func LogTimeFormat(epochSeconds int64, level LogLevel) string {
	t := time.Unix(epochSeconds, 0)
	buff := bytes.Buffer{}
	buff.WriteString(strconv.Itoa(t.Year()))
	buff.WriteString("-")
	buff.WriteString(intToString(int(t.Month())))
	buff.WriteString("-")
	buff.WriteString(intToString(t.Day()))
	buff.WriteString(" ")
	buff.WriteString(intToString(t.Hour()))
	buff.WriteString(":")
	buff.WriteString(intToString(t.Minute()))
	buff.WriteString(":")
	buff.WriteString(intToString(t.Second()))
	buff.WriteString(" ")
	buff.WriteString(string(level))
	buff.WriteString(" - ")
	return buff.String()
}

func intToString(i int) string {
	s := strconv.Itoa(i)
	if i < 10 {
		buff := bytes.Buffer{}
		buff.WriteString("0")
		buff.WriteString(s)
		return buff.String()
	}
	return s
}
