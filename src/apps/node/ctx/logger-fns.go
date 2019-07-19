package ctx

import (
	"github.com/cihub/seelog"
)

type defaultLoggerGetter struct{}

func (self *defaultLoggerGetter) GetLogger() seelog.LoggerInterface { return Ctx().GetLogger("app") }

type httpLoggerGetter struct{}

func (self *httpLoggerGetter) GetLogger() seelog.LoggerInterface { return Ctx().GetLogger("http") }

var (
	dfltLoggerGetter *defaultLoggerGetter = new(defaultLoggerGetter)
	hLoggerGetter    *httpLoggerGetter    = new(httpLoggerGetter)
)

func DefaultLogger() *defaultLoggerGetter { return dfltLoggerGetter }
func HTTPLogger() *httpLoggerGetter       { return hLoggerGetter }
