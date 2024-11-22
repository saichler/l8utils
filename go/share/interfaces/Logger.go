package interfaces

type ILogger interface {
	Trace(...interface{})
	Debug(...interface{})
	Info(...interface{})
	Warning(...interface{})
	Error(...interface{}) error
	Empty() bool
	Fail(interface{}, ...interface{})

	IsTraceEnabled() bool
	IsDebugEnabled() bool
	IsInfoEnabled() bool
	IsWarningEnabled() bool
}

var logger ILogger

func SetLogger(l ILogger) {
	logger = l
}

func Logger() ILogger {
	return logger
}

func Trace(args ...interface{}) {
	if logger.IsTraceEnabled() {
		logger.Trace(args...)
	}
}

func Debug(args ...interface{}) {
	if logger.IsDebugEnabled() {
		logger.Debug(args...)
	}
}

func Info(args ...interface{}) {
	if logger.IsInfoEnabled() {
		logger.Info(args...)
	}
}

func Warning(args ...interface{}) {
	if logger.IsWarningEnabled() {
		logger.Warning(args...)
	}
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
