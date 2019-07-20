package memstats

import (
	"time"

	"github.com/vany-egorov/ha-eta/lib/log"
)

const (
	defaultPeriod = 60 * time.Second
)

var (
	defaultFnLog = log.LogStd
)
