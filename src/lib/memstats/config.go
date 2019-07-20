package memstats

import (
	"time"

	"github.com/vany-egorov/ha-eta/lib/log"
)

type Config struct {
	period time.Duration
	fnLog  func(log.Level, string)
}

func (it *Config) defaultize() {
	it.period = defaultPeriod
	it.fnLog = defaultFnLog
}

type Arg func(*Config)

func FnLog(v func(log.Level, string)) Arg {
	return func(cfg *Config) { cfg.fnLog = v }
}

func Period(v time.Duration) Arg {
	return func(cfg *Config) { cfg.period = v }
}
