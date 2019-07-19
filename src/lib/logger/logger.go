package logger

type Logger interface {
	Debugf(format string, params ...interface{})
}
