package logger

import (
	"errors"
	"fmt"
	"log"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8utils/go/utils/queues"
)

type ILogMethod interface {
	Log(level ifs.LogLevel, msg string)
}

type LoggerImpl struct {
	queue      *queues.Queue
	logMethods []ILogMethod
	logLevel   ifs.LogLevel
}

type LoggerEntry struct {
	anys []interface{}
	t    int64
	l    ifs.LogLevel
}

func NewLoggerImpl(logMethods ...ILogMethod) *LoggerImpl {
	logImpl := &LoggerImpl{}
	logImpl.logMethods = logMethods
	logImpl.queue = queues.NewQueue("Logger Queue", 50000)
	go logImpl.processQueue()
	return logImpl
}

func (loggerImpl *LoggerImpl) processQueue() {
	for {
		entry := loggerImpl.queue.Next().(*LoggerEntry)
		str := FormatLog(entry.l, entry.t, entry.anys...)
		if len(loggerImpl.logMethods) == 1 {
			loggerImpl.logMethods[0].Log(entry.l, str)
		} else if len(loggerImpl.logMethods) == 2 {
			loggerImpl.logMethods[0].Log(entry.l, str)
			loggerImpl.logMethods[1].Log(entry.l, str)
		} else if len(loggerImpl.logMethods) == 3 {
			loggerImpl.logMethods[0].Log(entry.l, str)
			loggerImpl.logMethods[1].Log(entry.l, str)
			loggerImpl.logMethods[2].Log(entry.l, str)
		}
	}
}

func newEntry(l ifs.LogLevel, anys ...interface{}) *LoggerEntry {
	return &LoggerEntry{
		t:    time.Now().Unix(),
		l:    l,
		anys: anys,
	}
}

func (loggerImpl *LoggerImpl) Empty() bool {
	return loggerImpl.queue.Size() == 0
}

func (loggerImpl *LoggerImpl) Trace(anys ...interface{}) {
	if loggerImpl.logLevel > ifs.Trace_Level {
		return
	}
	loggerImpl.queue.Add(newEntry(ifs.Trace_Level, anys...))
}

func (loggerImpl *LoggerImpl) Debug(anys ...interface{}) {
	if loggerImpl.logLevel > ifs.Debug_Level {
		return
	}
	loggerImpl.queue.Add(newEntry(ifs.Debug_Level, anys...))
}

func (loggerImpl *LoggerImpl) Info(anys ...interface{}) {
	if loggerImpl.logLevel > ifs.Info_Level {
		return
	}
	loggerImpl.queue.Add(newEntry(ifs.Info_Level, anys...))
}

func (loggerImpl *LoggerImpl) Warning(anys ...interface{}) {
	if loggerImpl.logLevel > ifs.Warning_Level {
		return
	}
	loggerImpl.queue.Add(newEntry(ifs.Warning_Level, anys...))
}

func (loggerImpl *LoggerImpl) Error(anys ...interface{}) error {
	anys = append(anys, FileAndLine(".go", false))
	loggerImpl.queue.Add(newEntry(ifs.Error_Level, anys...))
	err := FormatLog(ifs.Error_Level, time.Now().Unix(), anys...)
	return errors.New(err)
}

func (loggerImpl *LoggerImpl) Fail(t interface{}, args ...interface{}) {
	args = append(args, FileAndLine("tests", true))
	loggerImpl.Error(args...)
	ts, ok := t.(*testing.T)
	if ok {
		ts.Fail()
	}
}

func (loggerImpl *LoggerImpl) SetLogLevel(level ifs.LogLevel) {
	loggerImpl.logLevel = level
}

func init() {
	os.MkdirAll("/data/logs", 0777)
	hostname := os.Getenv("HOSTNAME")
	if hostname == "" {
		hostname = "localhost"
	}

	uuid := ifs.NewUuid()
	panicFileName := "/data/logs/" + hostname + "-" + uuid + ".err"
	panicFile, err := os.Create(panicFileName)

	logFileName := "/data/logs/" + hostname + "-" + uuid + ".log"
	logFile, err := os.Create(logFileName)

	if err == nil {
		err = syscall.Dup2(int(panicFile.Fd()), int(os.Stderr.Fd()))
		if err != nil {
			log.Fatalf("Failed to redirect stderr: %v", err)
		}
		err = syscall.Dup2(int(logFile.Fd()), int(os.Stdout.Fd()))
		if err != nil {
			log.Fatalf("Failed to redirect stdout: %v", err)
		}
	} else {
		fmt.Println("Failed to create error file:", err)
	}
}
