package lib

import (
	"sync"

	"github.com/cihub/seelog"
)

type LoggerInterface struct {
	seelog.LoggerInterface
}

type LoggersMap struct {
	sync.RWMutex
	m map[string]seelog.LoggerInterface
}

func (it *LoggersMap) Get(name string) seelog.LoggerInterface {
	it.RLock()
	defer it.RUnlock()
	return it.m[name]
}

func (it *LoggersMap) Store(name string, logger seelog.LoggerInterface) *LoggersMap {
	it.Lock()
	defer it.Unlock()
	it.m[name] = logger
	return it
}

func (it *LoggersMap) Close() *LoggersMap {
	it.Lock()
	defer it.Unlock()
	for _, logger := range it.m {
		logger.Close()
	}
	return it
}

func NewLoggersMap() *LoggersMap {
	it := &LoggersMap{
		m: make(map[string]seelog.LoggerInterface),
	}
	return it
}
