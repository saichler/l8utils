package interfaces

type LogLevel int

const (
	Trace_Level   LogLevel = 1
	Debug_Level   LogLevel = 2
	Info_Level    LogLevel = 3
	Warning_Level LogLevel = 4
	Error_Level   LogLevel = 5
)

func (l LogLevel) String() string {
	switch l {
	case Trace_Level:
		return "(Trace)"
	case Debug_Level:
		return "(Debug)"
	case Info_Level:
		return "(Info) "
	case Warning_Level:
		return "(Warn )"
	case Error_Level:
		return "(Error)"
	}
	return ""
}

type ILogger interface {
	Trace(...interface{})
	Debug(...interface{})
	Info(...interface{})
	Warning(...interface{})
	Error(...interface{}) error
	Empty() bool
	Fail(interface{}, ...interface{})
	SetLogLevel(LogLevel)
	LoggerLock()
	LoggerUnlock()
	EnableLoggerSync(bool)
}

var logger ILogger

func SetLogger(l ILogger) {
	logger = l
}

func Logger() ILogger {
	return logger
}

func Trace(args ...interface{}) {
	logger.Trace(args...)
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Warning(args ...interface{}) {
	logger.Warning(args...)
}

func Error(args ...interface{}) error {
	return logger.Error(args...)
}

func Empty() bool {
	return logger.Empty()
}

func Fail(t interface{}, args ...interface{}) {
	logger.Fail(t, args...)
}
