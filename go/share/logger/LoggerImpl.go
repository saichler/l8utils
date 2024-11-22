package logger

import (
	"errors"
	"github.com/saichler/shared/go/share/interfaces"
	"github.com/saichler/shared/go/share/queues"
	"sync"
	"testing"
	"time"
)

type ILogMethod interface {
	Log(level interfaces.LogLevel, msg string)
}

type LoggerImpl struct {
	queue      *queues.Queue
	logMethods []ILogMethod
	logLevel   interfaces.LogLevel
	mtx        *sync.Mutex
	enableMtx  bool
}

type LoggerEntry struct {
	anys []interface{}
	t    int64
	l    interfaces.LogLevel
}

func NewLoggerImpl(logMethods ...ILogMethod) *LoggerImpl {
	logImpl := &LoggerImpl{}
	logImpl.logMethods = logMethods
	logImpl.queue = queues.NewQueue("Logger Queue", 50000)
	logImpl.mtx = &sync.Mutex{}
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

func newEntry(l interfaces.LogLevel, anys ...interface{}) *LoggerEntry {
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
	if loggerImpl.logLevel > interfaces.Trace_Level {
		return
	}
	loggerImpl.queue.Add(newEntry(interfaces.Trace_Level, anys...))
}

func (loggerImpl *LoggerImpl) Debug(anys ...interface{}) {
	if loggerImpl.logLevel > interfaces.Debug_Level {
		return
	}
	loggerImpl.queue.Add(newEntry(interfaces.Debug_Level, anys...))
}

func (loggerImpl *LoggerImpl) Info(anys ...interface{}) {
	if loggerImpl.logLevel > interfaces.Info_Level {
		return
	}
	loggerImpl.queue.Add(newEntry(interfaces.Info_Level, anys...))
}

func (loggerImpl *LoggerImpl) Warning(anys ...interface{}) {
	if loggerImpl.logLevel > interfaces.Warning_Level {
		return
	}
	loggerImpl.queue.Add(newEntry(interfaces.Warning_Level, anys...))
}

func (loggerImpl *LoggerImpl) Error(anys ...interface{}) error {
	loggerImpl.queue.Add(newEntry(interfaces.Error_Level, anys...))
	err := FormatLog(interfaces.Error_Level, time.Now().Unix(), anys...)
	return errors.New(err)
}

func (loggerImpl *LoggerImpl) Fail(t interface{}, args ...interface{}) {
	args = append(args, FileAndLine("tests"))
	loggerImpl.Error(args...)
	ts, ok := t.(*testing.T)
	if ok {
		ts.Fail()
	}
}

func (loggerImpl *LoggerImpl) SetLogLevel(level interfaces.LogLevel) {
	loggerImpl.logLevel = level
}

func (loggerImpl *LoggerImpl) LoggerLock() {
	if loggerImpl.enableMtx {
		loggerImpl.mtx.Lock()
	}
}
func (loggerImpl *LoggerImpl) LoggerUnlock() {
	if loggerImpl.enableMtx {
		loggerImpl.mtx.Unlock()
	}
}
func (loggerImpl *LoggerImpl) EnableLoggerSync(enable bool) {
	loggerImpl.enableMtx = enable
}
