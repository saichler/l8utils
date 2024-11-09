package interfaces

type ILogger interface {
	Trace(...interface{})
	Debug(...interface{})
	Info(...interface{})
	Warning(...interface{})
	Error(...interface{}) error
	Empty() bool
	Fail(...interface{})
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

func Fail(args ...interface{}) {
	logger.Fail(args...)
}
