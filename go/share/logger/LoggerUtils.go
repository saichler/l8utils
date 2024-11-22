package logger

import (
	"bytes"
	"github.com/saichler/shared/go/share/interfaces"
	"github.com/saichler/shared/go/share/string_utils"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func FormatLog(level interfaces.LogLevel, t int64, args ...interface{}) string {
	str := string_utils.New()
	str.Add(LogTimeFormat(t, level))
	if args != nil {
		for _, arg := range args {
			str.Add(str.StringOf(arg))
		}
	}
	return str.String()
}

func LogTimeFormat(epochSeconds int64, level interfaces.LogLevel) string {
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
	buff.WriteString(level.String())
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

func FileAndLine(path string) string {
	filename, line := findFileAndLine(path)
	buff := bytes.Buffer{}
	buff.WriteString(" (")
	buff.WriteString(filename)
	buff.WriteString(".")
	buff.WriteString(strconv.Itoa(line))
	buff.WriteString(")")
	return buff.String()
}

func findFileAndLine(path string) (string, int) {
	index := 2
	ok := true
	filename := "/Unknown"
	line := -1
	for ok == true {
		_, filename, line, ok = runtime.Caller(index)
		if strings.Contains(filename, path) {
			break
		}
		index++
	}
	index = strings.LastIndex(filename, "/")
	return filename[index+1:], line
}
