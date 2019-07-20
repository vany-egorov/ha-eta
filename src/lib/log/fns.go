package log

import (
	"io"
	"io/ioutil"
	"os"
	"strings"

	"go.uber.org/zap"
)

var (
	nope      = func(string) {}
	nopeLevel = func(string) {}

	toWriter = func(w io.Writer, msg string) {
		b := strings.Builder{}
		b.WriteString(msg)
		b.WriteString("\n")
		w.Write([]byte(b.String()))
	}

	toDiscard = func(msg string) { toWriter(ioutil.Discard, msg) }
	toStderr  = func(msg string) { toWriter(os.Stderr, msg) }
	toStdout  = func(msg string) { toWriter(os.Stdout, msg) }
)

func LogFn(any logger, l Level) func(string) {
	switch lgr := any.(type) {
	case *zap.Logger:
		return func(msg string) {
			switch l {
			case Quiet:
			case Trace, Debug:
				lgr.Debug(msg)
			case Info:
				lgr.Info(msg)
			case Warn:
				lgr.Warn(msg)
			case Error:
				lgr.Warn(msg) // TODO: research for non-panicable Error for zap
			case Critical:
				lgr.Fatal(msg)
			}

			lgr.Sync()
		}
	}

	switch l {
	case Quiet:
		return toDiscard
	case Trace, Debug, Info:
		return toStdout
	case Warn, Error, Critical:
		return toStderr
	}

	return nope
}

func Log(any logger, l Level, msg string) { LogFn(any, l)(msg) }

func LogStd(l Level, msg string) { LogFn(nil, l)(msg) }

func LogFnWithLogger(any logger) func(l Level, msg string) {
	return func(l Level, msg string) {
		LogFn(any, l)(msg)
	}
}
