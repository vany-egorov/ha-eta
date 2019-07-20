package node

import "sync"

type Context struct {
	cfgLock sync.RWMutex
	config  *config
}

func (it *Context) setCfg(v *config) *config {
	it.cfgLock.Lock()
	defer it.cfgLock.Unlock()
	old := it.config
	it.config = v
	return old
}

func (it *Context) cfg() *config {
	it.cfgLock.Lock()
	defer it.cfgLock.Unlock()
	return it.config
}
