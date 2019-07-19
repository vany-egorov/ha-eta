package ctx

import (
	"fmt"
	"sync"
	"time"

	"github.com/cihub/seelog"

	"github.com/vany-egorov/ha-eta/apps/node/initializers"
	"github.com/vany-egorov/ha-eta/lib"
	"github.com/vany-egorov/ha-eta/lib/logger"
)

var ctx *Context

type Context struct {
	version string
	doneAt  *lib.DoneAt

	configMutex sync.RWMutex
	config      *initializers.Config
	flags       *initializers.Flags

	loggersMutex sync.RWMutex
	loggers      *lib.LoggersMap
}

func (it *Context) Version() string     { return it.version }
func (it *Context) SetVersion(v string) { it.version = v }

func (it *Context) DoneAt() *lib.DoneAt              { return it.doneAt }
func (it *Context) SetDoneAt(v *lib.DoneAt) *Context { it.doneAt = v; return it }

func (it *Context) Flags() *initializers.Flags              { return it.flags }
func (it *Context) SetFlags(v *initializers.Flags) *Context { it.flags = v; return it }

func (it *Context) Config() *initializers.Config {
	it.configMutex.RLock()
	defer it.configMutex.RUnlock()
	return it.config
}
func (it *Context) SetConfig(v *initializers.Config) *Context {
	it.configMutex.Lock()
	defer it.configMutex.Unlock()
	it.config = v
	return it
}

func (it *Context) GetLogger(name string) seelog.LoggerInterface {
	it.loggersMutex.RLock()
	defer it.loggersMutex.RUnlock()
	return it.loggers.Get(name)
}
func (it *Context) Logger(name string) seelog.LoggerInterface { return it.GetLogger(name) }
func (it *Context) Loggers() *lib.LoggersMap {
	it.loggersMutex.Lock()
	defer it.loggersMutex.Unlock()
	return it.loggers
}
func (it *Context) SetLoggers(v *lib.LoggersMap) *Context {
	it.loggersMutex.Lock()
	defer it.loggersMutex.Unlock()
	it.loggers = v
	return it
}
func (it *Context) DefaultLogger() logger.Getter { return DefaultLogger() }
func (it *Context) HTTPLogger() logger.Getter    { return HTTPLogger() }

func (it *Context) ConfigPeriodMemstats() time.Duration { return it.Config().Period.Memstats }

func Ctx() *Context {
	if ctx == nil {
		ctx = NewCtx()
	}
	return ctx
}

func NewCtx() *Context {
	ctx = new(Context)
	ctx.version = "MAJOR.MINOR.PATCH.YYYYMMDD-HHddSS ~ UNKNOWN"
	return ctx
}

func LoadLoggers() error {
	oldLoggers := ctx.Loggers()

	ctx.Config().GetLog().Defaultize(ctx.Config().PathLogIsSet(), ctx.Config().PathLog(), ctx.Config().Environment)
	if loggers, e := ctx.Config().GetLog().ToLoggersMap(nil); e != nil {
		return fmt.Errorf("loggers initialization failed: %s", e.Error())
	} else {
		ctx.SetLoggers(loggers)
	}

	if oldLoggers != nil {
		oldLoggers.Close()
	}

	return nil
}
